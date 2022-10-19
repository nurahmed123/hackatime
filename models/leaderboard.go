package models

import (
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"strings"
	"time"
)

type LeaderboardItem struct {
	ID        uint          `json:"-" gorm:"primary_key; size:32"`
	User      *User         `json:"-" gorm:"not null; constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UserID    string        `json:"user_id" gorm:"not null; index:idx_leaderboard_user"`
	Rank      uint          `json:"rank" gorm:"-:migration"`
	Interval  string        `json:"interval" gorm:"not null; size:32; index:idx_leaderboard_combined"`
	By        *uint8        `json:"aggregated_by" gorm:"index:idx_leaderboard_combined"` // pointer because nullable
	Total     time.Duration `json:"total" gorm:"not null" swaggertype:"primitive,integer"`
	Key       *string       `json:"key" gorm:"size:255"` // pointer because nullable
	CreatedAt CustomTime    `gorm:"type:timestamp; default:CURRENT_TIMESTAMP" swaggertype:"string" format:"date" example:"2006-01-02 15:04:05.000"`
}

func (l1 *LeaderboardItem) Equals(l2 *LeaderboardItem) bool {
	return l1.ID == l2.ID
}

type Leaderboard []*LeaderboardItem

func (l *Leaderboard) Add(item *LeaderboardItem) {
	if _, found := slice.Find[*LeaderboardItem](*l, func(i int, item2 *LeaderboardItem) bool {
		return item.Equals(item2)
	}); !found {
		*l = append(*l, item)
	}
}

func (l *Leaderboard) AddMany(items []*LeaderboardItem) {
	for _, item := range items {
		l.Add(item)
	}
}

func (l Leaderboard) UserIDs() []string {
	return slice.Unique[string](slice.Map[*LeaderboardItem, string](l, func(i int, item *LeaderboardItem) string {
		return item.UserID
	}))
}

func (l Leaderboard) HasUser(userId string) bool {
	return slice.Contain(l.UserIDs(), userId)
}

func (l Leaderboard) TopByKey(by uint8, key string) Leaderboard {
	return slice.Filter[*LeaderboardItem](l, func(i int, item *LeaderboardItem) bool {
		return item.By != nil && *item.By == by && item.Key != nil && strings.ToLower(*item.Key) == strings.ToLower(key)
	})
}

func (l Leaderboard) TopKeys(by uint8) []string {
	type keyTotal struct {
		Key   string
		Total time.Duration
	}

	totalsMapped := make(map[string]*keyTotal, len(l))

	for _, item := range l {
		if item.Key == nil || item.By == nil || *item.By != by {
			continue
		}
		key := strings.ToLower(*item.Key)
		if _, ok := totalsMapped[key]; !ok {
			totalsMapped[key] = &keyTotal{Key: *item.Key, Total: 0}
		}
		totalsMapped[key].Total += item.Total
	}

	totals := slice.Map[*keyTotal, keyTotal](maputil.Values[string, *keyTotal](totalsMapped), func(i int, item *keyTotal) keyTotal {
		return *item
	})
	if err := slice.SortByField(totals, "Total", "desc"); err != nil {
		return []string{} // TODO
	}

	return slice.Map[keyTotal, string](totals, func(i int, item keyTotal) string {
		return item.Key
	})
}

func (l Leaderboard) TopKeysByUser(by uint8, userId string) []string {
	return Leaderboard(slice.Filter[*LeaderboardItem](l, func(i int, item *LeaderboardItem) bool {
		return item.UserID == userId
	})).TopKeys(by)
}

func (l Leaderboard) LastUpdate() time.Time {
	lastUpdate := time.Time{}
	for _, item := range l {
		if item.CreatedAt.T().After(lastUpdate) {
			lastUpdate = item.CreatedAt.T()
		}
	}
	return lastUpdate
}
