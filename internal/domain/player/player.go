package player

import (
	"errors"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
)

var (
	ErrInvalidPerson            = errors.New("player has to be a valid person")
	ErrPlayerAlreadyActivated   = errors.New("player is already activated")
	ErrPlayerAlreadyDeactivated = errors.New("player is inactive")
)

// Player is a aggregate that combines all entities needed to represent a player.
type Player struct {
	person    *entity.Person
	activated bool

	changes []Event
	version int
}

// NewPlayer is a factory to create a new Player aggregate.
func NewPlayer(name string) (*Player, error) {
	p := &Player{}

	if name == "" {
		return p, ErrInvalidPerson
	}

	person := &entity.Person{
		ID:   uuid.New(),
		Name: name,
	}

	p.register(&PlayerRegistered{
		ID:   person.ID,
		Name: person.Name,
	})

	return p, nil
}

// NewFromEvents is a helper method that creates a new player
// from a series of events.
func NewPlayerFromEvents(events []Event) *Player {
	p := &Player{}

	for _, event := range events {
		p.Apply(event, false)
	}

	return p
}

// GetID returns the player root entity GetID.
func (c *Player) GetID() uuid.UUID {
	return c.person.ID
}

// GetName returns the name of the player.
func (p *Player) GetName() string {
	return p.person.Name
}

// IsActivated returns whether the player is activated.
func (p *Player) IsActivated() bool {
	return p.activated
}

// Activate activates player.
func (p *Player) Activate() error {
	if p.activated {
		return ErrPlayerAlreadyActivated
	}

	p.register(&PlayerActivated{
		ID: p.person.ID,
	})

	return nil
}

// Deactivate deactivates player.
func (p *Player) Deactivate() error {
	if !p.activated {
		return ErrPlayerAlreadyDeactivated
	}

	p.register(&PlayerDeactivated{
		ID: p.person.ID,
	})

	return nil
}

// Apply applies player events to the player aggregate.
func (p *Player) Apply(event Event, new bool) {
	switch e := event.(type) {
	case *PlayerRegistered:
		p.person = &entity.Person{
			ID:   e.ID,
			Name: e.Name,
		}
		p.activated = true

	case *PlayerDeactivated:
		p.activated = false

	case *PlayerActivated:
		p.activated = true
	}

	if !new {
		p.version++
	}
}

// Events returns the uncommitted events from the player aggregate.
func (p Player) Events() []Event {
	return p.changes
}

// Version returns the last version of the aggregate before changes.
func (p Player) Version() int {
	return p.version
}

func (p *Player) register(event Event) {
	p.changes = append(p.changes, event)
	p.Apply(event, true)
}
