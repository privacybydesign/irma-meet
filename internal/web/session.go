package web

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/privacybydesign/irma-meet/internal/identity"
)

const (
	sessionLifetime     = 2 * time.Hour
	sessionCookie       = "session"
	secureSessionCookie = "__Host-session"

	keyLocale   = "locale"
	keyPending  = "pending"
	keyVerified = "verified"
)

// NewSessionManager returns a session manager over store. Session data is
// encrypted with key before it reaches the store, because it holds the names
// and email addresses people disclosed. The cookie is Secure when appURL is
// https.
func NewSessionManager(store scs.Store, key []byte, appURL *url.URL) (*scs.SessionManager, error) {
	codec, err := newSealedCodec(key)
	if err != nil {
		return nil, err
	}

	sm := scs.New()
	sm.Store = store
	sm.Codec = codec
	sm.Lifetime = sessionLifetime
	sm.Cookie.Name = sessionCookie
	sm.Cookie.HttpOnly = true
	sm.Cookie.SameSite = http.SameSiteLaxMode
	sm.Cookie.Secure = appURL.Scheme == "https"

	if sm.Cookie.Secure {
		sm.Cookie.Name = secureSessionCookie
	}

	return sm, nil
}

// sealedCodec encrypts what the gob codec produces with AES-GCM.
type sealedCodec struct {
	gcm   cipher.AEAD
	inner scs.GobCodec
}

func newSealedCodec(key []byte) (*sealedCodec, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return &sealedCodec{gcm: gcm}, nil
}

func (c *sealedCodec) Encode(deadline time.Time, values map[string]any) ([]byte, error) {
	plain, err := c.inner.Encode(deadline, values)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, c.gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return c.gcm.Seal(nonce, nonce, plain, nil), nil
}

// Decode treats data it cannot open (a rotated key, or tampering) as an
// empty session instead of failing every request that carries the cookie.
func (c *sealedCodec) Decode(data []byte) (time.Time, map[string]any, error) {
	nonceSize := c.gcm.NonceSize()
	if len(data) < nonceSize {
		return c.empty()
	}

	plain, err := c.gcm.Open(nil, data[:nonceSize], data[nonceSize:], nil)
	if err != nil {
		slog.Warn("session data could not be decrypted, starting a new session")
		return c.empty()
	}

	return c.inner.Decode(plain)
}

func (c *sealedCodec) empty() (time.Time, map[string]any, error) {
	return time.Now().Add(sessionLifetime), map[string]any{}, nil
}

// pending is a Yivi session that was started but not yet verified.
type pending struct {
	Token string        `json:"token"`
	Type  identity.Type `json:"type"`
}

// state wraps the per-visitor session data.
type state struct {
	sm *scs.SessionManager
}

func (s state) locale(r *http.Request) string {
	return s.sm.GetString(r.Context(), keyLocale)
}

func (s state) setLocale(r *http.Request, locale string) {
	s.sm.Put(r.Context(), keyLocale, locale)
}

func (s state) pending(r *http.Request) (pending, bool) {
	var p pending

	ok := s.load(r, keyPending, &p)

	return p, ok && p.Token != ""
}

func (s state) setPending(r *http.Request, p pending) {
	s.store(r, keyPending, p)
}

func (s state) clearPending(r *http.Request) {
	s.sm.Remove(r.Context(), keyPending)
}

func (s state) verified(r *http.Request) (identity.Verified, bool) {
	var v identity.Verified

	ok := s.load(r, keyVerified, &v)

	return v, ok && v.Type != ""
}

// setVerified stores a new identity under a fresh session token, so a token
// planted before the disclosure is of no use afterwards.
func (s state) setVerified(r *http.Request, v identity.Verified) error {
	if err := s.sm.RenewToken(r.Context()); err != nil {
		return err
	}

	s.store(r, keyVerified, v)

	return nil
}

func (s state) load(r *http.Request, key string, into any) bool {
	raw := s.sm.GetBytes(r.Context(), key)
	if raw == nil {
		return false
	}

	return json.NewDecoder(bytes.NewReader(raw)).Decode(into) == nil
}

func (s state) store(r *http.Request, key string, value any) {
	raw, err := json.Marshal(value)
	if err != nil {
		panic(errors.New("session value is not JSON serializable: " + err.Error()))
	}

	s.sm.Put(r.Context(), key, raw)
}
