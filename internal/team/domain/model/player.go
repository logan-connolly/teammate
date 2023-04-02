package model

import (
	"errors"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/entity"
	"github.com/logan-connolly/teammate/internal/team/domain/event"
)

var (
	ErrInvalidPerson            = errors.New("model: player has to be a valid person")
	ErrPlayerAlreadyActivated   = errors.New("model: player is already activated")
	ErrPlayerAlreadyDeactivated = errors.New("model: player is already inactive")
	ErrTeamAlreadyAssigned      = errors.New("model: team already assigned to player")
	ErrTeamNotAssignedToPlayer  = errors.New("model: team not assigned to player")
)

// TeamMapping stores the mapping of registered team ids to group entities.
type TeamMapping map[uuid.UUID]*entity.Group

// Player is a aggregate that combines all entities needed to represent a player.
type Player struct {
	person    *entity.Person
	activated bool
	teams     TeamMapping

	changes []event.Event
	version int
}

// NewPlayer is a factory to create a new Player aggregate.
func NewPlayer(p *entity.Person) (*Player, error) {
	player := &Player{}

	if p.Name == "" {
		return player, ErrInvalidPerson
	}

	player.register(&event.PlayerCreated{
		ID:   p.ID,
		Name: p.Name,
	})

	return player, nil
}

// NewFromEvents is a helper method that creates a new player
// from a series of events.
func NewPlayerFromEvents(events []event.Event) *Player {
	p := &Player{}

	for _, event := range events {
		p.Apply(event, false)
	}

	return p
}

// GetID returns the player root entity GetID.
func (p *Player) GetID() uuid.UUID {
	return p.person.ID
}

// GetName returns the name of the player.
func (p *Player) GetName() string {
	return p.person.Name
}

func (p *Player) GetTeams() (teams []*entity.Group) {
	for _, team := range p.teams {
		teams = append(teams, team)
	}
	return teams
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

	p.register(&event.PlayerActivated{
		ID: p.person.ID,
	})

	return nil
}

// Deactivate deactivates player.
func (p *Player) Deactivate() error {
	if !p.activated {
		return ErrPlayerAlreadyDeactivated
	}

	p.register(&event.PlayerDeactivated{
		ID: p.person.ID,
	})

	return nil
}

// AssignTeam assigns team to player.
func (p *Player) AssignTeam(t *Team) error {
	if _, ok := p.teams[t.group.ID]; ok {
		return ErrTeamAlreadyAssigned
	}

	p.register(&event.TeamAssignedToPlayer{
		ID:       p.person.ID,
		TeamId:   t.group.ID,
		TeamName: t.group.Name,
	})

	return nil
}

// UnassignTeam unassigns team from player.
func (p *Player) UnassignTeam(t *Team) error {
	if _, ok := p.teams[t.group.ID]; !ok {
		return ErrTeamNotAssignedToPlayer
	}

	p.register(&event.TeamUnassignedFromPlayer{
		ID:       p.person.ID,
		TeamId:   t.group.ID,
		TeamName: t.group.Name,
	})

	return nil
}

// Apply applies player events to the player aggregate.
func (p *Player) Apply(e event.Event, new bool) {
	switch pe := e.(type) {
	case *event.PlayerCreated:
		p.person = &entity.Person{
			ID:   pe.ID,
			Name: pe.Name,
		}
		p.activated = true
		p.teams = make(map[uuid.UUID]*entity.Group)

	case *event.PlayerDeactivated:
		p.activated = false

	case *event.PlayerActivated:
		p.activated = true

	case *event.TeamAssignedToPlayer:
		p.teams[pe.TeamId] = &entity.Group{ID: pe.TeamId, Name: pe.TeamName}

	case *event.TeamUnassignedFromPlayer:
		delete(p.teams, pe.TeamId)
	}

	if !new {
		p.version++
	}
}

// Events returns the uncommitted events from the player aggregate.
func (p Player) Events() []event.Event {
	return p.changes
}

// Version returns the last version of the aggregate before changes.
func (p Player) Version() int {
	return p.version
}

func (p *Player) register(event event.Event) {
	p.changes = append(p.changes, event)
	p.Apply(event, true)
}
