package workers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"forgeflow/internal/database"
	"forgeflow/internal/execution"
	"forgeflow/internal/repositories"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Worker interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type worker struct {
	id         uuid.UUID
	name       string
	queueID    *uuid.UUID
	workerRepo repositories.WorkerRepository
	jobRepo    repositories.JobRepository
	db         *gorm.DB
	engine     execution.Engine
	pollInt    time.Duration
	hbInt      time.Duration
	cancel     context.CancelFunc
	wg         sync.WaitGroup
}

type WorkerConfig struct {
	Name              string
	QueueID           *uuid.UUID
	PollInterval      time.Duration
	HeartbeatInterval time.Duration
}

func NewWorker(cfg WorkerConfig, workerRepo repositories.WorkerRepository, jobRepo repositories.JobRepository, db *gorm.DB, engine execution.Engine) Worker {
	if cfg.PollInterval == 0 {
		cfg.PollInterval = 5 * time.Second
	}
	if cfg.HeartbeatInterval == 0 {
		cfg.HeartbeatInterval = 10 * time.Second
	}

	return &worker{
		id:         uuid.New(),
		name:       cfg.Name,
		queueID:    cfg.QueueID,
		workerRepo: workerRepo,
		jobRepo:    jobRepo,
		db:         db,
		engine:     engine,
		pollInt:    cfg.PollInterval,
		hbInt:      cfg.HeartbeatInterval,
	}
}

func (w *worker) Start(ctx context.Context) error {
	wCtx, cancel := context.WithCancel(ctx)
	w.cancel = cancel

	dbWorker := &database.Worker{
		BaseModel: database.BaseModel{ID: w.id},
		Name:      w.name,
		Status:    database.WorkerStatusOnline,
	}

	// Register
	if err := w.workerRepo.Register(wCtx, dbWorker); err != nil {
		return fmt.Errorf("failed to register worker: %w", err)
	}

	log.Printf("Worker %s (%s) registered successfully", w.name, w.id)

	w.wg.Add(2)
	go w.heartbeatLoop(wCtx)
	go w.pollLoop(wCtx)

	return nil
}

func (w *worker) Stop(ctx context.Context) error {
	log.Printf("Worker %s stopping...", w.id)
	if w.cancel != nil {
		w.cancel()
	}

	// wait for loops with timeout
	c := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(c)
	}()

	select {
	case <-c:
	case <-ctx.Done():
		return ctx.Err()
	}

	// Mark offline
	offlineCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := w.workerRepo.MarkOffline(offlineCtx, w.id); err != nil {
		return fmt.Errorf("failed to mark worker offline: %w", err)
	}

	log.Printf("Worker %s stopped.", w.id)
	return nil
}

func (w *worker) heartbeatLoop(ctx context.Context) {
	defer w.wg.Done()
	ticker := time.NewTicker(w.hbInt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := w.workerRepo.UpdateHeartbeat(ctx, w.id); err != nil {
				log.Printf("Worker %s failed to update heartbeat: %v", w.id, err)
			}
		}
	}
}

func (w *worker) pollLoop(ctx context.Context) {
	defer w.wg.Done()
	ticker := time.NewTicker(w.pollInt)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.claimAndExecute(ctx)
		}
	}
}

func (w *worker) claimAndExecute(ctx context.Context) {
	job, err := w.jobRepo.ClaimNextJob(ctx, w.queueID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// No jobs to claim
			return
		}
		log.Printf("Error claiming job: %v", err)
		return
	}

	log.Printf("Claimed job %s", job.ID)

	// Update Status to Running
	if err := w.jobRepo.UpdateStatus(ctx, job.ID, database.JobStatusRunning); err != nil {
		log.Printf("Error updating job status to running: %v", err)
		return
	}

	// Execute
	result := w.engine.ExecuteJob(ctx, job.Payload)

	// Upload Logs
	if w.db != nil {
		msg := fmt.Sprintf("Exit Code: %d\nDuration: %v\nError: %v\nStdout:\n%s\nStderr:\n%s",
			result.ExitCode, result.Duration, result.Error, result.Stdout, result.Stderr)
		logRecord := database.JobLog{
			JobID:   job.ID,
			Message: msg,
		}
		if err := w.db.WithContext(ctx).Create(&logRecord).Error; err != nil {
			log.Printf("Error uploading job log: %v", err)
		}
	}

	// Update Status
	finalStatus := database.JobStatusCompleted
	if result.ExitCode != 0 || result.Error != nil {
		finalStatus = database.JobStatusFailed
	}

	if err := w.jobRepo.UpdateStatus(ctx, job.ID, finalStatus); err != nil {
		log.Printf("Error updating final job status: %v", err)
	}

	log.Printf("Finished job %s with status %s", job.ID, finalStatus)
}
