package mail

import (
	"net/smtp"

	"github.com/arthurbailao/suno-kindle/domain"
	"github.com/arthurbailao/suno-kindle/suno"
	"github.com/domodwyer/mailyak"
	"github.com/pkg/errors"
)

// Mail ...
type Mail struct {
	config Config
}

// Config ...
type Config struct {
	Host     string
	Auth     smtp.Auth
	From     string
	FromName string
}

// New ...
func New(cfg Config) *Mail {
	return &Mail{config: cfg}
}

// SendMail ...
func (m Mail) SendMail(device domain.Device, r suno.Report) error {
	mail := mailyak.New(m.config.Host, m.config.Auth)
	mail.From(m.config.From)
	mail.FromName(m.config.FromName)
	mail.To(device.Email)
	mail.Subject("Convert: " + r.Title)

	file, err := r.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	mail.Attach(r.Filename, file)

	if err := mail.Send(); err != nil {
		return errors.Wrap(err, "failed to send mail")
	}
	return nil
}
