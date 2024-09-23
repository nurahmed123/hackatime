package imports

import (
	"strings"
	"time"

	"github.com/hackclub/hackatime/config"
	"github.com/hackclub/hackatime/models"
)

type WakatimeImporter struct {
	apiKey      string
	forceLegacy bool
}

func NewWakatimeImporter(apiKey string, forceLegacy bool) *WakatimeImporter {
	return &WakatimeImporter{apiKey: apiKey, forceLegacy: forceLegacy}
}

func (w *WakatimeImporter) Import(user *models.User, minFrom time.Time, maxTo time.Time) (<-chan *models.Heartbeat, error) {
	if strings.Contains(user.WakaTimeURL(config.WakatimeApiUrl), "wakatime.com") && !w.forceLegacy {
		return NewWakatimeDumpImporter(w.apiKey).Import(user, minFrom, maxTo)
	}
	return NewWakatimeHeartbeatImporter(w.apiKey).Import(user, minFrom, maxTo)
}

func (w *WakatimeImporter) ImportAll(user *models.User) (<-chan *models.Heartbeat, error) {
	if strings.Contains(user.WakaTimeURL(config.WakatimeApiUrl), "wakatime.com") && !w.forceLegacy {
		return NewWakatimeDumpImporter(w.apiKey).ImportAll(user)
	}
	return NewWakatimeHeartbeatImporter(w.apiKey).ImportAll(user)
}
