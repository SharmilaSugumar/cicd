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
	Name string `gorm:"not null"`

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
	OrganizationID uuid.UUID `gorm:"type:uuid;index;not null"`
	Name           string    `gorm:"not null"`
	Description    string

	Pipelines    []Pipeline   `gorm:"foreignKey:ProjectID"`
	Organization Organization `gorm:"foreignKey:OrganizationID"`
}

type Queue struct {
	BaseModel
	Name   string      `gorm:"not null"`
	Status QueueStatus `gorm:"type:varchar(20);index;not null;default:'ACTIVE'"`
}

type Pipeline struct {
	BaseModel
	ProjectID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Name        string    `gorm:"not null"`
	Description string

	Runs    []PipelineRun `gorm:"foreignKey:PipelineID"`
	Project Project       `gorm:"foreignKey:ProjectID"`
}

type PipelineRun struct {
	BaseModel
	PipelineID uuid.UUID         `gorm:"type:uuid;index;not null"`
	Status     PipelineRunStatus `gorm:"type:varchar(20);index;not null;default:'PENDING'"`

	Jobs     []Job    `gorm:"foreignKey:PipelineRunID"`
	Pipeline Pipeline `gorm:"foreignKey:PipelineID"`
}

type Job struct {
	BaseModel
	PipelineRunID uuid.UUID  `gorm:"type:uuid;index;not null"`
	QueueID       *uuid.UUID `gorm:"type:uuid;index"`
	Name          string     `gorm:"not null"`
	Status        JobStatus  `gorm:"type:varchar(20);index;not null;default:'CREATED'"`
	Priority      int        `gorm:"index;default:0"`
	Payload       string     // JSON payload for execution

	PipelineRun PipelineRun `gorm:"foreignKey:PipelineRunID"`
	Queue       *Queue      `gorm:"foreignKey:QueueID"`
	Logs        []JobLog    `gorm:"foreignKey:JobID"`
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
	JobID   uuid.UUID `gorm:"type:uuid;index;not null"`
	Message string    `gorm:"not null"`

	Job Job `gorm:"foreignKey:JobID"`
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
	JobID        uuid.UUID `gorm:"type:uuid;index;not null"`
	ErrorMessage string
	Payload      string

	Job Job `gorm:"foreignKey:JobID"`
}
