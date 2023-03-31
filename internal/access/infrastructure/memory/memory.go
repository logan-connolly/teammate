package memory

import (
	"sync"

	"github.com/google/uuid"

	"github.com/logan-connolly/teammate/internal/access/domain/event"
	"github.com/logan-connolly/teammate/internal/access/domain/model"
	"github.com/logan-connolly/teammate/internal/access/domain/repository"
	"github.com/logan-connolly/teammate/internal/entity"
)

// MemoryUserRepository is an in-memory user repository.
type MemoryAccessRepository struct {
	users map[uuid.UUID][]event.Event
	sync.Mutex
}

// NewMemoryUserRepository intializes an in-memory user repository.
func NewMemoryAccessRepository() *MemoryAccessRepository {
	return &MemoryAccessRepository{
		users: make(map[uuid.UUID][]event.Event),
	}
}

// Get retrieves a user by ID.
func (r *MemoryAccessRepository) Get(p *entity.Person) (*model.User, error) {
	if events, ok := r.users[p.ID]; ok {
		return model.NewUserFromEvents(events), nil
	}

	return &model.User{}, repository.ErrUserNotFound
}

// Add stores a new user in the repository.
func (r *MemoryAccessRepository) Add(p *model.User) error {
	if _, ok := r.users[p.GetID()]; ok {
		return repository.ErrUserAlreadyExists
	}

	r.Lock()
	r.users[p.GetID()] = p.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to user in the repository.
func (r *MemoryAccessRepository) Update(p *model.User) error {
	storedEvents, ok := r.users[p.GetID()]
	if !ok {
		return repository.ErrUserNotFound
	}

	newEvents := p.Events()
	if len(newEvents) == 0 {
		return repository.ErrUserHasNoUpdates
	}

	r.Lock()
	r.users[p.GetID()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
