package models

type TaskStatus int64

const (
	TaskStatusCreated        TaskStatus = 0
	TaskStatusProcessing     TaskStatus = 1
	TaskStatusCompleted      TaskStatus = 2
	TaskStatusFailed         TaskStatus = 3
	TaskStatusNoAttemptsLeft TaskStatus = 4
)
