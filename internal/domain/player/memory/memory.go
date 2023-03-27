package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/domain/player"
	"github.com/logan-connolly/teammate/internal/entity"
)

// MemoryPlayerRepository is an in-memory player repository.
type MemoryPlayerRepository struct {
	players map[uuid.UUID][]player.Event
	sync.Mutex
}

// NewMemoryPlayerRepository intializes an in-memory player repository.
func NewMemoryPlayerRepository() *MemoryPlayerRepository {
	return &MemoryPlayerRepository{
		players: make(map[uuid.UUID][]player.Event),
	}
}

// Get retrieves a player by ID.
func (r *MemoryPlayerRepository) Get(p *entity.Person) (*player.Player, error) {
	if events, ok := r.players[p.ID]; ok {
		return player.NewPlayerFromEvents(events), nil
	}

	return &player.Player{}, player.ErrPlayerNotFound
}

// Add stores a new player in the repository.
func (r *MemoryPlayerRepository) Add(p *player.Player) error {
	if _, ok := r.players[p.GetID()]; ok {
		return player.ErrPlayerAlreadyExists
	}

	r.Lock()
	r.players[p.GetID()] = p.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to player in the repository.
func (r *MemoryPlayerRepository) Update(p *player.Player) error {
	storedEvents, ok := r.players[p.GetID()]
	if !ok {
		return player.ErrPlayerNotFound
	}

	newEvents := p.Events()
	if len(newEvents) == 0 {
		return player.ErrPlayerHasNoUpdates
	}

	r.Lock()
	r.players[p.GetID()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
