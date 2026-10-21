package web_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/meeting"
)

func TestStoreNeedsDisclosure(t *testing.T) {
	e := newEnv(t)
	anon := e.newUser()

	resp := anon.post("/irma_session/store", createForm("free", "Sneaky"))

	if resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("status %d, want a redirect to the disclosure page", resp.StatusCode)
	}

	if want := "/irma_session/authenticate/default?return=%2Firma_session%2Fcreate%2Ffree"; resp.location() != want {
		t.Errorf("redirect %q, want %q", resp.location(), want)
	}

	if len(e.store.Meetings()) != 0 || len(e.smtp.Mails()) != 0 {
		t.Error("an unauthenticated request created a meeting or sent mail")
	}
}

func TestCreateFormNeedsHostType(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	resp := user.get("/irma_session/create/exam")
	if want := "/irma_session/authenticate/teacher?return=%2Firma_session%2Fcreate%2Fexam"; resp.location() != want {
		t.Errorf("an email and name must not open the exam form: redirect %q, want %q", resp.location(), want)
	}
}

func TestJoinRoutesNeedDisclosure(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")
	anon := e.newUser()

	for _, path := range []string{"/irma_session/join_host/", "/irma_session/join_participant/"} {
		resp := anon.get(path + m.ID)

		if resp.StatusCode != http.StatusSeeOther || !strings.HasPrefix(resp.location(), "/irma_session/authenticate/default?return=") {
			t.Errorf("%s: status %d, location %q", path, resp.StatusCode, resp.location())
		}
	}

	if len(e.bbb.Calls("create")) != 0 {
		t.Error("an unauthenticated request created a BigBlueButton room")
	}
}

func TestUnknownMeetingIs404(t *testing.T) {
	e := newEnv(t)
	anon := e.newUser()

	for _, id := range []string{"deadbeefdeadbeefdeadbeefdeadbeef", "not-hex", "../etc/passwd", strings.Repeat("a", 200)} {
		for _, route := range []string{"join", "join_host", "join_participant"} {
			if resp := anon.get("/irma_session/" + route + "/" + id); resp.StatusCode == http.StatusOK || resp.StatusCode >= http.StatusInternalServerError {
				t.Errorf("%s/%s: status %d, want a rejection", route, id, resp.StatusCode)
			}
		}
	}

	if len(e.bbb.Calls("create")) != 0 {
		t.Error("an unknown meeting created a BigBlueButton room")
	}
}

func TestInvalidMeetingType(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	for _, name := range []string{"evil-not-a-type", "medical_consult", "../free", ""} {
		if resp := user.get("/irma_session/create/" + url.PathEscape(name)); resp.StatusCode != http.StatusNotFound {
			t.Errorf("create/%q: status %d, want 404", name, resp.StatusCode)
		}

		if resp := user.post("/irma_session/store", createForm(name, "x")); resp.StatusCode != http.StatusNotFound {
			t.Errorf("store with type %q: status %d, want 404", name, resp.StatusCode)
		}
	}

	if _, ok := meeting.ParseType("free"); !ok {
		t.Error("the configured meeting type free was not accepted")
	}
}

func TestStartKnownTypesOnly(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	for _, name := range []string{"medical", "patient", "evil", "Default"} {
		if resp := user.post("/irma_auth/start/"+name, nil); resp.StatusCode != http.StatusBadRequest {
			t.Errorf("start/%s: status %d, want 400", name, resp.StatusCode)
		}
	}

	if e.irma.StartedCount() != 0 {
		t.Error("an unknown disclosure type started a Yivi session")
	}
}

func TestStartOneSessionPerRequest(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	resp := user.post("/irma_auth/start/default", nil)
	if resp.StatusCode != http.StatusOK || !strings.Contains(resp.Body, `"u"`) {
		t.Fatalf("status %d, body %s: want the session pointer", resp.StatusCode, resp.Body)
	}

	if got := e.irma.StartedCount(); got != 1 {
		t.Errorf("%d Yivi sessions for one start request, want 1", got)
	}

	if strings.Contains(resp.Body, "FakeToken") {
		t.Error("the session token must stay on the server")
	}
}

func TestStartRequestsAttributes(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.post("/irma_auth/start/teacher", nil)

	request := string(e.irma.Requests()[0])
	for _, want := range []string{`"@context":"https://irma.app/ld/request/disclosure/v2"`, `{"type":"pbdf.pbdf.surfnet-2.type","value":"employee"}`, `"pbdf.pbdf.surfnet-2.email"`} {
		if !strings.Contains(request, want) {
			t.Errorf("teacher request misses %s: %s", want, request)
		}
	}
}

