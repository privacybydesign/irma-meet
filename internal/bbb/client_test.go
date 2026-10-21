package bbb_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

const secret = "bbb-secret"

func checksum(call, query string) string {
	sum := sha256.Sum256([]byte(call + query + secret))

	return hex.EncodeToString(sum[:])
}

func TestJoinURLIsSigned(t *testing.T) {
	c := bbb.New("https://bbb.example/bigbluebutton/", secret, nil)

	got := c.JoinURL("room1", "Alice & Bob", "pw")

	query := "fullName=Alice+%26+Bob&meetingID=room1&password=pw&redirect=true"
	want := "https://bbb.example/bigbluebutton/api/join?" + query + "&checksum=" + checksum("join", query)

	if got != want {
		t.Errorf("\n got %s\nwant %s", got, want)
	}
}

func TestEnsureRoomCreatesOnce(t *testing.T) {
	fake := testutil.NewFakeBBB(t, secret)
	c := bbb.New(fake.BaseURL(), secret, nil)
	room := bbb.Room{ID: "room1", Name: "Weekly sync", ModeratorPassword: "mod", AttendeePassword: "att"}

	for range 2 {
		if err := c.EnsureRoom(context.Background(), room); err != nil {
			t.Fatal(err)
		}
	}

	if got := len(fake.Calls("create")); got != 1 {
		t.Fatalf("%d create calls for two joins, want 1", got)
	}

	p := fake.Calls("create")[0].Params
	if p.Get("meetingID") != "room1" || p.Get("name") != "Weekly sync" || p.Get("moderatorPW") != "mod" || p.Get("attendeePW") != "att" {
		t.Errorf("create parameters %v", p)
	}

	if got := len(fake.Calls("isMeetingRunning")); got != 2 {
		t.Errorf("%d isMeetingRunning calls, want 2", got)
	}
}

func TestEnsureRoomReportsFailures(t *testing.T) {
	fake := testutil.NewFakeBBB(t, secret)
	fake.FailCreate = true

	err := bbb.New(fake.BaseURL(), secret, nil).EnsureRoom(context.Background(), bbb.Room{ID: "r"})
	if err == nil || !strings.Contains(err.Error(), "idNotUnique") {
		t.Errorf("error %v, want the BigBlueButton message key", err)
	}
}

func TestWrongSecretRejected(t *testing.T) {
	fake := testutil.NewFakeBBB(t, secret)

	err := bbb.New(fake.BaseURL(), "other-secret", nil).EnsureRoom(context.Background(), bbb.Room{ID: "r"})
	if err == nil || !strings.Contains(err.Error(), "checksumError") {
		t.Errorf("error %v, want checksumError", err)
	}
}

func TestUnreachableServer(t *testing.T) {
	err := bbb.New("http://127.0.0.1:1/bigbluebutton/", secret, nil).EnsureRoom(context.Background(), bbb.Room{ID: "r"})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestJoinURLHidesSecret(t *testing.T) {
	u, err := url.Parse(bbb.New("https://bbb.example/bigbluebutton/", secret, nil).JoinURL("r", "n", "p"))
	if err != nil || strings.Contains(u.String(), secret) {
		t.Errorf("join URL %s", u)
	}
}
