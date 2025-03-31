package model

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Description string
	Tasks       []Task
	Budget      time.Duration
}
