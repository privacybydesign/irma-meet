package web_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

// keysInTemplates matches {{.T `key`}}, {{.HTML `key`}} and {{.HTMLf `key` ...}}.
var keysInTemplates = regexp.MustCompile("\\.(?:T|HTML|HTMLf) `([^`]*)`")

// keysInCode are the translation keys the Go code looks up.
var keysInCode = []string{
	"Page not found.",
	"Something went wrong. Please try again later.",
	"This meeting link has expired.",
	"Your attributes could not be verified by IRMA.",
	"Can't create room! please contact our administrator.",
	"The meeting could not be started because the Yivi server could not be reached. Please try again in a moment.",
	"Join as host", "Join as participant", "Join as teacher", "Join as student",
	"Meeting name", "Invitation note", "Participant email address",
	"The :attribute field is required.",
	"The :attribute field must not be greater than :max characters.",
	"The :attribute field must be a valid email address.",
	"Please check the box to allow processing.",
	"Yivi Meet link for :name",
	"Please send the above meeting link to people that you wish to invite to the video meeting.",
	"The meeting link has also been sent to the participant.",
	"The meeting link has also been sent to the participant, together with the following personal message: ",
	"This mail mostly contains general information, but specifically for this meeting the host has added the following personal message: ",
}

func TestTemplateKeysTranslated(t *testing.T) {
	raw, err := os.ReadFile("../i18n/locales/nl.json")
	if err != nil {
		t.Fatal(err)
	}

	dutch := map[string]string{}
	if err := json.Unmarshal(raw, &dutch); err != nil {
		t.Fatal(err)
	}

	keys := map[string]bool{}
	for _, k := range keysInCode {
		keys[k] = true
	}

	files, err := filepath.Glob("templates/*.html")
	if err != nil || len(files) == 0 {
		t.Fatalf("no templates found: %v", err)
	}

	for _, file := range files {
		src, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}

		for _, m := range keysInTemplates.FindAllStringSubmatch(string(src), -1) {
			keys[m[1]] = true
		}
	}

	for key := range keys {
		if dutch[key] == "" {
			t.Errorf("no Dutch translation for %.90q", key)
		}
	}
}
