package config_test

import (
	"encoding/base64"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/identity"
)

const (
	keyLen = 32

	testDatabaseURL = "postgresql://irmameet:s3cret@db:5432/irmameet?connect_timeout=5"

	envDatabaseURL = "DATABASE_URL"
)

func validEnv() map[string]string {
	return map[string]string{
		"APP_URL":               "https://meet.example.org/",
		"SESSION_KEY":           base64.StdEncoding.EncodeToString(make([]byte, keyLen)),
		envDatabaseURL:          testDatabaseURL,
		"MAIL_HOST":             "mail",
		"MAIL_FROM_ADDRESS":     "noreply@example.org",
		"IRMA_SERVER_URL":       "https://irma.example.org/",
		"IRMA_API_TOKEN":        "token",
		"BBB_SERVER_BASE_URL":   "https://bbb.example.org/bigbluebutton",
		"BBB_SECRET":            "bbb-secret",
		"IRMA_MEET_NAME_FORMAT": "JSON",
		"MAIL_ENCRYPTION":       "tls",
	}
}

func load(env map[string]string) (config.Config, error) {
	return config.Load(func(k string) string { return env[k] })
}

func TestLoadValidConfig(t *testing.T) {
	cfg, err := load(validEnv())
	if err != nil {
		t.Fatal(err)
	}

	if cfg.AppURL.String() != "https://meet.example.org" || cfg.IrmaServerURL != "https://irma.example.org" {
		t.Errorf("urls %s %s", cfg.AppURL, cfg.IrmaServerURL)
	}

	if cfg.BBBBaseURL != "https://bbb.example.org/bigbluebutton/" {
		t.Errorf("BBB base url %q must end in a slash", cfg.BBBBaseURL)
	}

	if cfg.NameFormat != identity.NameFormatJSON || cfg.Mail.Encryption != config.MailImplicitTLS {
		t.Errorf("name format %q, encryption %q", cfg.NameFormat, cfg.Mail.Encryption)
	}

	if cfg.DatabaseURL != testDatabaseURL || cfg.Mail.Port != 25 {
		t.Errorf("database url %q, default mail port %d", cfg.DatabaseURL, cfg.Mail.Port)
	}

	if len(cfg.SessionKey) != keyLen {
		t.Errorf("session key has %d bytes", len(cfg.SessionKey))
	}
}

func TestMissingVarsAllReported(t *testing.T) {
	_, err := load(map[string]string{})
	if err == nil {
		t.Fatal("expected an error")
	}

	for _, name := range []string{"APP_URL", "SESSION_KEY", "DATABASE_URL", "MAIL_HOST", "MAIL_FROM_ADDRESS", "IRMA_SERVER_URL", "IRMA_API_TOKEN", "BBB_SERVER_BASE_URL", "BBB_SECRET"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("error does not mention %s: %v", name, err)
		}
	}
}

func TestSessionKeyFormats(t *testing.T) {
	good := map[string]string{
		"hex":      hex.EncodeToString(make([]byte, keyLen)),
		"base64":   base64.StdEncoding.EncodeToString(make([]byte, keyLen)),
		"unpadded": base64.RawStdEncoding.EncodeToString(append(make([]byte, keyLen-1), 1)),
	}

	for name, key := range good {
		env := validEnv()
		env["SESSION_KEY"] = key

		if _, err := load(env); err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}

	for name, key := range map[string]string{
		"too short":         base64.StdEncoding.EncodeToString(make([]byte, 16)),
		"too long":          hex.EncodeToString(make([]byte, keyLen+1)),
		"garbage":           "not a key!",
		"old APP_KEY style": "base64:" + base64.StdEncoding.EncodeToString(make([]byte, keyLen)),
	} {
		env := validEnv()
		env["SESSION_KEY"] = key

		if _, err := load(env); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestAppURLMustBeAnOrigin(t *testing.T) {
	for _, bad := range []string{"meet.example.org", "ftp://meet.example.org", "https://meet.example.org/app", "https://"} {
		env := validEnv()
		env["APP_URL"] = bad

		if _, err := load(env); err == nil {
			t.Errorf("APP_URL %q accepted", bad)
		}
	}
}

func TestMailEncryptionAndFormat(t *testing.T) {
	for value, want := range map[string]config.MailEncryption{"": config.MailStartTLS, "starttls": config.MailStartTLS, "tls": config.MailImplicitTLS, "TLS": config.MailImplicitTLS, "none": config.MailNone} {
		env := validEnv()
		env["MAIL_ENCRYPTION"] = value

		cfg, err := load(env)
		if err != nil || cfg.Mail.Encryption != want {
			t.Errorf("MAIL_ENCRYPTION=%q: %q, %v", value, cfg.Mail.Encryption, err)
		}
	}

	for _, old := range []string{"ssl", "null"} {
		env := validEnv()
		env["MAIL_ENCRYPTION"] = old

		if _, err := load(env); err == nil {
			t.Errorf("MAIL_ENCRYPTION=%q must be rejected", old)
		}
	}

	env := validEnv()
	env["IRMA_MEET_NAME_FORMAT"] = "xml"

	if _, err := load(env); err == nil {
		t.Error("an unknown IRMA_MEET_NAME_FORMAT must be rejected")
	}

	env = validEnv()
	env["IRMA_MEET_NAME_FORMAT"] = ""

	if cfg, err := load(env); err != nil || cfg.NameFormat != identity.NameFormatString {
		t.Errorf("default name format %q, %v", cfg.NameFormat, err)
	}
}

func TestPortsMustBeNumbers(t *testing.T) {
	for _, bad := range []string{"abc", "0", "70000"} {
		env := validEnv()
		env["MAIL_PORT"] = bad

		if _, err := load(env); err == nil {
			t.Errorf("MAIL_PORT=%q accepted", bad)
		}
	}
}

func TestDatabaseURL(t *testing.T) {
	for _, good := range []string{testDatabaseURL, "postgres://u:p@localhost/db", "postgresql://u@localhost:5432/db?sslmode=require"} {
		env := validEnv()
		env[envDatabaseURL] = good

		if _, err := load(env); err != nil {
			t.Errorf("DATABASE_URL %q rejected: %v", good, err)
		}
	}

	for _, bad := range []string{"mysql://u:p@localhost/db", "u:p@localhost/db", "host=localhost user=u", "postgresql://u:p@local host/db"} {
		env := validEnv()
		env[envDatabaseURL] = bad

		_, err := load(env)
		if err == nil {
			t.Errorf("DATABASE_URL %q accepted", bad)
			continue
		}

		if strings.Contains(err.Error(), "p@") || strings.Contains(err.Error(), bad) {
			t.Errorf("the error repeats the connection string: %v", err)
		}
	}
}

func TestLoadDatabaseURLNeedsNoOtherVariable(t *testing.T) {
	got, err := config.LoadDatabaseURL(func(k string) string { return map[string]string{envDatabaseURL: testDatabaseURL}[k] })
	if err != nil || got != testDatabaseURL {
		t.Errorf("got %q, %v", got, err)
	}

	if _, err := config.LoadDatabaseURL(func(string) string { return "" }); err == nil || !strings.Contains(err.Error(), envDatabaseURL) {
		t.Errorf("a missing DATABASE_URL gave %v", err)
	}
}
