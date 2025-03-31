package model

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type TaskUpdate struct {
	ID             uuid.UUID
	Original       string
	Todo           string
	Priority       string
	Projects       []string
	Contexts       []string
	AdditionalTags map[string]string
	CreatedDate    time.Time
	DueDate        time.Time
	CompletedDate  time.Time
	Completed      bool
	UsedTime       time.Duration
	Source         string
}

type TaskStartTime struct {
	ID        uuid.UUID
	TaskID    uuid.UUID
	StartTime time.Time
}

type TaskEndTime struct {
	ID      uuid.UUID
	TaskID  uuid.UUID
	EndTime time.Time
}

type SetTimeBudgetForProject struct {
	ID          uuid.UUID
	ProjectID   uuid.UUID
	Budget      time.Duration
	UsedTime    time.Duration
	Description string
}