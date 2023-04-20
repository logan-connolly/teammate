package services

import (
	"errors"

	"git.sr.ht/~loges/teammate/internal/access/domain/model"
	"git.sr.ht/~loges/teammate/internal/access/domain/repository"
	"git.sr.ht/~loges/teammate/internal/access/infrastructure/memory"
	"git.sr.ht/~loges/teammate/internal/entity"
	"github.com/google/uuid"
)

var ErrInvalidRegistrationConfig = errors.New("services: invalid registration configuration")

// RegistrationConfigs defines the configurations to intialize the service with.
var RegistrationConfigs = []RegistrationConfiguration{
	WithMemoryRepositories(),
}

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// WithMemoryRepositories attaches in memory repostories to service.
func WithMemoryRepositories() RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.users = memory.NewMemoryUserRepository()
		return nil
	}
}

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	users repository.UserRepository
}

// NewRegistrationService accepts configs and returns a new service.
func NewRegistrationService() (*RegistrationService, error) {
	s := &RegistrationService{}

	for _, cfg := range RegistrationConfigs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

// RegisterUser registers a user if the email is not already registered.
func (s *RegistrationService) RegisterUser(name, email string) error {
	u, err := model.NewUser(&entity.Person{ID: uuid.New(), Name: name}, email)
	if err != nil {
		return err
	}

	if err = s.users.Add(u); err != nil {
		return err
	}

	return nil
}
