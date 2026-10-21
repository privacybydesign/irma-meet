package web_test

import (
	"net/http"
	"strings"
	"testing"
)

// The probes of the deployment request /.
func TestHomeAnswers200(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	if resp := user.get("/"); resp.StatusCode != http.StatusOK || !strings.Contains(resp.Body, "Yivi Meet") {
		t.Errorf("GET /: status %d", resp.StatusCode)
	}

	req, _ := http.NewRequest(http.MethodHead, e.srv.URL+"/", nil)

	if resp := user.do(req); resp.StatusCode != http.StatusOK {
		t.Errorf("HEAD /: status %d", resp.StatusCode)
	}
}

func TestUnknownPageIs404(t *testing.T) {
	e := newEnv(t)

	if resp := e.newUser().get("/admin"); resp.StatusCode != http.StatusNotFound {
		t.Errorf("GET /admin: status %d, want 404 (the Laravel admin page is gone)", resp.StatusCode)
	}

	for _, path := range []string{"/login", "/register", "/password/reset", "/home"} {
		if resp := e.newUser().get(path); resp.StatusCode != http.StatusNotFound {
			t.Errorf("GET %s: status %d, want 404 (the Laravel auth routes are gone)", path, resp.StatusCode)
		}
	}
}

// The failure messages are injected for the Yivi frontend, in the visitor's language.
func TestEnglishFailureMessagesAreInjected(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.get("/lang/en")

	resp := user.get("/")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	for _, want := range []string{
		"window.yiviSessionMessages",
		"the Yivi session was cancelled",
		"Yivi session timed out",
		"Read which credentials you need and how to load them",
		"https://privacybydesign.foundation/uitgifte/",
	} {
		if !strings.Contains(resp.Body, want) {
			t.Errorf("home page misses %q", want)
		}
	}
}

func TestDutchFailureMessagesAreTranslated(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.get("/lang/nl")

	resp := user.get("/")
	for _, want := range []string{"window.yiviSessionMessages", "de Yivi-sessie is geannuleerd", "Lees welke gegevens je nodig hebt"} {
		if !strings.Contains(resp.Body, want) {
			t.Errorf("home page misses %q", want)
		}
	}
}

func TestLanguageSwitch(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	if resp := user.get("/"); !strings.Contains(resp.Body, `<html lang="nl">`) {
		t.Error("the default language is Dutch")
	}

	if resp := user.get("/lang/en"); resp.StatusCode != http.StatusSeeOther || resp.location() != "/" {
		t.Errorf("switching: status %d, location %q", resp.StatusCode, resp.location())
	}

	if resp := user.get("/"); !strings.Contains(resp.Body, `<html lang="en">`) {
		t.Error("the language did not switch to English")
	}

	user.get("/lang/xx")

	if resp := user.get("/"); !strings.Contains(resp.Body, `<html lang="en">`) {
		t.Error("an unsupported language must be ignored")
	}
}

func TestEveryPageRendersInBothLanguages(t *testing.T) {
	e := newEnv(t)

	for _, lang := range []string{"nl", "en"} {
		user := e.newUser()
		user.get("/lang/" + lang)

		for _, path := range []string{"/", "/irma_session/authenticate/default?return=%2F", "/nope"} {
			resp := user.get(path)
			if resp.StatusCode >= http.StatusInternalServerError || strings.Contains(resp.Body, "{{") {
				t.Errorf("%s (%s): status %d, unrendered template markup or error", path, lang, resp.StatusCode)
			}
		}

		user.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

		if resp := user.get("/irma_session/create/free"); resp.StatusCode != http.StatusOK {
			t.Errorf("create form (%s): status %d", lang, resp.StatusCode)
		}
	}
}
