package services

import "errors"

var (
	ErrNotFound               = errors.New("resource not found")
	ErrUnauthorized           = errors.New("unauthorized access")
	ErrForbidden              = errors.New("forbidden action")
	ErrInvalidInput           = errors.New("invalid input")
	ErrDuplicateResource      = errors.New("resource already exists")
	ErrInternalError          = errors.New("internal server error")
	ErrLimitExceeded          = errors.New("concurrency limit exceeded")
	ErrDependencyNotMet       = errors.New("job dependencies not met")
	ErrInvalidStateTransition = errors.New("invalid state transition")
)
