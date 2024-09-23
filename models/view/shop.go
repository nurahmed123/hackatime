package view

import (
	"github.com/hackclub/hackatime/models"
	"github.com/hackclub/hackatime/utils"
)

type ShopViewModel struct {
	SharedLoggedInViewModel
	Products   []*models.Product
	PageParams *utils.PageParams
}

func (s *ShopViewModel) LangIcon(lang string) string {
	return GetLanguageIcon(lang)
}

func (s *ShopViewModel) WithSuccess(m string) *ShopViewModel {
	s.SetSuccess(m)
	return s
}

func (s *ShopViewModel) WithError(m string) *ShopViewModel {
	s.SetError(m)
	return s
}
