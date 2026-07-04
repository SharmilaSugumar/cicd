package database

type JobStatus string

const (
	JobStatusCreated    JobStatus = "CREATED"
	JobStatusQueued     JobStatus = "QUEUED"
	JobStatusClaimed    JobStatus = "CLAIMED"
	JobStatusRunning    JobStatus = "RUNNING"
	JobStatusCompleted  JobStatus = "COMPLETED"
	JobStatusFailed     JobStatus = "FAILED"
	JobStatusRetrying   JobStatus = "RETRYING"
	JobStatusDeadLetter JobStatus = "DEAD_LETTER"
)

type WorkerStatus string

const (
	WorkerStatusOnline  WorkerStatus = "ONLINE"
	WorkerStatusOffline WorkerStatus = "OFFLINE"
)

type QueueStatus string

const (
	QueueStatusActive QueueStatus = "ACTIVE"
	QueueStatusPaused QueueStatus = "PAUSED"
)

type PipelineRunStatus string

const (
	PipelineRunStatusPending PipelineRunStatus = "PENDING"
	PipelineRunStatusRunning PipelineRunStatus = "RUNNING"
	PipelineRunStatusSuccess PipelineRunStatus = "SUCCESS"
	PipelineRunStatusFailed  PipelineRunStatus = "FAILED"
)

type OrganizationRole string

const (
	OrganizationRoleOwner      OrganizationRole = "OWNER"
	OrganizationRoleMaintainer OrganizationRole = "MAINTAINER"
	OrganizationRoleDeveloper  OrganizationRole = "DEVELOPER"
	OrganizationRoleViewer     OrganizationRole = "VIEWER"
)
