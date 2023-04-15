package application

import (
	"github.com/logan-connolly/teammate/internal/team/application/services"
)

// TeamApplication holds all services related to team management.
type TeamApplication struct {
	rosterService *services.RosterService
}

// NewTeamApplication intitializes the team application.
func NewTeamApplication() (*TeamApplication, error) {
	rs, err := services.NewRosterService()
	if err != nil {
		return &TeamApplication{}, services.ErrInvalidRosterConfig
	}

	return &TeamApplication{rosterService: rs}, nil
}

// GetRosterService returns the roster service from the app.
func (a *TeamApplication) GetRosterService() *services.RosterService {
	return a.rosterService
}
