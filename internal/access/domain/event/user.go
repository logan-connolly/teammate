package event

import (
	"reflect"

	"github.com/google/uuid"
)

// UserRegistered event.
type UserRegistered struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}

func (e UserRegistered) eventName() string {
	return reflect.TypeOf(e).Name()
}

// UserNameChanged event.
type UserNameChanged struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (e UserNameChanged) eventName() string {
	return reflect.TypeOf(e).Name()
}

// UserEmailChanged event.
type UserEmailChanged struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

func (e UserEmailChanged) eventName() string {
	return reflect.TypeOf(e).Name()
}

// UserActivated event.
type UserActivated struct {
	ID uuid.UUID `json:"id"`
}

func (e UserActivated) eventName() string {
	return reflect.TypeOf(e).Name()
}

// UserDeactivated event.
type UserDeactivated struct {
	ID uuid.UUID `json:"id"`
}

func (e UserDeactivated) eventName() string {
	return reflect.TypeOf(e).Name()
}
