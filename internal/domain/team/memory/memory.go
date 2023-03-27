package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/domain/team"
)

// MemoryTeamRepository is an in-memory team repository.
type MemoryTeamRepository struct {
	teams map[uuid.UUID][]team.Event
	sync.Mutex
}

// NewMemoryTeamRepository intializes an in-memory team repository.
func NewMemoryTeamRepository() *MemoryTeamRepository {
	return &MemoryTeamRepository{
		teams: make(map[uuid.UUID][]team.Event),
	}
}

// Get retrieves a team by ID.
func (r *MemoryTeamRepository) Get(id uuid.UUID) (*team.Team, error) {
	if events, ok := r.teams[id]; ok {
		return team.NewTeamFromEvents(events), nil
	}

	return &team.Team{}, team.ErrTeamNotFound
}

// Add stores a new team in the repository.
func (r *MemoryTeamRepository) Add(t *team.Team) error {
	if _, ok := r.teams[t.GetID()]; ok {
		return team.ErrTeamAlreadyExists
	}

	r.Lock()
	r.teams[t.GetID()] = t.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to team in the repository.
func (r *MemoryTeamRepository) Update(t *team.Team) error {
	storedEvents, ok := r.teams[t.GetID()]
	if !ok {
		return team.ErrTeamNotFound
	}

	newEvents := t.Events()
	if len(newEvents) == 0 {
		return team.ErrTeamHasNoUpdates
	}

	r.Lock()
	r.teams[t.GetID()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