func TestFailureKeepsType(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.post("/irma_auth/start/teacher", nil)
	e.irma.FailLatest("INVALID_TIMESTAMP")

	resp := user.get("/irma_session/create/exam")
	if want := "/irma_session/authenticate/teacher?return=%2Firma_session%2Fcreate%2Fexam"; resp.location() != want {
		t.Errorf("redirect %q, want %q", resp.location(), want)
	}

	// The pending session is spent: asking again does not hit the IRMA server again.
	before := e.irma.ResultFetches()
	user.get("/irma_session/create/exam")

	if e.irma.ResultFetches() != before {
		t.Error("a failed session was verified twice")
	}
}

func TestExpiredSessionRestartsDisclosure(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.post("/irma_auth/start/default", nil)
	e.irma.ExpireLatest()

	for range 2 {
		resp := user.get("/irma_session/create/free")
		if want := "/irma_session/authenticate/default?return=%2Firma_session%2Fcreate%2Ffree"; resp.StatusCode != http.StatusSeeOther || resp.location() != want {
			t.Errorf("status %d, redirect %q, want a redirect to %q", resp.StatusCode, resp.location(), want)
		}
	}

	if got := e.irma.ResultFetches(); got != 1 {
		t.Errorf("an expired session was verified %d times, want once", got)
	}
}

func TestUncompletedDisclosure(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	user.post("/irma_auth/start/default", nil)

	resp := user.get("/irma_session/create/free")
	if resp.StatusCode != http.StatusSeeOther {
		t.Errorf("status %d: a disclosure the person never completed must not authenticate", resp.StatusCode)
	}
}

func TestDisclosureIsPerVisitor(t *testing.T) {
	e := newEnv(t)
	alice, mallory := e.newUser(), e.newUser()

	alice.disclose("default", defaultAttrs("Alice", "alice@example.com"))

	if resp := mallory.get("/irma_session/create/free"); resp.StatusCode != http.StatusSeeOther {
		t.Errorf("status %d, want mallory to be sent to the disclosure page", resp.StatusCode)
	}
}

func TestSessionTokenRenewed(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	cookieValue := func() string {
		u, _ := url.Parse(e.srv.URL)

		for _, c := range user.client.Jar.Cookies(u) {
			if c.Name == "session" {
				return c.Value
			}
		}

		return ""
	}

	user.post("/irma_auth/start/default", nil)
	before := cookieValue()

	e.irma.CompleteLatest(defaultAttrs("Hannah Host", "host@example.com"))
	user.get("/irma_session/create/free")

	if after := cookieValue(); before == "" || after == "" || before == after {
		t.Errorf("session token %q -> %q, want a new token after the disclosure", before, after)
	}
}

func TestAuthenticateLocalOnly(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	bad := []string{
		"https://evil.example/", "//evil.example/", `/\evil.example`, "javascript:alert(1)", "evil", "", "/ok\r\nSet-Cookie: x=y",
	}

	for _, target := range bad {
		resp := user.get("/irma_session/authenticate/default?return=" + url.QueryEscape(target))
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("return=%q: status %d, want 400", target, resp.StatusCode)
		}
	}

	if resp := user.get("/irma_session/authenticate/medical?return=%2F"); resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unknown disclosure type: status %d, want 400", resp.StatusCode)
	}

	ok := user.get("/irma_session/authenticate/default?return=" + url.QueryEscape("/irma_session/create/free"))
	if ok.StatusCode != http.StatusOK || !strings.Contains(ok.Body, "irma_auth/start/default") {
		t.Errorf("local return path: status %d", ok.StatusCode)
	}

	user.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	done := user.get("/irma_session/authenticate/default?return=" + url.QueryEscape("/irma_session/create/free"))
	if done.StatusCode != http.StatusSeeOther || done.location() != "/irma_session/create/free" {
		t.Errorf("already disclosed: status %d, location %q", done.StatusCode, done.location())
	}
}

