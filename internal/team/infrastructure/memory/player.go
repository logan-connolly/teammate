package memory

import (
	"sync"

	"git.sr.ht/~loges/teammate/internal/entity"
	"git.sr.ht/~loges/teammate/internal/team/domain/event"
	"git.sr.ht/~loges/teammate/internal/team/domain/model"
	"git.sr.ht/~loges/teammate/internal/team/domain/repository"
	"github.com/google/uuid"
)

// MemoryPlayerRepository is an in-memory player repository.
type MemoryPlayerRepository struct {
	players map[uuid.UUID][]event.Event
	sync.Mutex
}

// NewMemoryPlayerRepository intializes an in-memory player repository.
func NewMemoryPlayerRepository() *MemoryPlayerRepository {
	return &MemoryPlayerRepository{
		players: make(map[uuid.UUID][]event.Event),
	}
}

// Get retrieves a player by ID.
func (r *MemoryPlayerRepository) Get(p *entity.Person) (*model.Player, error) {
	if events, ok := r.players[p.ID]; ok {
		return model.NewPlayerFromEvents(events), nil
	}

	return &model.Player{}, repository.ErrPlayerNotFound
}

// GetTeams retrieves teams assigned to players.
func (r *MemoryPlayerRepository) GetTeams(p *entity.Person) ([]*entity.Group, error) {
	events, ok := r.players[p.ID]
	if !ok {
		return []*entity.Group{}, repository.ErrPlayerNotFound
	}
	return model.NewPlayerFromEvents(events).GetTeams(), nil
}

// Add stores a new player in the repository.
func (r *MemoryPlayerRepository) Add(p *model.Player) error {
	if _, ok := r.players[p.GetID()]; ok {
		return repository.ErrPlayerAlreadyExists
	}

	r.Lock()
	r.players[p.GetID()] = p.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to player in the repository.
func (r *MemoryPlayerRepository) Update(p *model.Player) error {
	storedEvents, ok := r.players[p.GetID()]
	if !ok {
		return repository.ErrPlayerNotFound
	}

	newEvents := p.Events()
	if len(newEvents) == 0 {
		return repository.ErrPlayerHasNoUpdates
	}

	r.Lock()
	r.players[p.GetID()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
