package database

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BaseModel provides UUID primary key, timestamps, and soft deletes
type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	BaseModel
	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Name         string

	Memberships []OrganizationMember `gorm:"foreignKey:UserID"`
}

type Organization struct {
	BaseModel
	Name string `gorm:"not null" json:"name"`

	Members  []OrganizationMember `gorm:"foreignKey:OrganizationID"`
	Projects []Project            `gorm:"foreignKey:OrganizationID"`
}

type OrganizationMember struct {
	BaseModel
	OrganizationID uuid.UUID        `gorm:"type:uuid;index;not null"`
	UserID         uuid.UUID        `gorm:"type:uuid;index;not null"`
	Role           OrganizationRole `gorm:"type:varchar(20);not null"`

	Organization Organization `gorm:"foreignKey:OrganizationID"`
	User         User         `gorm:"foreignKey:UserID"`
}

type Project struct {
	BaseModel
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null" json:"organization_id"`
	Name           string    `gorm:"not null" json:"name"`
	Description    string    `json:"description"`

	Pipelines    []Pipeline   `gorm:"foreignKey:ProjectID"`
	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

type Queue struct {
	BaseModel
	Name   string      `gorm:"not null" json:"name"`
	Status QueueStatus `gorm:"type:varchar(20);index;not null;default:'ACTIVE'" json:"status"`
}

type Pipeline struct {
	BaseModel
	ProjectID   uuid.UUID `gorm:"type:uuid;index;not null" json:"project_id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	YamlConfig  string    `json:"yaml_config"`

	Runs    []PipelineRun `gorm:"foreignKey:PipelineID"`
	Project Project       `gorm:"foreignKey:ProjectID"`
}

type PipelineRun struct {
	BaseModel
	PipelineID uuid.UUID         `gorm:"type:uuid;index;not null" json:"pipeline_id"`
	Status     PipelineRunStatus `gorm:"type:varchar(20);index;not null;default:'PENDING'" json:"status"`

	Jobs     []Job    `gorm:"foreignKey:PipelineRunID"`
	Pipeline Pipeline `gorm:"foreignKey:PipelineID"`
}

type Job struct {
	BaseModel
	PipelineRunID uuid.UUID  `gorm:"type:uuid;index;not null" json:"pipeline_run_id"`
	QueueID       *uuid.UUID `gorm:"type:uuid;index" json:"queue_id"`
	Name          string     `gorm:"not null" json:"name"`
	Status        JobStatus  `gorm:"type:varchar(20);index;not null;default:'CREATED'" json:"status"`
	Priority      int        `gorm:"index;default:0" json:"priority"`
	Payload       string     `json:"payload"` // JSON payload for execution

	PipelineRun PipelineRun `gorm:"foreignKey:PipelineRunID"`
	Queue       *Queue           `gorm:"foreignKey:QueueID"`
	Logs        []JobLog         `gorm:"foreignKey:JobID" json:"logs"`
	DLQ         *DeadLetterQueue `gorm:"foreignKey:JobID" json:"dlq"`
}

type JobDependency struct {
	BaseModel
	ParentJobID uuid.UUID `gorm:"type:uuid;index;not null"`
	ChildJobID  uuid.UUID `gorm:"type:uuid;index;not null"`

	ParentJob Job `gorm:"foreignKey:ParentJobID"`
	ChildJob  Job `gorm:"foreignKey:ChildJobID"`
}

type JobLog struct {
	BaseModel
	JobID   uuid.UUID `gorm:"type:uuid;index;not null" json:"job_id"`
	Message string    `gorm:"not null" json:"message"`

	Job Job `gorm:"foreignKey:JobID" json:"-"`
}

type RetryPolicy struct {
	BaseModel
	JobID      uuid.UUID `gorm:"type:uuid;index;not null;unique"`
	MaxRetries int       `gorm:"default:3"`
	Delay      int       // In seconds

	Job Job `gorm:"foreignKey:JobID"`
}

type Worker struct {
	BaseModel
	Name   string       `gorm:"not null"`
	Status WorkerStatus `gorm:"type:varchar(20);index;not null;default:'OFFLINE'"`

	Heartbeats []WorkerHeartbeat `gorm:"foreignKey:WorkerID"`
}

type WorkerHeartbeat struct {
	BaseModel
	WorkerID uuid.UUID `gorm:"type:uuid;index;not null"`

	Worker Worker `gorm:"foreignKey:WorkerID"`
}

type DeadLetterQueue struct {
	BaseModel
	JobID        uuid.UUID `gorm:"type:uuid;index;not null" json:"job_id"`
	ErrorMessage string    `json:"error_message"`
	Payload      string    `json:"payload"`

	Job Job `gorm:"foreignKey:JobID" json:"-"`
}
