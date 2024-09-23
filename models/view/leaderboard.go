package view

import (
	"time"

	"github.com/hackclub/hackatime/models"
	"github.com/hackclub/hackatime/utils"
)

type LeaderboardViewModel struct {
	SharedLoggedInViewModel
	By            string
	Key           string
	Items         []*models.LeaderboardItemRanked
	TopKeys       []string
	UserLanguages map[string][]string
	IntervalLabel string
	PageParams    *utils.PageParams
}

func (s *LeaderboardViewModel) WithSuccess(m string) *LeaderboardViewModel {
	s.SetSuccess(m)
	return s
}

func (s *LeaderboardViewModel) WithError(m string) *LeaderboardViewModel {
	s.SetError(m)
	return s
}

func (s *LeaderboardViewModel) ColorModifier(item *models.LeaderboardItemRanked, principal *models.User) string {
	if principal != nil && item.UserID == principal.ID {
		return "border-accent-primary dark:border-accent-dark-primary border-3"
	}
	return ""
}

func (s *LeaderboardViewModel) LangIcon(lang string) string {
	return GetLanguageIcon(lang)
}

func (s *LeaderboardViewModel) LastUpdate() time.Time {
	return models.Leaderboard(s.Items).LastUpdate()
}
