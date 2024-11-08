package interfaces

type EmailSender interface {
	SendConfirmationEmail(email, code string) error
	SendRecoveryEmail(email, code string) error
}
