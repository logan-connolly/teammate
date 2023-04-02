package event

import (
	"reflect"

	"github.com/google/uuid"
)

// TeamCreated event.
type TeamCreated struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e TeamCreated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// TeamActivated event.
type TeamActivated struct {
	ID uuid.UUID `json:"id"`
}

func (e TeamActivated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// TeamDeactivated event.
type TeamDeactivated struct {
	ID uuid.UUID `json:"id"`
}

func (e TeamDeactivated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// PlayerAssignedToTeam event.
type PlayerAssignedToTeam struct {
	ID         uuid.UUID `json:"id"`
	PlayerId   uuid.UUID `json:"player_id"`
	PlayerName string    `json:"player_name"`
}

func (e PlayerAssignedToTeam) eventName() string {
	return reflect.TypeOf(e).Name()
}

// PlayerUnassignedFromTeam event.
type PlayerUnassignedFromTeam struct {
	ID         uuid.UUID `json:"id"`
	PlayerId   uuid.UUID `json:"player_id"`
	PlayerName string    `json:"player_name"`
}

func (e PlayerUnassignedFromTeam) eventName() string {
	return reflect.TypeOf(e).Name()
}
