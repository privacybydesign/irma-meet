package web_test

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/mailer"
	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/testutil"
	"github.com/privacybydesign/irma-meet/internal/web"
)

const (
	irmaToken = "irma-test-token"
	bbbSecret = "bbb-test-secret"

	attrEmail      = "pbdf.pbdf.email.email"
	attrFullName   = "pbdf.gemeente.personalData.fullname"
	attrSurfInst   = "pbdf.pbdf.surfnet-2.institute"
	attrSurfType   = "pbdf.pbdf.surfnet-2.type"
	attrSurfFamily = "pbdf.pbdf.surfnet-2.familyname"
	attrSurfEmail  = "pbdf.pbdf.surfnet-2.email"
)

type clock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *clock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

func (c *clock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)
}

// env is the whole service wired to fake IRMA, BBB and SMTP servers.
type env struct {
	t        *testing.T
	srv      *httptest.Server
	irma     *testutil.FakeIRMA
	bbb      *testutil.FakeBBB
	smtp     *testutil.FakeSMTP
	store    *testutil.MemStore
	clock    *clock
	meetings *meeting.Service
}

type envOption func(*envConfig)

type envConfig struct {
	nameFormat identity.NameFormat
	smtpPort   int
	log        *slog.Logger
}

func withNameFormat(f identity.NameFormat) envOption {
	return func(c *envConfig) { c.nameFormat = f }
}

func withLogger(l *slog.Logger) envOption {
	return func(c *envConfig) { c.log = l }
}

func withSMTPPort(port int) envOption {
	return func(c *envConfig) { c.smtpPort = port }
}

func newEnv(t *testing.T, opts ...envOption) *env {
	t.Helper()

	cfg := envConfig{nameFormat: identity.NameFormatString, log: slog.New(slog.NewTextHandler(io.Discard, nil))}
	for _, opt := range opts {
		opt(&cfg)
	}

	e := &env{
		t:     t,
		irma:  testutil.NewFakeIRMA(t, irmaToken),
		bbb:   testutil.NewFakeBBB(t, bbbSecret),
		smtp:  testutil.NewFakeSMTP(t),
		store: testutil.NewMemStore(),
		clock: &clock{now: time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)},
	}

	var handler http.Handler
	e.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	t.Cleanup(e.srv.Close)

	smtpPort := e.smtp.Port
	if cfg.smtpPort != 0 {
		smtpPort = cfg.smtpPort
	}

	mails, err := mailer.New(config.Mail{Host: e.smtp.Host, Port: smtpPort, FromAddress: "noreply@yivi-meet.test", FromName: "Yivi Meet"})
	if err != nil {
		t.Fatal(err)
	}

	text, err := i18n.New()
	if err != nil {
		t.Fatal(err)
	}

	e.meetings = meeting.NewService(meeting.Options{
		Store: e.store, Rooms: bbb.New(e.bbb.BaseURL(), bbbSecret, nil), Mails: mails, Text: text,
		BaseURL: e.srv.URL, NameFormat: cfg.nameFormat, Now: e.clock.Now,
	})

	appURL, _ := url.Parse(e.srv.URL)

	sessions, err := web.NewSessionManager(scs.New().Store, make([]byte, config.SessionKeyLen), false)
	if err != nil {
		t.Fatal(err)
	}

	handler, err = web.New(web.Options{
		Meetings:   e.meetings,
		Identities: identity.NewService(e.irma.URL(), irmaToken, nil),
		Text:       text,
		Sessions:   sessions,
		AppURL:     appURL,
		Log:        cfg.log,
	})
	if err != nil {
		t.Fatal(err)
	}

	return e
}

// user is one browser: it has its own cookies.
type user struct {
	e      *env
	client *http.Client
}

func (e *env) newUser() *user {
	jar, _ := cookiejar.New(nil)

	return &user{e: e, client: &http.Client{
		Jar:           jar,
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}}
}

type response struct {
	*http.Response
	Body string
}

func (u *user) do(req *http.Request) response {
	u.e.t.Helper()

	resp, err := u.client.Do(req)
	if err != nil {
		u.e.t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return response{Response: resp, Body: string(body)}
}

func (u *user) get(path string) response {
	u.e.t.Helper()

	req, _ := http.NewRequest(http.MethodGet, u.e.srv.URL+path, nil)

	return u.do(req)
}

func (u *user) post(path string, form url.Values, headers ...string) response {
	u.e.t.Helper()

	req, _ := http.NewRequest(http.MethodPost, u.e.srv.URL+path, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	for i := 0; i+1 < len(headers); i += 2 {
		req.Header.Set(headers[i], headers[i+1])
	}

	return u.do(req)
}

// disclose starts a Yivi session of the given type and has the person
// complete it with attrs, as the Yivi app would.
func (u *user) disclose(disclosureType string, attrs map[string]string) {
	u.e.t.Helper()

	resp := u.post("/irma_auth/start/"+disclosureType, nil)
	if resp.StatusCode != http.StatusOK {
		u.e.t.Fatalf("start %s: status %d: %s", disclosureType, resp.StatusCode, resp.Body)
	}

	u.e.irma.CompleteLatest(attrs)
}

func (r response) location() string {
	return r.Header.Get("Location")
}

func defaultAttrs(name, email string) map[string]string {
	return map[string]string{attrFullName: name, attrEmail: email}
}

func teacherAttrs(name, email string) map[string]string {
	return map[string]string{attrSurfInst: "uni.example", attrSurfType: "employee", attrSurfFamily: name, attrSurfEmail: email}
}

func studentAttrs(name, email string) map[string]string {
	return map[string]string{attrSurfInst: "uni.example", attrSurfType: "student", attrSurfFamily: name, attrSurfEmail: email}
}
