package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
)

var (
	ErrInvalidGroup           = errors.New("model: team has to be a valid group")
	ErrTeamAlreadyActivated   = errors.New("model: team is already activated")
	ErrTeamAlreadyDeactivated = errors.New("model: team is already inactive")
	ErrPlayerAlreadyAssigned  = errors.New("model: player already assigned to team")
)

// TeamMapping stores the mapping of registered team ids to group entities.
type PlayerMapping map[uuid.UUID]*entity.Person

// Team is a aggregate that combines all entities needed to represent a team.
type Team struct {
	group     *entity.Group
	activated bool
	players   PlayerMapping

	changes []event.Event
	version int
}

// NewTeam is a factory to create a new Team aggregate.
func NewTeam(g *entity.Group) (*Team, error) {
	t := &Team{}

	if g.Name == "" {
		return t, ErrInvalidGroup
	}

	t.register(&event.TeamCreated{
		ID:   g.ID,
		Name: g.Name,
	})

	return t, nil
}

// NewFromEvents is a helper method that creates a new team
// from a series of events.
func NewTeamFromEvents(events []event.Event) *Team {
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

func (t *Team) GetPlayers() (players []*entity.Person) {
	for _, player := range t.players {
		players = append(players, player)
	}
	return players
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

	t.register(&event.TeamActivated{
		ID: t.group.ID,
	})

	return nil
}

// Deactivate deactivates team.
func (t *Team) Deactivate() error {
	if !t.activated {
		return ErrTeamAlreadyDeactivated
	}

	t.register(&event.TeamDeactivated{
		ID: t.group.ID,
	})

	return nil
}

// AssignPlayer assigns player to team.
func (t *Team) AssignPlayer(p *Player) error {
	if _, ok := t.players[p.person.ID]; ok {
		return ErrPlayerAlreadyAssigned
	}

	t.register(&event.PlayerAssignedToTeam{
		ID:         t.group.ID,
		PlayerId:   p.person.ID,
		PlayerName: p.person.Name,
	})

	return nil
}

// Apply applies team events to the team aggregate.
func (t *Team) Apply(e event.Event, new bool) {
	switch te := e.(type) {
	case *event.TeamCreated:
		t.group = &entity.Group{
			ID:   te.ID,
			Name: te.Name,
		}
		t.activated = true
		t.players = make(map[uuid.UUID]*entity.Person)

	case *event.TeamDeactivated:
		t.activated = false

	case *event.TeamActivated:
		t.activated = true

	case *event.PlayerAssignedToTeam:
		t.players[te.PlayerId] = &entity.Person{ID: te.PlayerId, Name: te.PlayerName}
	}

	if !new {
		t.version++
	}
}

// Events returns the uncommitted events from the team aggregate.
func (t Team) Events() []event.Event {
	return t.changes
}

// Version returns the last version of the aggregate before changes.
func (t Team) Version() int {
	return t.version
}

func (t *Team) register(event event.Event) {
	t.changes = append(t.changes, event)
	t.Apply(event, true)
}
