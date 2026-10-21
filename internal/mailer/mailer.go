// Package mailer renders the meeting mails and sends them over SMTP.
package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"strings"
	"time"

	mail "github.com/wneessen/go-mail"

	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/i18n"
)

const smtpTimeout = 15 * time.Second

//go:embed templates
var templateFS embed.FS

// Kind selects the mail template.
type Kind string

const (
	// Confirmation goes to the host with the meeting link.
	Confirmation Kind = "confirmation"
	// Invitation goes to a participant with the meeting link.
	Invitation Kind = "invitation"
)

// Data fills the templates. Everything in it is HTML-escaped on render.
type Data struct {
	HosterName     string
	InvitationLink string
	InvitationNote string
}

// Message is one mail to send.
type Message struct {
	Kind        Kind
	MeetingType string
	Locale      string
	To          string
	ReplyTo     string // optional
	Subject     string
	Data        Data
}

type Mailer struct {
	cfg       config.Mail
	templates map[string]*template.Template
}

func New(cfg config.Mail) (*Mailer, error) {
	m := &Mailer{cfg: cfg, templates: map[string]*template.Template{}}

	for _, locale := range []string{i18n.English, i18n.Dutch} {
		entries, err := templateFS.ReadDir("templates/" + locale)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			path := "templates/" + locale + "/" + entry.Name()

			tmpl, err := template.ParseFS(templateFS, path)
			if err != nil {
				return nil, err
			}

			m.templates[locale+"/"+strings.TrimSuffix(entry.Name(), ".html")] = tmpl
		}
	}

	return m, nil
}

// render returns the HTML body of msg.
func (m *Mailer) render(msg Message) (string, error) {
	name := i18n.Normalize(msg.Locale) + "/" + string(msg.Kind) + "_" + msg.MeetingType

	tmpl, ok := m.templates[name]
	if !ok {
		return "", fmt.Errorf("no mail template %q", name)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, msg.Data); err != nil {
		return "", err
	}

	return body.String(), nil
}

// Send delivers all messages over one SMTP connection.
func (m *Mailer) Send(ctx context.Context, msgs ...Message) error {
	client, err := m.client()
	if err != nil {
		return err
	}

	built := make([]*mail.Msg, 0, len(msgs))

	for _, msg := range msgs {
		out, err := m.build(msg)
		if err != nil {
			return err
		}

		built = append(built, out)
	}

	return client.DialAndSendWithContext(ctx, built...)
}

func (m *Mailer) build(msg Message) (*mail.Msg, error) {
	body, err := m.render(msg)
	if err != nil {
		return nil, err
	}

	out := mail.NewMsg()

	if err := out.FromFormat(m.cfg.FromName, m.cfg.FromAddress); err != nil {
		return nil, err
	}

	if err := out.To(msg.To); err != nil {
		return nil, err
	}

	if msg.ReplyTo != "" {
		if err := out.ReplyTo(msg.ReplyTo); err != nil {
			return nil, err
		}
	}

	out.Subject(singleLine(msg.Subject))
	out.SetBodyString(mail.TypeTextHTML, body)

	return out, nil
}

func (m *Mailer) client() (*mail.Client, error) {
	opts := []mail.Option{
		mail.WithPort(m.cfg.Port),
		mail.WithTimeout(smtpTimeout),
	}

	switch m.cfg.Encryption {
	case config.MailImplicitTLS:
		opts = append(opts, mail.WithSSL())
	case config.MailStartTLS:
		opts = append(opts, mail.WithTLSPolicy(mail.TLSMandatory))
	default:
		opts = append(opts, mail.WithTLSPolicy(mail.TLSOpportunistic))
	}

	if m.cfg.Username != "" {
		opts = append(opts,
			mail.WithSMTPAuth(mail.SMTPAuthAutoDiscover),
			mail.WithUsername(m.cfg.Username),
			mail.WithPassword(m.cfg.Password),
		)
	}

	return mail.NewClient(m.cfg.Host, opts...)
}

// singleLine keeps user-supplied text from adding headers to a mail.
func singleLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
