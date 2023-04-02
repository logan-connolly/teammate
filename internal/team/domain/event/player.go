package event

import (
	"reflect"

	"github.com/google/uuid"
)

// PlayerCreated event.
type PlayerCreated struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e PlayerCreated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// PlayerActivated event.
type PlayerActivated struct {
	ID uuid.UUID `json:"id"`
}

func (e PlayerActivated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// PlayerDeactivated event.
type PlayerDeactivated struct {
	ID uuid.UUID `json:"id"`
}

func (e PlayerDeactivated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// TeamAssignedToPlayer event.
type TeamAssignedToPlayer struct {
	ID       uuid.UUID `json:"id"`
	TeamId   uuid.UUID `json:"team_id"`
	TeamName string    `json:"team_name"`
}

func (e TeamAssignedToPlayer) eventName() string {
	return reflect.TypeOf(e).Name()
}

// TeamUnassignedFromPlayer event.
type TeamUnassignedFromPlayer struct {
	ID       uuid.UUID `json:"id"`
	TeamId   uuid.UUID `json:"team_id"`
	TeamName string    `json:"team_name"`
}

func (e TeamUnassignedFromPlayer) eventName() string {
	return reflect.TypeOf(e).Name()
}
