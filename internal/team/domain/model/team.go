package model

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/event"
	"github.com/google/uuid"
)

var (
	ErrInvalidGroup     = errors.New("model: team has to be a valid group")
	ErrTeamUpdateFailed = errors.New("model: team update failed")
)

// Team is a aggregate that combines all entities needed to represent a team.
type Team struct {
	group     *entity.Group
	activated bool
	players   map[uuid.UUID]*entity.Person

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
		return ErrTeamUpdateFailed
	}

	t.register(&event.TeamActivated{
		ID: t.group.ID,
	})

	return nil
}

// Deactivate deactivates team.
func (t *Team) Deactivate() error {
	if !t.activated {
		return ErrTeamUpdateFailed
	}

	t.register(&event.TeamDeactivated{
		ID: t.group.ID,
	})

	return nil
}

// AssignPlayer assigns player to team.
func (t *Team) AssignPlayer(p *Player) error {
	if _, ok := t.players[p.person.ID]; ok {
		return ErrTeamUpdateFailed
	}

	t.register(&event.PlayerAssignedToTeam{
		ID:         t.group.ID,
		PlayerId:   p.person.ID,
		PlayerName: p.person.Name,
	})

	return nil
}

// UassignPlayer assigns player from team.
func (t *Team) UnassignPlayer(p *Player) error {
	if _, ok := t.players[p.person.ID]; !ok {
		return ErrTeamUpdateFailed
	}

	t.register(&event.PlayerUnassignedFromTeam{
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

	case *event.PlayerUnassignedFromTeam:
		delete(t.players, te.PlayerId)
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
