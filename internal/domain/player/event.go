package player

import "github.com/google/uuid"

// Event is a domain event marker.
type Event interface {
	isEvent()
}

func (e PlayerRegistered) isEvent()  {}
func (e PlayerActivated) isEvent()   {}
func (e PlayerDeactivated) isEvent() {}

// PlayerRegistered event.
type PlayerRegistered struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// PlayerActivated event.
type PlayerActivated struct {
	ID uuid.UUID `json:"id"`
}

// PlayerDeactivated event.
type PlayerDeactivated struct {
	ID uuid.UUID `json:"id"`
}
