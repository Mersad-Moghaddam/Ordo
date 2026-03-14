package task

import "errors"

var ErrTaskNotFound = errors.New("task not found")
var ErrTaskAlreadyExists = errors.New("task already exists")
var ErrInvalidTaskStatusTransition = errors.New("invalid task status transition")
var ErrOutboxPersistFailure = errors.New("outbox persist failure")
