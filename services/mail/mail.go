package mail

import (
	"bytes"
	"fmt"
	"time"

	"github.com/hackclub/hackatime/helpers"
	"github.com/hackclub/hackatime/models"
	"github.com/hackclub/hackatime/routes"
	"github.com/hackclub/hackatime/services"
	"github.com/hackclub/hackatime/utils"
	"github.com/hackclub/hackatime/views/mail"

	conf "github.com/hackclub/hackatime/config"
)

const (
	tplNameWelcome                     = "welcome"
	tplNamePasswordReset               = "reset_password"
	tplNameImportNotification          = "import_finished"
	tplNameWakatimeFailureNotification = "wakatime_connection_failure"
	tplNameReport                      = "report"
	tplNameSubscriptionNotification    = "subscription_expiring"
	subjectWelcome                     = "Hackatime - Welcome!"
	subjectPasswordReset               = "Hackatime - Password Reset"
	subjectImportNotification          = "Hackatime - Data Import Finished"
	subjectWakatimeFailureNotification = "Hackatime - WakaTime Connection Failure"
	subjectReport                      = "Hackatime - Report from %s"
	subjectSubscriptionNotification    = "Hackatime - Subscription expiring / expired"
)

type SendingService interface {
	Send(*models.Mail) error
}

type MailService struct {
	config         *conf.Config
	sendingService SendingService
	templates      utils.TemplateMap
}

func NewMailService() services.IMailService {
	config := conf.Get()

	var sendingService SendingService
	sendingService = &NoopSendingService{}

	if config.Mail.Enabled {
		if config.Mail.Provider == conf.MailProviderSmtp {
			sendingService = NewSMTPSendingService(config.Mail.Smtp)
		}
	}

	// Use local file system when in 'dev' environment, go embed file system otherwise
	templateFs := conf.ChooseFS("views/mail", mail.TemplateFiles)
	templates, err := utils.LoadTemplates(templateFs, routes.DefaultTemplateFuncs())
	if err != nil {
		panic(err)
	}

	return &MailService{sendingService: sendingService, config: config, templates: templates}
}

func (m *MailService) SendWelcome(recipient *models.User) error {
	tpl, err := m.getWelcomeTemplate(WelcomeTplData{PublicUrl: m.config.Server.PublicUrl, Name: recipient.Name, Email: recipient.Email, Id: recipient.ID})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: subjectWelcome,
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) SendPasswordReset(recipient *models.User, resetLink string) error {
	tpl, err := m.getPasswordResetTemplate(PasswordResetTplData{ResetLink: resetLink})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: subjectPasswordReset,
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) SendWakatimeFailureNotification(recipient *models.User, numFailures int) error {
	tpl, err := m.getWakatimeFailureNotificationTemplate(WakatimeFailureNotificationNotificationTplData{
		PublicUrl:   m.config.Server.PublicUrl,
		NumFailures: numFailures,
	})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: subjectWakatimeFailureNotification,
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) SendImportNotification(recipient *models.User, duration time.Duration, numHeartbeats int) error {
	tpl, err := m.getImportNotificationTemplate(ImportNotificationTplData{
		PublicUrl:     m.config.Server.PublicUrl,
		Duration:      fmt.Sprintf("%.0f seconds", duration.Seconds()),
		NumHeartbeats: numHeartbeats,
	})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: subjectImportNotification,
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) SendReport(recipient *models.User, report *models.Report) error {
	tpl, err := m.getReportTemplate(ReportTplData{report})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: fmt.Sprintf(subjectReport, helpers.FormatDateHuman(time.Now().In(recipient.TZ()))),
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) SendSubscriptionNotification(recipient *models.User, hasExpired bool) error {
	tpl, err := m.getSubscriptionNotificationTemplate(SubscriptionNotificationTplData{
		PublicUrl:           m.config.Server.PublicUrl,
		DataRetentionMonths: m.config.App.DataRetentionMonths,
		HasExpired:          hasExpired,
	})
	if err != nil {
		return err
	}
	mail := &models.Mail{
		From:    models.MailAddress(m.config.Mail.Sender),
		To:      models.MailAddresses([]models.MailAddress{models.MailAddress(recipient.Email)}),
		Subject: subjectSubscriptionNotification,
	}
	mail.WithHTML(tpl.String())
	return m.sendingService.Send(mail)
}

func (m *MailService) getWelcomeTemplate(data WelcomeTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNameWelcome)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) getPasswordResetTemplate(data PasswordResetTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNamePasswordReset)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) getWakatimeFailureNotificationTemplate(data WakatimeFailureNotificationNotificationTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNameWakatimeFailureNotification)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) getImportNotificationTemplate(data ImportNotificationTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNameImportNotification)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) getReportTemplate(data ReportTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNameReport)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) getSubscriptionNotificationTemplate(data SubscriptionNotificationTplData) (*bytes.Buffer, error) {
	var rendered bytes.Buffer
	if err := m.templates[m.fmtName(tplNameSubscriptionNotification)].Execute(&rendered, data); err != nil {
		return nil, err
	}
	return &rendered, nil
}

func (m *MailService) fmtName(name string) string {
	return fmt.Sprintf("%s.tpl.html", name)
}
