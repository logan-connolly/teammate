package entity

import (
	"github.com/google/uuid"
)

// Group is an entity that represents a group in all domains.
type Group struct {
	ID   uuid.UUID
	Name string
}
