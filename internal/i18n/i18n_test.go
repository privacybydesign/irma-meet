package i18n_test

import (
	"testing"

	"github.com/privacybydesign/irma-meet/internal/i18n"
)

func TestTranslations(t *testing.T) {
	tr, err := i18n.New()
	if err != nil {
		t.Fatal(err)
	}

	if got := tr.T("nl", "Send invite"); got != "Stuur uitnodiging" {
		t.Errorf("nl: %q", got)
	}

	if got := tr.T("en", "Send invite"); got != "Send invite" {
		t.Errorf("an empty or missing English entry falls back to the key, got %q", got)
	}

	if got := tr.T("en", "demo_notice_label"); got == "demo_notice_label" {
		t.Errorf("English overrides must apply, got %q", got)
	}

	if got := tr.T("nl", "no such key"); got != "no such key" {
		t.Errorf("unknown key: %q", got)
	}

	if got := tr.T("fr", "Send invite"); got != "Stuur uitnodiging" {
		t.Errorf("an unsupported locale uses the default (nl), got %q", got)
	}
}

func TestPlaceholders(t *testing.T) {
	tr, _ := i18n.New()

	if got := tr.Tr("nl", "Yivi Meet link for :name", map[string]string{"name": "Sync"}); got != "Yivi Meet link voor Sync" {
		t.Errorf("%q", got)
	}

	if got := tr.Tr("en", "Yivi Meet link for :name", map[string]string{"name": "Sync"}); got != "Yivi Meet link for Sync" {
		t.Errorf("%q", got)
	}
}

func TestNormalizeAndSupported(t *testing.T) {
	for in, want := range map[string]string{"en": "en", "nl": "nl", "": "nl", "de": "nl", "EN": "nl"} {
		if got := i18n.Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}

	if !i18n.Supported("en") || !i18n.Supported("nl") || i18n.Supported("de") {
		t.Error("Supported is wrong")
	}
}
