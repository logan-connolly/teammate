package access

import "github.com/logan-connolly/teammate/internal/access/services"

var memoryConfigs = []services.RegistrationConfiguration{
	services.WithMemoryUserRepository(),
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
