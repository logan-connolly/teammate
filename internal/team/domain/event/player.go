package event

import (
	"reflect"

	"github.com/google/uuid"
)

// PlayerRegistered event.
type PlayerRegistered struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e PlayerRegistered) eventName() string {
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
