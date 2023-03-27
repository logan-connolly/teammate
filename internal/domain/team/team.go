package team

import (
	"errors"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
)

var (
	ErrInvalidGroup           = errors.New("team has to be a valid group")
	ErrTeamAlreadyActivated   = errors.New("team is already activated")
	ErrTeamAlreadyDeactivated = errors.New("team is already inactive")
)

// Team is a aggregate that combines all entities needed to represent a team.
type Team struct {
	group     *entity.Group
	activated bool

	changes []Event
	version int
}

// NewTeam is a factory to create a new Team aggregate.
func NewTeam(name string) (*Team, error) {
	t := &Team{}

	if name == "" {
		return t, ErrInvalidGroup
	}

	group := &entity.Group{
		ID:   uuid.New(),
		Name: name,
	}

	t.register(&TeamRegistered{
		ID:   group.ID,
		Name: group.Name,
	})

	return t, nil
}

// NewFromEvents is a helper method that creates a new team
// from a series of events.
func NewTeamFromEvents(events []Event) *Team {
	t := &Team{}

	for _, event := range events {
		t.Apply(event, false)
	}

	return t
}

// GetID returns the team root entity GetID.
func (t *Team) GetID() uuid.UUID {
	return t.group.ID
}

// GetName returns the name of the team.
func (t *Team) GetName() string {
	return t.group.Name
}

// IsActivated returns whether the team is activated.
func (t *Team) IsActivated() bool {
	return t.activated
}

// Activate activates team.
func (t *Team) Activate() error {
	if t.activated {
		return ErrTeamAlreadyActivated
	}

	t.register(&TeamActivated{
		ID: t.group.ID,
	})

	return nil
}

// Deactivate deactivates team.
func (t *Team) Deactivate() error {
	if !t.activated {
		return ErrTeamAlreadyDeactivated
	}

	t.register(&TeamDeactivated{
		ID: t.group.ID,
	})

	return nil
}

// Apply applies team events to the team aggregate.
func (t *Team) Apply(event Event, new bool) {
	switch e := event.(type) {
	case *TeamRegistered:
		t.group = &entity.Group{
			ID:   e.ID,
			Name: e.Name,
		}
		t.activated = true

	case *TeamDeactivated:
		t.activated = false

	case *TeamActivated:
		t.activated = true
	}

	if !new {
		t.version++
	}
}

// Events returns the uncommitted events from the team aggregate.
func (t Team) Events() []Event {
	return t.changes
}

// Version returns the last version of the aggregate before changes.
func (t Team) Version() int {
	return t.version
}

func (t *Team) register(event Event) {
	t.changes = append(t.changes, event)
	t.Apply(event, true)
}
