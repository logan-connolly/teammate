package team

import (
	"github.com/logan-connolly/teammate/internal/team/services"
)

type Application struct {
	registrationService services.RegistrationService
}

func (a *Application) GetRegistrationService() services.RegistrationService {
	return a.registrationService
}

func NewMemoryApplication() (*Application, error) {
	rs, err := services.NewRegistrationService(
		services.WithMemoryPlayerRepository(),
		services.WithMemoryTeamRepository(),
	)
	if err != nil {
		return &Application{}, services.ErrInvalidRegistrationConfig
	}

	return &Application{
		registrationService: *rs,
	}, nil
}
