package access

import "git.sr.ht/~loges/teammate/internal/access/application/services"

// AccessApplication holds all services related to access management.
type AccessApplication struct {
	registrationService *services.RegistrationService
}

// NewAccessApplication intitializes the access application.
func NewAccessApplication() (*AccessApplication, error) {
	rs, err := services.NewRegistrationService()
	if err != nil {
		return &AccessApplication{}, services.ErrInvalidRegistrationConfig
	}

	return &AccessApplication{registrationService: rs}, nil
}

// GetRosterService returns the roster service from the app.
func (a *AccessApplication) GetRegistrationService() *services.RegistrationService {
	return a.registrationService
}
