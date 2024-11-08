package repository

import (
	"awesomeProject4/user-auth-service/internal/config"
	"awesomeProject4/user-auth-service/internal/domains/interfaces"
	"fmt"
	"net/smtp"
)

type SMTPEmailSender struct {
	SMTPHost string
	SMTPPort string
	Username string
	Password string
}

func NewSMTPEmailSender(cfg *config.Config) interfaces.EmailSender {
	return &SMTPEmailSender{
		SMTPHost: cfg.SMTPHost,
		SMTPPort: cfg.SMTPPort,
		Username: cfg.SMTPUser,
		Password: cfg.SMTPPassword,
	}
}

func (s *SMTPEmailSender) SendConfirmationEmail(email, code string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.SMTPHost)
	message := []byte(fmt.Sprintf(
		"Subject: Подтверждение регистрации\n\nВаш код подтверждения: %s",
		code,
	))
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	return smtp.SendMail(addr, auth, s.Username, []string{email}, message)
}

func (s *SMTPEmailSender) SendRecoveryEmail(email, code string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.SMTPHost)
	message := []byte(fmt.Sprintf(
		"Subject: Восстановление пароля\n\nВаш код восстановления: %s",
		code,
	))
	addr := fmt.Sprintf("%s:%s", s.SMTPHost, s.SMTPPort)
	return smtp.SendMail(addr, auth, s.Username, []string{email}, message)
}
