package memory

import (
	"sync"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/event"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
	"git.sr.ht/~loges/teammate/internal/team/domain/repository"
	"github.com/google/uuid"
)

// MemoryTeamRepository is an in-memory team repository.
type MemoryTeamRepository struct {
	teams map[uuid.UUID][]event.Event
	sync.Mutex
}

// NewMemoryTeamRepository intializes an in-memory team repository.
func NewMemoryTeamRepository() *MemoryTeamRepository {
	return &MemoryTeamRepository{
		teams: make(map[uuid.UUID][]event.Event),
	}
}

// Get retrieves a team by ID.
func (r *MemoryTeamRepository) Get(g *entity.Group) (*model.Team, error) {
	if events, ok := r.teams[g.ID]; ok {
		return model.NewTeamFromEvents(events), nil
	}

	return &model.Team{}, repository.ErrTeamNotFound
}

// GetPlayers retrieves a team by ID.
func (r *MemoryTeamRepository) GetPlayers(p *entity.Group) ([]*entity.Person, error) {
	events, ok := r.teams[p.ID]
	if !ok {
		return []*entity.Person{}, repository.ErrTeamNotFound
	}
	return model.NewTeamFromEvents(events).GetPlayers(), nil
}

// Add stores a new team in the repository.
func (r *MemoryTeamRepository) Add(t *model.Team) error {
	if _, ok := r.teams[t.GetID()]; ok {
		return repository.ErrTeamAlreadyExists
	}

	r.Lock()
	r.teams[t.GetID()] = t.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to team in the repository.
func (r *MemoryTeamRepository) Update(t *model.Team) error {
	storedEvents, ok := r.teams[t.GetID()]
	if !ok {
		return repository.ErrTeamNotFound
	}

	newEvents := t.Events()
	if len(newEvents) == 0 {
		return repository.ErrTeamHasNoUpdates
	}

	r.Lock()
	r.teams[t.GetID()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