func TestCrossOriginRejected(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	cases := map[string][]string{
		"foreign origin":     {headerOrigin, "https://evil.example"},
		"cross-site fetch":   {headerSecFetchSite, "cross-site"},
		"same-site sibling":  {headerSecFetchSite, "same-site"},
		"foreign origin too": {headerOrigin, "http://127.0.0.1:1", headerSecFetchSite, "cross-site"},
	}

	for name, headers := range cases {
		store := host.post("/irma_session/store", createForm("free", "Forged"), headers...)
		if store.StatusCode != http.StatusForbidden {
			t.Errorf("%s: store status %d, want 403", name, store.StatusCode)
		}

		start := host.post("/irma_auth/start/default", nil, headers...)
		if start.StatusCode != http.StatusForbidden {
			t.Errorf("%s: start status %d, want 403", name, start.StatusCode)
		}
	}

	if len(e.store.Meetings()) != 0 || len(e.smtp.Mails()) != 0 {
		t.Error("a cross-origin request created a meeting")
	}

	same := host.post("/irma_session/store", createForm("free", "Genuine"), headerOrigin, e.srv.URL, headerSecFetchSite, "same-origin")
	if same.StatusCode != http.StatusOK {
		t.Errorf("same-origin browser request: status %d, want 200", same.StatusCode)
	}
}

func TestMeetingNameIsHTMLEscaped(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	const payload = `<script>alert(1)</script>`

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	resp := host.post("/irma_session/store", createForm("free", payload))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %s", resp.StatusCode, resp.Body)
	}

	if strings.Contains(resp.Body, payload) {
		t.Error("the raw script tag reached the success page")
	}

	if !strings.Contains(resp.Body, "&lt;script&gt;alert(1)&lt;/script&gt;") {
		t.Error("the meeting name is not shown HTML-encoded")
	}

	for _, mail := range e.smtp.Mails() {
		if strings.Contains(mail.Body(), payload) {
			t.Error("the raw script tag reached a mail body")
		}
	}
}

func TestUserInputEscaped(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	const (
		name   = `Eve <img src=x onerror=alert(1)>`
		note   = `<b onmouseover=alert(2)>note</b>`
		quotes = `Quote " and ' and </script>`
	)

	host.disclose("teacher", teacherAttrs(name, "teacher@uni.example"))

	form := createForm("exam", quotes)
	form.Set("invitation_note", note)
	form.Set("participant_email_address1", "student@uni.example")

	page := host.get("/irma_session/create/exam")
	if strings.Contains(page.Body, "<img src=x") {
		t.Error("the disclosed name is not escaped on the create form")
	}

	resp := host.post("/irma_session/store", form)
	if strings.Contains(resp.Body, quotes) || strings.Contains(resp.Body, "</script>"+"Quote") {
		t.Error("the meeting name is not escaped on the success page")
	}

	for _, mail := range e.smtp.Mails() {
		body := mail.Body()
		for _, raw := range []string{"<img src=x", "<b onmouseover"} {
			if strings.Contains(body, raw) {
				t.Errorf("unescaped %q in mail to %v", raw, mail.To)
			}
		}
	}

	invalid := host.post("/irma_session/store", url.Values{"meeting_type": {"exam"}, "meeting_name": {`"><script>alert(3)</script>`}})
	if invalid.StatusCode != http.StatusUnprocessableEntity || strings.Contains(invalid.Body, "<script>alert(3)") {
		t.Errorf("status %d: submitted values must be escaped when the form is shown again", invalid.StatusCode)
	}
}

func TestSubjectNoHeaderInjection(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	resp := host.post("/irma_session/store", createForm("free", "Hi\r\nBcc: victim@example.com"))
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}

	for _, mail := range e.smtp.Mails() {
		if mail.Header(headerBcc) != "" || len(mail.To) != 1 {
			t.Errorf("header injection: Bcc %q, recipients %v", mail.Header(headerBcc), mail.To)
		}
	}
}

func TestFormValidation(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("teacher", teacherAttrs("Tanja Teacher", "tanja@uni.example"))

	cases := []struct {
		name  string
		form  url.Values
		field string
	}{
		{"missing name", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}}, "meeting_name"},
		{"too long name", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}, "meeting_name": {strings.Repeat("x", 256)}}, "meeting_name"},
		{"consent missing", url.Values{"meeting_type": {"exam"}, "meeting_name": {"x"}}, "agreed"},
		{"student without note", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}, "meeting_name": {"x"}, "participant_email_address1": {"s@uni.example"}}, "invitation_note"},
		{"note without student", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}, "meeting_name": {"x"}, "invitation_note": {"hi"}}, "participant_email_address1"},
		{"bad email", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}, "meeting_name": {"x"}, "invitation_note": {"hi"}, "participant_email_address1": {"not an email"}}, "participant_email_address1"},
		{"display name email", url.Values{"meeting_type": {"exam"}, "agreed": {"on"}, "meeting_name": {"x"}, "invitation_note": {"hi"}, "participant_email_address1": {"Eve <eve@uni.example>"}}, "participant_email_address1"},
	}

	for _, c := range cases {
		resp := host.post("/irma_session/store", c.form)

		if resp.StatusCode != http.StatusUnprocessableEntity || !strings.Contains(resp.Body, "alert-danger") || !strings.Contains(resp.Body, "is-invalid") {
			t.Errorf("%s: status %d, want the form again with errors on %s", c.name, resp.StatusCode, c.field)
		}
	}

	if len(e.store.Meetings()) != 0 {
		t.Error("an invalid form created a meeting")
	}
}

