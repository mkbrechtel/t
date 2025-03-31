package core

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type Project struct {
	ID          uuid.UUID
	Name        string
	Budget      time.Duration
}
