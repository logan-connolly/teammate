package model

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/access/domain/event"
	"git.sr.ht/~loges/teammate/internal/entity"
	"github.com/google/uuid"
)

var (
	ErrInputIsEmpty     = errors.New("model: non-empty value must be provided")
	ErrUserUpdateFailed = errors.New("model: user update failed")
)

// User is a aggregate that combines all entities needed to represent a user.
type User struct {
	person    *entity.Person
	email     string
	activated bool

	changes []event.Event
	version int
}

// NewUser is a factory to create a new User aggregate.
func NewUser(p *entity.Person, email string) (*User, error) {
	user := &User{}

	if p.Name == "" || email == "" {
		return user, ErrInputIsEmpty
	}

	user.register(&event.UserRegistered{
		ID:    p.ID,
		Name:  p.Name,
		Email: email,
	})

	return user, nil
}

// NewFromEvents is a helper method that creates a new user
// from a series of events.
func NewUserFromEvents(events []event.Event) *User {
	u := &User{}

	for _, event := range events {
		u.Apply(event, false)
	}

	return u
}

// GetID returns the user root entity GetID.
func (u *User) GetID() uuid.UUID {
	return u.person.ID
}

// GetName returns the name of the user.
func (u *User) GetName() string {
	return u.person.Name
}

// GetEmail returns the email of the user.
func (u *User) GetEmail() string {
	return u.email
}

// IsActivated returns whether the user is activated.
func (u *User) IsActivated() bool {
	return u.activated
}

// UpdateName updates the user's name.
func (u *User) UpdateName(name string) error {
	if u.person.Name == name {
		return ErrUserUpdateFailed
	}

	u.register(&event.UserNameChanged{
		ID:   u.person.ID,
		Name: name,
	})

	return nil
}

// UpdateEmail updates the user's email.
func (u *User) UpdateEmail(email string) error {
	if u.email == email {
		return ErrUserUpdateFailed
	}

	u.register(&event.UserEmailChanged{
		ID:    u.person.ID,
		Email: email,
	})

	return nil
}

// Activate activates user.
func (u *User) Activate() error {
	if u.activated {
		return ErrUserUpdateFailed
	}

	u.register(&event.UserActivated{
		ID: u.person.ID,
	})

	return nil
}

// Deactivate deactivates user.
func (u *User) Deactivate() error {
	if !u.activated {
		return ErrUserUpdateFailed
	}

	u.register(&event.UserDeactivated{
		ID: u.person.ID,
	})

	return nil
}

// Apply applies user events to the user aggregate.
func (u *User) Apply(e event.Event, new bool) {
	switch ue := e.(type) {
	case *event.UserRegistered:
		u.person = &entity.Person{
			ID:   ue.ID,
			Name: ue.Name,
		}
		u.email = ue.Email
		u.activated = true

	case *event.UserNameChanged:
		u.person.Name = ue.Name

	case *event.UserEmailChanged:
		u.email = ue.Email

	case *event.UserDeactivated:
		u.activated = false

	case *event.UserActivated:
		u.activated = true
	}

	if !new {
		u.version++
	}
}

// Events returns the uncommitted events from the user aggregate.
func (u User) Events() []event.Event {
	return u.changes
}

// Version returns the last version of the aggregate before changes.
func (u User) Version() int {
	return u.version
}

func (u *User) register(event event.Event) {
	u.changes = append(u.changes, event)
	u.Apply(event, true)
}
