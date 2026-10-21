// Package i18n looks up UI and mail strings. Keys are the English source
// text, so a missing or empty translation falls back to the key itself.
package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	English = "en"
	Dutch   = "nl"

	DefaultLocale = Dutch
)

//go:embed locales/*.json
var localeFS embed.FS

type Translator struct {
	catalogs map[string]map[string]string
}

func New() (*Translator, error) {
	t := &Translator{catalogs: map[string]map[string]string{}}

	for _, locale := range []string{English, Dutch} {
		data, err := localeFS.ReadFile("locales/" + locale + ".json")
		if err != nil {
			return nil, err
		}

		catalog := map[string]string{}
		if err := json.Unmarshal(data, &catalog); err != nil {
			return nil, fmt.Errorf("locale %s: %w", locale, err)
		}
		t.catalogs[locale] = catalog
	}

	return t, nil
}

// Normalize maps any input to a supported locale.
func Normalize(locale string) string {
	if locale == English || locale == Dutch {
		return locale
	}

	return DefaultLocale
}

func Supported(locale string) bool {
	return locale == English || locale == Dutch
}

// T returns the translation of key, or key itself when there is none.
func (t *Translator) T(locale, key string) string {
	if line := t.catalogs[Normalize(locale)][key]; line != "" {
		return line
	}

	return key
}

// Tr is T followed by replacing ":name" style placeholders.
func (t *Translator) Tr(locale, key string, replace map[string]string) string {
	line := t.T(locale, key)
	for name, value := range replace {
		line = strings.ReplaceAll(line, ":"+name, value)
	}

	return line
}

// Has reports whether key has a non-empty translation in locale.
func (t *Translator) Has(locale, key string) bool {
	return t.catalogs[Normalize(locale)][key] != ""
}
