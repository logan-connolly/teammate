package team

import (
	"reflect"

	"github.com/google/uuid"
)

// Event is a domain event marker.
type Event interface {
	eventName() string
}

// TeamRegistered event.
type TeamRegistered struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e TeamRegistered) eventName() string {
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