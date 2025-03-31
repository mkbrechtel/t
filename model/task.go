package model

import (
	"time"

	uuid "github.com/gofrs/uuid/v5"
)

type Task struct {
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
