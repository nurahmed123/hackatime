package view

import (
	conf "github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
)

type BasicViewModel interface {
	SetError(string)
	SetSuccess(string)
}

type Messages struct {
	Success string
	Error   string
}

type SharedViewModel struct {
	Messages
	LeaderboardEnabled bool
	ShopEnabled        bool
	InvitesEnabled     bool
}

type SharedLoggedInViewModel struct {
	SharedViewModel
	User   *models.User
	ApiKey string
}

func NewSharedViewModel(c *conf.Config, messages *Messages) SharedViewModel {
	vm := SharedViewModel{
		LeaderboardEnabled: c.App.LeaderboardEnabled,
		ShopEnabled:        c.Shop.Enabled,
		InvitesEnabled:     c.Security.InviteCodes,
	}
	if messages != nil {
		vm.Messages = *messages
	}
	return vm
}

func (m *Messages) SetError(message string) {
	m.Error = message
}

func (m *Messages) SetSuccess(message string) {
	m.Success = message
}
