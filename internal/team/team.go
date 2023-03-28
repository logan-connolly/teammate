package team

import (
	"github.com/logan-connolly/teammate/internal/team/services"
)

var memoryConfigs = []services.RegistrationConfiguration{
	services.WithMemoryPlayerRepository(),
	services.WithMemoryTeamRepository(),
}

type Application struct {
	registrationService services.RegistrationService
}

func NewMemoryApplication() (*Application, error) {
	rs, err := services.NewRegistrationService(memoryConfigs...)
	if err != nil {
		return &Application{}, services.ErrInvalidRegistrationConfig
	}

	return &Application{
		registrationService: *rs,
	}, nil
}