func TestOnlyStaticFilesAreServed(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	for _, path := range []string{
		"/.git/config", "/.git/HEAD", "/.env", "/.env.docker", "/go.mod", "/go.sum", "/main.go", "/Dockerfile", "/Dockerfile.test",
		"/docker-compose.yml", "/README.md", "/CLAUDE.md", "/Makefile", "/vendor/modules.txt", "/node_modules/", "/composer.json",
		"/internal/web/app.go", "/internal/store/migrations/00001_irma_meet_tables.sql", "/templates/layout.html", "/static/css/custom.css",
		"/locales/nl.json", "/css/../go.mod", "/css/%2e%2e/go.mod", "/css/..%2fgo.mod", "/img/../main.go", "/css/", "/js/", "/img/", "/fonts/", "/favicon/",
		"/index.php", "/artisan", "/storage/logs/laravel.log", "/web.go",
	} {
		resp := user.get(path)
		if resp.StatusCode == http.StatusOK {
			t.Errorf("GET %s: status 200, nothing but the static assets may be served", path)
		}

		if strings.Contains(resp.Body, "module github.com") || strings.Contains(resp.Body, "package main") {
			t.Errorf("GET %s leaked source", path)
		}
	}

	for path, contentType := range map[string]string{
		"/css/custom.css": "text/css", "/css/bootstrap.min.css": "text/css", "/js/custom.js": "text/javascript", "/js/irma.js": "text/javascript",
		"/js/bootstrap.bundle.min.js": "text/javascript", "/img/yivi-logo.svg": "image/svg+xml", "/favicon.ico": "image/", "/robots.txt": "text/plain",
	} {
		resp := user.get(path)
		if resp.StatusCode != http.StatusOK || !strings.Contains(resp.Header.Get(headerContentType), contentType) {
			t.Errorf("GET %s: status %d, content type %q", path, resp.StatusCode, resp.Header.Get(headerContentType))
		}
	}
}

func TestEmbeddedFilesArePublic(t *testing.T) {
	var files []string

	err := filepath.WalkDir("static", func(path string, d os.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, path)
		}

		return err
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		base := filepath.Base(f)
		if strings.HasPrefix(base, ".") || strings.HasSuffix(base, ".go") || strings.HasSuffix(base, ".md") || strings.HasSuffix(base, ".map") || strings.HasSuffix(base, ".pem") || strings.HasSuffix(base, ".key") {
			t.Errorf("static/%s does not belong in a public directory", f)
		}
	}
}

func TestSecurityHeaders(t *testing.T) {
	e := newEnv(t)
	user := e.newUser()

	resp := user.get("/")
	if resp.Header.Get(headerContentTypeOptions) != "nosniff" || resp.Header.Get(headerFrameOptions) != "DENY" {
		t.Errorf("security headers missing: %v", resp.Header)
	}

	if got := resp.Header.Get(headerCacheControl); got != "no-store" {
		t.Errorf("pages that show names and emails must not be cached, Cache-Control %q", got)
	}

	if got := user.get("/css/custom.css").Header.Get(headerCacheControl); got != "no-cache" {
		t.Errorf("static files are revalidated, Cache-Control %q", got)
	}
}

func TestAccessLogHidesMeetingIDs(t *testing.T) {
	var buf bytes.Buffer

	e := newEnv(t, withLogger(slog.New(slog.NewTextHandler(&buf, nil))))
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")
	host.get("/irma_session/join/" + m.ID)
	host.get("/irma_session/join_host/" + m.ID)

	log := buf.String()
	if strings.Contains(log, m.ID) {
		t.Errorf("the meeting id is in the access log:\n%s", log)
	}

	if !strings.Contains(log, "path=/irma_session/join_host/{id}") {
		t.Errorf("join requests are not logged:\n%s", log)
	}
}
