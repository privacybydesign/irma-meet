package web

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/privacybydesign/irma-meet/internal/identity"
)

const authPath = "/irma_session/authenticate/"

// startSession starts one Yivi session and hands its pointer to the Yivi
// frontend. The session token stays on the server.
func (a *app) startSession(w http.ResponseWriter, r *http.Request) {
	t, ok := identity.Parse(r.PathValue("disclosureType"))
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	session, err := a.identities.Start(r.Context(), t)
	if err != nil {
		a.log.Error("starting yivi session failed", "error", err)
		http.Error(w, "Internal server error", http.StatusBadGateway)

		return
	}

	a.state.setPending(r, pending{Token: session.Token, Type: t})

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(session.Ptr)
}

type authenticatePage struct {
	StartURL  string
	ReturnURL string
}

// authenticate shows the page that starts a disclosure and then continues to
// the return path.
func (a *app) authenticate(w http.ResponseWriter, r *http.Request) {
	t, ok := identity.Parse(r.PathValue("disclosureType"))
	next := r.URL.Query().Get("return")

	if !ok || !isLocalPath(next) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if v, ok, err := a.identityOf(r); err == nil && ok && v.Type == t {
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	a.page(w, r, pageAuthenticate, http.StatusOK, authenticatePage{StartURL: startPath(t), ReturnURL: next})
}

// identityOf returns who the visitor is. A Yivi session that was started and
// has not been looked at yet is verified first and replaces what the session
// held, so the latest disclosure wins.
func (a *app) identityOf(r *http.Request) (identity.Verified, bool, error) {
	if p, ok := a.state.pending(r); ok {
		v, err := a.identities.Verify(r.Context(), p.Token, p.Type)

		switch {
		case err == nil:
			a.state.clearPending(r)

			if err := a.state.setVerified(r, v); err != nil {
				return identity.Verified{}, false, err
			}

			return v, true, nil
		case errors.Is(err, identity.ErrNotVerified):
			a.state.clearPending(r)
		case !errors.Is(err, identity.ErrPending):
			return identity.Verified{}, false, err
		}
	}

	v, ok := a.state.verified(r)

	return v, ok, nil
}

// identityFor returns who the visitor has proven to be, as t. When they have
// not (yet), it responds itself, by sending them to the disclosure page for
// t with returnTo as the place to continue, and reports false.
func (a *app) identityFor(w http.ResponseWriter, r *http.Request, t identity.Type, returnTo string) (identity.Verified, bool) {
	v, ok, err := a.identityOf(r)
	if err != nil {
		a.log.Error("verifying yivi session failed", "error", err)
		a.fail(w, r, http.StatusBadGateway, "The meeting could not be started because the Yivi server could not be reached. Please try again in a moment.")

		return identity.Verified{}, false
	}

	if ok && v.Type == t {
		return v, true
	}

	http.Redirect(w, r, authPath+url.PathEscape(string(t))+"?return="+url.QueryEscape(returnTo), http.StatusSeeOther)

	return identity.Verified{}, false
}

func startPath(t identity.Type) string {
	return "/irma_auth/start/" + string(t)
}

// isLocalPath accepts only paths on this site, so the return parameter
// cannot be used as an open redirect.
func isLocalPath(p string) bool {
	if !strings.HasPrefix(p, "/") || strings.HasPrefix(p, "//") || strings.HasPrefix(p, `/\`) {
		return false
	}

	if strings.ContainsAny(p, "\\\r\n\t") {
		return false
	}

	u, err := url.Parse(p)

	return err == nil && u.Scheme == "" && u.Host == ""
}
