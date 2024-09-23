package imports

import (
	"time"

	"github.com/hackclub/hackatime/models"
)

type DataImporter interface {
	Import(*models.User, time.Time, time.Time) (<-chan *models.Heartbeat, error)
	ImportAll(*models.User) (<-chan *models.Heartbeat, error)
}
