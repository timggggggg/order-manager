package models

import "time"

type TaskStatus int64

const (
	TaskStatusCreated        TaskStatus = 0
	TaskStatusProcessing     TaskStatus = 1
	TaskStatusCompleted      TaskStatus = 2
	TaskStatusFailed         TaskStatus = 3
	TaskStatusNoAttemptsLeft TaskStatus = 4
)

type TaskType int64

const (
	TaskTypeUnknown         TaskType = 0
	TaskTypeStatusChangeLog TaskType = 1
	TaskTypeHttpReqLog      TaskType = 2
	TaskTypeHttpRespLog     TaskType = 3
)

type OutboxTask struct {
	ID           int64
	Type         TaskType
	Payload      []byte
	Status       TaskStatus
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CompletedAt  time.Time
	AttemptsLeft int64
}

func NewOutboxTask(taskType TaskType, payload []byte, attempts int64) *OutboxTask {
	return &OutboxTask{
		Type:         taskType,
		Payload:      payload,
		Status:       TaskStatusCreated,
		AttemptsLeft: attempts,
	}
}
