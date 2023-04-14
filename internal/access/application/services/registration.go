package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/logan-connolly/teammate/internal/access/domain/model"
	"github.com/logan-connolly/teammate/internal/access/domain/repository"
	"github.com/logan-connolly/teammate/internal/access/infrastructure/memory"
	"github.com/logan-connolly/teammate/internal/entity"
)

var ErrInvalidRegistrationConfig = errors.New("services: invalid registration configuration")

// RegistrationConfiguration is a function that modifies the service.
type RegistrationConfiguration func(s *RegistrationService) error

// RegistrationService is a implementation of the RegistrationService.
type RegistrationService struct {
	users repository.UserRepository
}

// NewRegistrationService accepts configs and returns a new service.
func NewRegistrationService(cfgs ...RegistrationConfiguration) (*RegistrationService, error) {
	s := &RegistrationService{}

	for _, cfg := range cfgs {
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

// WithUserRepository applies a given user repository to the service.
func WithUserRepository(r repository.UserRepository) RegistrationConfiguration {
	return func(s *RegistrationService) error {
		s.users = r
		return nil
	}
}

// WithMemoryUserRepository applies a memory user repository to the service.
func WithMemoryUserRepository() RegistrationConfiguration {
	r := memory.NewMemoryUserRepository()
	return WithUserRepository(r)
}
