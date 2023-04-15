package access

import "github.com/logan-connolly/teammate/internal/access/application/services"

var memoryConfigs = []services.RegistrationConfiguration{
	services.WithMemoryUserRepository(),
}

// AccessApplication holds all services related to access management.
type AccessApplication struct {
	registrationService services.RegistrationService
}

// NewAccessApplication intitializes the access application.
func NewAccessApplication() (*AccessApplication, error) {
	rs, err := services.NewRegistrationService(memoryConfigs...)
	if err != nil {
		return &AccessApplication{}, services.ErrInvalidRegistrationConfig
	}

	return &AccessApplication{registrationService: *rs}, nil
}
