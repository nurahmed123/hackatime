package mail

import (
	"log/slog"

	"github.com/hackclub/hackatime/models"
)

type NoopSendingService struct{}

func (n *NoopSendingService) Send(mail *models.Mail) error {
	slog.Info("noop mail service doing nothing instead of sending password reset mail", "to", mail.To.Strings())
	return nil
}
