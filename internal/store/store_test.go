package store_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/store"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

// The tests need a Postgres server, see testutil.Postgres.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	return testutil.Postgres(t)
}

func migrated(t *testing.T) (*sql.DB, *store.Store) {
	t.Helper()

	db := openTestDB(t)
	if err := store.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}

	return db, store.New(db)
}

func sample(id string, start time.Time) meeting.Meeting {
	return meeting.Meeting{
		ID: id, Type: meeting.Exam, Name: "Oral exam ünïcode", HostName: "Tanja Teacher", HostEmail: "tanja@uni.example",
		InvitationNote: "Good luck", RoomID: "room-" + id, ModeratorPassword: "mod-" + id, AttendeePassword: "att-" + id, StartTime: start,
	}
}

func TestMigrateIsIdempotent(t *testing.T) {
	db := openTestDB(t)

	for range 2 {
		if err := store.Migrate(context.Background(), db); err != nil {
			t.Fatal(err)
		}
	}

	for _, table := range []string{"irma_meet_sessions", "irma_meet_participants", "sessions"} {
		var n int
		if err := db.QueryRow(`SELECT COUNT(*) FROM ` + table).Scan(&n); err != nil {
			t.Errorf("%s: %v", table, err)
		}
	}
}

func TestCreateAndFind(t *testing.T) {
	_, s := migrated(t)
	ctx := context.Background()
	start := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

	if err := s.Create(ctx, sample("aaa111", start), []string{"student@uni.example", "other@uni.example"}); err != nil {
		t.Fatal(err)
	}

	got, err := s.Find(ctx, "aaa111")
	if err != nil {
		t.Fatal(err)
	}

	if got != sample("aaa111", start) {
		t.Errorf("round trip\n got %+v\nwant %+v", got, sample("aaa111", start))
	}

	if _, err := s.Find(ctx, "missing"); err != meeting.ErrNotFound {
		t.Errorf("unknown id: %v, want ErrNotFound", err)
	}
}

func TestCreateIsAtomic(t *testing.T) {
	db, s := migrated(t)
	ctx := context.Background()

	if err := s.Create(ctx, sample("dup", time.Now()), []string{"a@example.com"}); err != nil {
		t.Fatal(err)
	}

	if err := s.Create(ctx, sample("dup", time.Now()), []string{"b@example.com"}); err == nil {
		t.Fatal("a duplicate link id must be rejected")
	}

	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM irma_meet_participants").Scan(&n); err != nil || n != 1 {
		t.Errorf("%d participant rows, %v, want only the first meeting's", n, err)
	}
}

func TestSetPasswordsFillsEmpty(t *testing.T) {
	db, s := migrated(t)
	ctx := context.Background()

	if _, err := db.Exec("INSERT INTO irma_meet_sessions (irma_session_id, meeting_name, hoster_name, hoster_email_address, start_time, bbb_session_id, bbb_moderator_password, meeting_type) VALUES ('legacy', 'n', 'h', 'h@example.com', now(), 'r', 'kept', 'free')"); err != nil {
		t.Fatal(err)
	}

	if err := s.SetPasswords(ctx, "legacy", "new-mod", "new-att"); err != nil {
		t.Fatal(err)
	}

	m, _ := s.Find(ctx, "legacy")
	if m.ModeratorPassword != "kept" || m.AttendeePassword != "new-att" {
		t.Errorf("passwords %q %q, want the existing moderator password kept and the attendee password filled", m.ModeratorPassword, m.AttendeePassword)
	}
}

func TestDeleteExpired(t *testing.T) {
	db, s := migrated(t)
	ctx := context.Background()
	now := time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

	for _, m := range []struct {
		meeting meeting.Meeting
		emails  []string
	}{
		{sample("old1", now.Add(-25*time.Hour)), []string{"a@example.com", "b@example.com"}},
		{sample("old2", now.Add(-48*time.Hour)), nil},
		{sample("fresh", now.Add(-time.Hour)), []string{"c@example.com"}},
	} {
		if err := s.Create(ctx, m.meeting, m.emails); err != nil {
			t.Fatal(err)
		}
	}

	if _, err := db.Exec("INSERT INTO irma_meet_participants (irma_session_id, email_address, authentication) VALUES ('gone-long-ago', 'orphan@example.com', 1)"); err != nil {
		t.Fatal(err)
	}

	removed, err := s.DeleteExpired(ctx, now.Add(-meeting.LinkValidity))
	if err != nil {
		t.Fatal(err)
	}

	if removed.Meetings != 2 || removed.Participants != 3 {
		t.Errorf("removed %+v, want 2 meetings and 3 participants (two of old1, one orphan)", removed)
	}

	if _, err := s.Find(ctx, "fresh"); err != nil {
		t.Errorf("the fresh meeting was deleted: %v", err)
	}

	var n int
	if err := db.QueryRow("SELECT COUNT(*) FROM irma_meet_participants").Scan(&n); err != nil || n != 1 {
		t.Errorf("%d participant rows left, %v, want 1", n, err)
	}
}
