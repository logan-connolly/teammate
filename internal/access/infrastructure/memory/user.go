package memory

import (
	"sync"

	"github.com/logan-connolly/teammate/internal/access/domain/event"
	"github.com/logan-connolly/teammate/internal/access/domain/model"
	"github.com/logan-connolly/teammate/internal/access/domain/repository"
)

// MemoryUserRepository is an in-memory user repository.
type MemoryUserRepository struct {
	users map[string][]event.Event
	sync.Mutex
}

// NewMemoryUserRepository intializes an in-memory user repository.
func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{
		users: make(map[string][]event.Event),
	}
}

// Get retrieves a user by ID.
func (r *MemoryUserRepository) GetByEmail(email string) (*model.User, error) {
	if events, ok := r.users[email]; ok {
		return model.NewUserFromEvents(events), nil
	}

	return &model.User{}, repository.ErrUserNotFound
}

// Add stores a new user in the repository.
func (r *MemoryUserRepository) Add(p *model.User) error {
	if _, ok := r.users[p.GetEmail()]; ok {
		return repository.ErrUserAlreadyExists
	}

	r.Lock()
	r.users[p.GetEmail()] = p.Events()
	defer r.Unlock()

	return nil
}

// Update appends changes to user in the repository.
func (r *MemoryUserRepository) Update(p *model.User) error {
	storedEvents, ok := r.users[p.GetEmail()]
	if !ok {
		return repository.ErrUserNotFound
	}

	newEvents := p.Events()
	if len(newEvents) == 0 {
		return repository.ErrUserHasNoUpdates
	}

	r.Lock()
	r.users[p.GetEmail()] = append(storedEvents, newEvents...)
	defer r.Unlock()

	return nil
}
