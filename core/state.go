package core

import uuid "github.com/gofrs/uuid/v5"

type AppState struct {
	Tasks map[uuid.UUID]Task
	Projects map[string]Project
}
