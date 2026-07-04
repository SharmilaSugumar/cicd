package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	QueueLength = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "forgeflow_queue_length",
		Help: "Current number of jobs in the queue",
	}, []string{"queue_id"})

	RunningJobs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "forgeflow_running_jobs",
		Help: "Current number of running jobs",
	})

	ActiveWorkers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "forgeflow_active_workers",
		Help: "Current number of active workers",
	})

	Throughput = promauto.NewCounter(prometheus.CounterOpts{
		Name: "forgeflow_jobs_processed_total",
		Help: "Total number of processed jobs",
	})

	Retries = promauto.NewCounter(prometheus.CounterOpts{
		Name: "forgeflow_job_retries_total",
		Help: "Total number of job retries",
	})

	ExecutionDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "forgeflow_job_execution_duration_seconds",
		Help:    "Duration of job executions in seconds",
		Buckets: prometheus.DefBuckets,
	})
)
