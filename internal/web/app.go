// Package web is the HTTP layer: routing, sessions, forms and pages.
package web

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/alexedwards/scs/v2"

	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/meeting"
)

const (
	maxFormBytes = 64 << 10

	headerCacheControl       = "Cache-Control"
	headerContentType        = "Content-Type"
	headerContentTypeOptions = "X-Content-Type-Options"
	headerFrameOptions       = "X-Frame-Options"

	cacheNoCache    = "no-cache"
	cacheNoStore    = "no-store"
	contentTypeHTML = "text/html; charset=utf-8"
	contentTypeJSON = "application/json"
	frameDeny       = "DENY"
	nosniff         = "nosniff"

	// msgYiviUnreachable is a translation key.
	msgYiviUnreachable = "The meeting could not be started because the Yivi server could not be reached. Please try again in a moment."
)

// Meetings is what the web layer needs from the meeting domain.
type Meetings interface {
	Create(ctx context.Context, in meeting.NewMeeting) (meeting.Created, error)
	Lookup(ctx context.Context, id string) (meeting.Meeting, error)
	JoinURL(ctx context.Context, m meeting.Meeting, who identity.Verified) (string, error)
}

// Identities is what the web layer needs to run Yivi disclosures.
type Identities interface {
	Start(ctx context.Context, t identity.Type) (identity.Session, error)
	Verify(ctx context.Context, token string, t identity.Type) (identity.Verified, error)
}

type Options struct {
	Meetings   Meetings
	Identities Identities
	Text       *i18n.Translator
	Sessions   *scs.SessionManager
	// AppURL is the public origin of the service.
	AppURL *url.URL
	Log    *slog.Logger
}

type app struct {
	meetings   Meetings
	identities Identities
	text       *i18n.Translator
	state      state
	views      *renderer
	log        *slog.Logger
}

// New returns the HTTP handler of the whole service.
func New(o Options) (http.Handler, error) {
	views, err := newRenderer(o.Text)
	if err != nil {
		return nil, err
	}

	static, err := staticHandler()
	if err != nil {
		return nil, err
	}

	a := &app{
		meetings:   o.Meetings,
		identities: o.Identities,
		text:       o.Text,
		state:      state{sm: o.Sessions},
		views:      views,
		log:        o.Log,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.home)
	mux.HandleFunc("GET /lang/{locale}", a.setLocale)
	mux.HandleFunc("POST /irma_auth/start/{disclosureType}", a.startSession)
	mux.HandleFunc("GET /irma_session/authenticate/{disclosureType}", a.authenticate)
	mux.HandleFunc("GET /irma_session/create/{meetingType}", a.createForm)
	mux.HandleFunc("POST /irma_session/store", a.store)
	mux.HandleFunc("GET /irma_session/join/{irmaSessionId}", a.join)
	mux.HandleFunc("GET /irma_session/join_host/{irmaSessionId}", a.joinHost)
	mux.HandleFunc("GET /irma_session/join_participant/{irmaSessionId}", a.joinParticipant)
	mux.Handle("GET /css/", static)
	mux.Handle("GET /js/", static)
	mux.Handle("GET /img/", static)
	mux.Handle("GET /fonts/", static)
	mux.Handle("GET /favicon/", static)
	mux.Handle("GET /favicon.ico", static)
	mux.Handle("GET /robots.txt", static)
	mux.HandleFunc("/", a.notFound)

	csrf := http.NewCrossOriginProtection()
	if err := csrf.AddTrustedOrigin(o.AppURL.Scheme + "://" + o.AppURL.Host); err != nil {
		return nil, fmt.Errorf("APP_URL as trusted origin: %w", err)
	}

	var h http.Handler = mux
	h = csrf.Handler(h)
	h = o.Sessions.LoadAndSave(h)
	h = a.secureHeaders(h)
	h = a.recoverPanics(h)

	return a.accessLog(h), nil
}

func (a *app) lang(r *http.Request) string {
	return i18n.Normalize(a.state.locale(r))
}

func (a *app) page(w http.ResponseWriter, r *http.Request, page string, status int, data any) {
	if err := a.views.render(w, page, a.lang(r), status, data); err != nil {
		a.log.Error("render failed", "page", page, "error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type errorPage struct {
	// Message is a translation key.
	Message string
}

func (a *app) fail(w http.ResponseWriter, r *http.Request, status int, messageKey string) {
	a.page(w, r, pageError, status, errorPage{Message: messageKey})
}

func (a *app) notFound(w http.ResponseWriter, r *http.Request) {
	a.fail(w, r, http.StatusNotFound, "Page not found.")
}

func (a *app) serverError(w http.ResponseWriter, r *http.Request, what string, err error) {
	a.log.Error(what, "error", err)
	a.fail(w, r, http.StatusInternalServerError, "Something went wrong. Please try again later.")
}

func (a *app) home(w http.ResponseWriter, r *http.Request) {
	a.page(w, r, pageHome, http.StatusOK, nil)
}

// setLocale switches the interface language and returns to the home page.
func (a *app) setLocale(w http.ResponseWriter, r *http.Request) {
	if locale := r.PathValue("locale"); i18n.Supported(locale) {
		a.state.setLocale(r, locale)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
