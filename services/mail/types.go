package mail

import "github.com/hackclub/hackatime/models"

type WelcomeTplData struct {
	PublicUrl string
	Email     string
	Name      string
	Id        string
}

type PasswordResetTplData struct {
	ResetLink string
}

type ImportNotificationTplData struct {
	PublicUrl     string
	Duration      string
	NumHeartbeats int
}

type WakatimeFailureNotificationNotificationTplData struct {
	PublicUrl   string
	NumFailures int
}

type ReportTplData struct {
	Report *models.Report
}

type SubscriptionNotificationTplData struct {
	PublicUrl           string
	HasExpired          bool
	DataRetentionMonths int
}
