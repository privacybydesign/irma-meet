// Package config reads the process environment into a validated Config.
package config

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/privacybydesign/irma-meet/internal/identity"
)

const (
	defaultDBPort   = 3306
	defaultMailPort = 25

	// sessionKeyLen is the AES-256 key size used to encrypt stored sessions.
	sessionKeyLen = 32

	envNone = "null"
)

// MailEncryption is the MAIL_ENCRYPTION setting.
type MailEncryption string

const (
	// MailPlain connects without TLS, upgrading with STARTTLS when offered.
	MailPlain MailEncryption = ""
	// MailStartTLS requires STARTTLS.
	MailStartTLS MailEncryption = "tls"
	// MailImplicitTLS wraps the whole connection in TLS.
	MailImplicitTLS MailEncryption = "ssl"
)

type Database struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type Mail struct {
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  MailEncryption
	FromAddress string
	FromName    string
}

type Config struct {
	AppURL        *url.URL
	Database      Database
	Mail          Mail
	IrmaServerURL string
	IrmaAPIToken  string
	NameFormat    identity.NameFormat
	BBBBaseURL    string
	BBBSecret     string
	SessionKey    []byte
}

// Load builds a Config from getenv, reporting every problem at once.
func Load(getenv func(string) string) (Config, error) {
	l := loader{getenv: getenv}

	cfg := Config{
		AppURL:        l.appURL("APP_URL"),
		IrmaServerURL: strings.TrimRight(l.required("IRMA_SERVER_URL"), "/"),
		IrmaAPIToken:  l.required("IRMA_API_TOKEN"),
		NameFormat:    l.nameFormat("IRMA_MEET_NAME_FORMAT"),
		BBBBaseURL:    l.bbbBaseURL("BBB_SERVER_BASE_URL"),
		BBBSecret:     l.required("BBB_SECRET"),
		SessionKey:    l.sessionKey("SESSION_KEY"),
		Database: Database{
			Host:     l.required("DB_HOST"),
			Port:     l.port("DB_PORT", defaultDBPort),
			Name:     l.required("DB_DATABASE"),
			User:     l.required("DB_USERNAME"),
			Password: getenv("DB_PASSWORD"),
		},
		Mail: Mail{
			Host:        l.required("MAIL_HOST"),
			Port:        l.port("MAIL_PORT", defaultMailPort),
			Username:    getenv("MAIL_USERNAME"),
			Password:    getenv("MAIL_PASSWORD"),
			Encryption:  l.mailEncryption("MAIL_ENCRYPTION"),
			FromAddress: l.required("MAIL_FROM_ADDRESS"),
			FromName:    getenv("MAIL_FROM_NAME"),
		},
	}

	return cfg, l.err()
}

type loader struct {
	getenv func(string) string
	errs   []error
}

func (l *loader) fail(name, format string, args ...any) {
	l.errs = append(l.errs, fmt.Errorf("%s: %s", name, fmt.Sprintf(format, args...)))
}

func (l *loader) err() error {
	return errors.Join(l.errs...)
}

func (l *loader) required(name string) string {
	value := strings.TrimSpace(l.getenv(name))
	if value == "" {
		l.fail(name, "is required")
	}

	return value
}

func (l *loader) appURL(name string) *url.URL {
	raw := strings.TrimRight(l.required(name), "/")
	if raw == "" {
		return nil
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.Path != "" {
		l.fail(name, "must be an http(s) origin such as https://meet.example.org")
		return nil
	}

	return u
}

func (l *loader) bbbBaseURL(name string) string {
	raw := l.required(name)
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		l.fail(name, "must be an http(s) URL")
		return ""
	}

	return strings.TrimRight(raw, "/") + "/"
}

func (l *loader) nameFormat(name string) identity.NameFormat {
	switch strings.ToLower(strings.TrimSpace(l.getenv(name))) {
	case "", "string":
		return identity.NameFormatString
	case "json":
		return identity.NameFormatJSON
	}

	l.fail(name, `must be "string" or "json"`)

	return identity.NameFormatString
}

func (l *loader) port(name string, fallback int) int {
	raw := strings.TrimSpace(l.getenv(name))
	if raw == "" {
		return fallback
	}

	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		l.fail(name, "must be a port number")
		return fallback
	}

	return port
}

func (l *loader) mailEncryption(name string) MailEncryption {
	switch strings.ToLower(strings.TrimSpace(l.getenv(name))) {
	case "", envNone:
		return MailPlain
	case "tls":
		return MailStartTLS
	case "ssl":
		return MailImplicitTLS
	}

	l.fail(name, `must be "tls", "ssl" or "null"`)

	return MailPlain
}

func (l *loader) sessionKey(name string) []byte {
	raw := l.required(name)
	if raw == "" {
		return nil
	}

	key, err := decodeKey(raw)
	if err != nil || len(key) != sessionKeyLen {
		l.fail(name, "must be %d random bytes, hex or base64 encoded (for example `openssl rand -base64 32`)", sessionKeyLen)
		return nil
	}

	return key
}

func decodeKey(raw string) ([]byte, error) {
	if key, err := hex.DecodeString(raw); err == nil {
		return key, nil
	}

	if key, err := base64.StdEncoding.DecodeString(raw); err == nil {
		return key, nil
	}

	return base64.RawStdEncoding.DecodeString(raw)
}
