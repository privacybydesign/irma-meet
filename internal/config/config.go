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
	defaultMailPort = 25

	maxPort = 65535

	// sessionKeyLen is the AES-256 key size used to encrypt stored sessions.
	sessionKeyLen = 32

	schemeHTTP       = "http"
	schemeHTTPS      = "https"
	schemePostgres   = "postgres"
	schemePostgreSQL = "postgresql"
)

// MailEncryption is the MAIL_ENCRYPTION setting.
type MailEncryption string

const (
	// MailStartTLS requires STARTTLS. It is the default.
	MailStartTLS MailEncryption = "starttls"
	// MailImplicitTLS wraps the whole connection in TLS.
	MailImplicitTLS MailEncryption = "tls"
	// MailNone connects without TLS.
	MailNone MailEncryption = "none"
)

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
	DatabaseURL   string
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
		DatabaseURL:   l.databaseURL("DATABASE_URL"),
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

// LoadDatabaseURL reads only DATABASE_URL, for the migrate command, which
// needs no other setting.
func LoadDatabaseURL(getenv func(string) string) (string, error) {
	l := loader{getenv: getenv}
	databaseURL := l.databaseURL("DATABASE_URL")

	return databaseURL, l.err()
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
	if err != nil || (u.Scheme != schemeHTTP && u.Scheme != schemeHTTPS) || u.Host == "" || u.Path != "" {
		l.fail(name, "must be an http(s) origin such as https://meet.example.org")
		return nil
	}

	return u
}

// databaseURL checks the shape only and never repeats the value in an error,
// because it holds the password.
func (l *loader) databaseURL(name string) string {
	raw := l.required(name)
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != schemePostgres && u.Scheme != schemePostgreSQL) {
		l.fail(name, "must be a postgresql:// URL such as postgresql://user:password@host:5432/dbname")
		return ""
	}

	return raw
}

func (l *loader) bbbBaseURL(name string) string {
	raw := l.required(name)
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil || (u.Scheme != schemeHTTP && u.Scheme != schemeHTTPS) || u.Host == "" {
		l.fail(name, "must be an http(s) URL")
		return ""
	}

	return strings.TrimRight(raw, "/") + "/"
}

func (l *loader) nameFormat(name string) identity.NameFormat {
	switch strings.ToLower(strings.TrimSpace(l.getenv(name))) {
	case "", string(identity.NameFormatString):
		return identity.NameFormatString
	case string(identity.NameFormatJSON):
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
	if err != nil || port < 1 || port > maxPort {
		l.fail(name, "must be a port number")
		return fallback
	}

	return port
}

func (l *loader) mailEncryption(name string) MailEncryption {
	switch strings.ToLower(strings.TrimSpace(l.getenv(name))) {
	case "", string(MailStartTLS):
		return MailStartTLS
	case string(MailImplicitTLS):
		return MailImplicitTLS
	case string(MailNone):
		return MailNone
	}

	l.fail(name, `must be "starttls", "tls" or "none"`)

	return MailStartTLS
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
