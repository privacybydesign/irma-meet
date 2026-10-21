package testutil

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
)

// BBBCall is one API call the fake BigBlueButton server received.
type BBBCall struct {
	Name   string
	Params url.Values
}

// FakeBBB is a BigBlueButton server that checks the checksum of every call.
type FakeBBB struct {
	server *httptest.Server
	secret string

	// FailCreate makes create answer with a FAILED return code.
	FailCreate bool

	mu    sync.Mutex
	calls []BBBCall
	rooms map[string]bool
}

func NewFakeBBB(t testing.TB, secret string) *FakeBBB {
	f := &FakeBBB{secret: secret, rooms: map[string]bool{}}
	f.server = httptest.NewServer(http.HandlerFunc(f.serve))
	t.Cleanup(f.server.Close)

	return f
}

// BaseURL is the value for BBB_SERVER_BASE_URL.
func (f *FakeBBB) BaseURL() string { return f.server.URL + "/bigbluebutton/" }

func (f *FakeBBB) Calls(name string) []BBBCall {
	f.mu.Lock()
	defer f.mu.Unlock()

	var out []BBBCall

	for _, c := range f.calls {
		if c.Name == name {
			out = append(out, c)
		}
	}

	return out
}

func (f *FakeBBB) serve(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/bigbluebutton/api/")

	if !f.validChecksum(name, r.URL.RawQuery) {
		xml(w, "FAILED", "<messageKey>checksumError</messageKey>")
		return
	}

	params := r.URL.Query()

	f.mu.Lock()
	defer f.mu.Unlock()

	f.calls = append(f.calls, BBBCall{Name: name, Params: params})

	switch name {
	case "isMeetingRunning":
		running := "false"
		if f.rooms[params.Get("meetingID")] {
			running = "true"
		}

		xml(w, "SUCCESS", "<running>"+running+"</running>")
	case "create":
		if f.FailCreate {
			xml(w, "FAILED", "<messageKey>idNotUnique</messageKey>")
			return
		}

		f.rooms[params.Get("meetingID")] = true
		xml(w, "SUCCESS", "<meetingID>"+params.Get("meetingID")+"</meetingID>")
	default:
		http.NotFound(w, r)
	}
}

// validChecksum verifies sha256(call + query without checksum + secret).
func (f *FakeBBB) validChecksum(name, rawQuery string) bool {
	query, checksum, found := strings.Cut(rawQuery, "&checksum=")
	if !found {
		return false
	}

	sum := sha256.Sum256([]byte(name + query + f.secret))

	return hex.EncodeToString(sum[:]) == checksum
}

func xml(w http.ResponseWriter, code, extra string) {
	w.Header().Set("Content-Type", "text/xml")
	_, _ = w.Write([]byte("<response><returncode>" + code + "</returncode>" + extra + "</response>"))
}
