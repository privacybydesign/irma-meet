package identity

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

const (
	requestTimeout  = 15 * time.Second
	maxResponseSize = 1 << 20

	statusDone       = "DONE"
	proofStatusValid = "VALID"
	attributePresent = "PRESENT"
)

var (
	// ErrNotVerified means the Yivi session ended without a valid disclosure.
	ErrNotVerified = errors.New("yivi session not verified")
	// ErrPending means the Yivi session is still running.
	ErrPending = errors.New("yivi session not finished")
)

// tokenPattern matches the session tokens the IRMA server hands out. Tokens
// end up in a request path, so nothing else is let through.
var tokenPattern = regexp.MustCompile(`^[A-Za-z0-9]{1,64}$`)

// running holds the IRMA session statuses before the session has ended.
var running = map[string]bool{"INITIALIZED": true, "PAIRING": true, "CONNECTED": true}

// Session is a started Yivi session.
type Session struct {
	// Token is the requestor token used to fetch the result.
	Token string
	// Ptr is what the Yivi frontend needs to reach the session.
	Ptr json.RawMessage
}

// Service runs disclosure sessions on an IRMA server.
type Service struct {
	baseURL  string
	apiToken string
	client   *http.Client
}

func NewService(baseURL, apiToken string, client *http.Client) *Service {
	if client == nil {
		client = &http.Client{Timeout: requestTimeout}
	}

	return &Service{baseURL: baseURL, apiToken: apiToken, client: client}
}

// Start begins one disclosure session of the given type.
func (s *Service) Start(ctx context.Context, t Type) (Session, error) {
	body, err := json.Marshal(t.request())
	if err != nil {
		return Session{}, err
	}

	var started struct {
		Token      string          `json:"token"`
		SessionPtr json.RawMessage `json:"sessionPtr"`
	}

	if err := s.do(ctx, http.MethodPost, "/session", body, &started); err != nil {
		return Session{}, err
	}

	if !tokenPattern.MatchString(started.Token) || len(started.SessionPtr) == 0 {
		return Session{}, errors.New("irma server returned an incomplete session")
	}

	return Session{Token: started.Token, Ptr: started.SessionPtr}, nil
}

type attribute struct {
	ID       string `json:"id"`
	RawValue string `json:"rawvalue"`
	Status   string `json:"status"`
}

// Verify fetches the result of a session started as type t. It returns
// ErrPending while the session runs and ErrNotVerified unless it ended in a
// complete and cryptographically valid disclosure.
func (s *Service) Verify(ctx context.Context, token string, t Type) (Verified, error) {
	if !tokenPattern.MatchString(token) {
		return Verified{}, ErrNotVerified
	}

	var result struct {
		Status      string        `json:"status"`
		ProofStatus string        `json:"proofStatus"`
		Disclosed   [][]attribute `json:"disclosed"`
	}

	path := "/session/" + url.PathEscape(token) + "/result"
	if err := s.do(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Verified{}, err
	}

	if running[result.Status] {
		return Verified{}, ErrPending
	}

	if result.Status != statusDone || result.ProofStatus != proofStatusValid {
		return Verified{}, ErrNotVerified
	}

	attrs := map[string]string{}

	for _, group := range result.Disclosed {
		for _, a := range group {
			if a.Status == attributePresent && a.RawValue != "" {
				attrs[a.ID] = a.RawValue
			}
		}
	}

	return Verified{Type: t, Attrs: attrs}, nil
}

func (s *Service) do(ctx context.Context, method, path string, body []byte, out any) error {
	req, err := http.NewRequestWithContext(ctx, method, s.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", s.apiToken)
	req.Header.Set("Accept", "application/json")

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("irma server: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(io.LimitReader(resp.Body, maxResponseSize))
	if err != nil {
		return fmt.Errorf("irma server: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("irma server: %s %s returned %d", method, path, resp.StatusCode)
	}

	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("irma server: invalid response: %w", err)
	}

	return nil
}
