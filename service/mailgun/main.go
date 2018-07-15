package mailgun

import (
	"github.com/keller0/yxi-back/internal"
	mail "github.com/mailgun/mailgun-go"
)

var (
	apiKey = internal.GetEnv("MAILGUN_API_KEY", "private key")
	pubkey = internal.GetEnv("MAILGUN_PUB_KEY", "public key")
	domain = internal.GetEnv("MAILGUN_DOMAIN", "mail domain")
	// ServiceAccount account info used in emails
	ServiceAccount = "YXI <no-reply@yxi.io>"
)

// SimpleMessage send text to user
func SimpleMessage(subject, content, userEmail string) (string, error) {
	mg := mail.NewMailgun(domain, apiKey, pubkey)
	m := mg.NewMessage(
		ServiceAccount,
		subject,
		content,
		userEmail,
	)
	_, id, err := mg.Send(m)
	return id, err
}
