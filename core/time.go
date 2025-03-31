package core

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

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
