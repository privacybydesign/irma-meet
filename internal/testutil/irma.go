// Package testutil holds fakes of the services irma-meet talks to: an IRMA
// server, a BigBlueButton server, an SMTP server and an in-memory store.
package testutil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
)

// FakeIRMA is an IRMA server that starts sessions and reports the results
// the test gives it.
type FakeIRMA struct {
	Server *httptest.Server

	apiToken string

	mu       sync.Mutex
	started  []json.RawMessage
	tokens   []string
	results  map[string]fakeResult
	failures map[string]string
	fetches  int
}

type fakeResult struct {
	attrs map[string]string
}

func NewFakeIRMA(t testing.TB, apiToken string) *FakeIRMA {
	f := &FakeIRMA{apiToken: apiToken, results: map[string]fakeResult{}, failures: map[string]string{}}
	f.Server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.Server.Close)

	return f
}

func (f *FakeIRMA) URL() string { return f.Server.URL }

// StartedCount is how many sessions were started.
func (f *FakeIRMA) StartedCount() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return len(f.tokens)
}

// ResultFetches is how many times a session result was requested.
func (f *FakeIRMA) ResultFetches() int {
	f.mu.Lock()
	defer f.mu.Unlock()

	return f.fetches
}

// Requests are the disclosure requests of the started sessions, oldest first.
func (f *FakeIRMA) Requests() []json.RawMessage {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]json.RawMessage(nil), f.started...)
}

// CompleteLatest makes the newest session end in a valid disclosure of attrs.
func (f *FakeIRMA) CompleteLatest(attrs map[string]string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.results[f.tokens[len(f.tokens)-1]] = fakeResult{attrs: attrs}
}

// FailLatest makes the newest session end with the given proof status.
func (f *FakeIRMA) FailLatest(proofStatus string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.failures[f.tokens[len(f.tokens)-1]] = proofStatus
}

func (f *FakeIRMA) serve(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != f.apiToken {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	switch {
	case r.Method == http.MethodPost && r.URL.Path == "/session":
		f.start(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/session/") && strings.HasSuffix(r.URL.Path, "/result"):
		f.result(w, strings.TrimSuffix(strings.TrimPrefix(r.URL.Path, "/session/"), "/result"))
	default:
		http.NotFound(w, r)
	}
}

func (f *FakeIRMA) start(w http.ResponseWriter, r *http.Request) {
	var request json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	f.mu.Lock()
	token := fmt.Sprintf("FakeToken%011d", len(f.tokens)+1)
	f.tokens = append(f.tokens, token)
	f.started = append(f.started, request)
	f.mu.Unlock()

	writeJSON(w, map[string]any{
		"token":      token,
		"sessionPtr": map[string]string{"u": f.Server.URL + "/irma/" + strings.Replace(token, "FakeToken", "ClientToken", 1), "irmaqr": "disclosing"},
	})
}

func (f *FakeIRMA) result(w http.ResponseWriter, token string) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.fetches++

	known := false

	for _, t := range f.tokens {
		known = known || t == token
	}

	if !known {
		http.NotFound(w, nil)
		return
	}

	if proof, failed := f.failures[token]; failed {
		writeJSON(w, map[string]any{"token": token, "status": "DONE", "type": "disclosing", "proofStatus": proof})
		return
	}

	res, done := f.results[token]
	if !done {
		writeJSON(w, map[string]any{"token": token, "status": "INITIALIZED", "type": "disclosing"})
		return
	}

	ids := make([]string, 0, len(res.attrs))
	for id := range res.attrs {
		ids = append(ids, id)
	}

	sort.Strings(ids)

	group := make([]map[string]any, 0, len(ids))
	for _, id := range ids {
		group = append(group, map[string]any{"id": id, "rawvalue": res.attrs[id], "status": "PRESENT"})
	}

	writeJSON(w, map[string]any{
		"token": token, "status": "DONE", "type": "disclosing", "proofStatus": "VALID",
		"disclosed": [][]map[string]any{group},
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
