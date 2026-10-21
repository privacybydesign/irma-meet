package store_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"

	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/store"
)

// These tests need a MySQL database that they may wipe, e.g.
// TEST_MYSQL_DSN='root:secret@tcp(127.0.0.1:3307)/irmameet_test'.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN is not set")
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}

	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.Params = map[string]string{"time_zone": "'+00:00'"}

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		t.Fatal(err)
	}

	db := sql.OpenDB(connector)
	t.Cleanup(func() { _ = db.Close() })

	for _, table := range []string{"irma_meet_participants", "irma_meet_sessions", "sessions", "goose_db_version", "goose_lock"} {
		if _, err := db.Exec("DROP TABLE IF EXISTS `" + table + "`"); err != nil {
			t.Fatal(err)
		}
	}

	return db
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

func TestMigrateCreatesTheTablesAndIsIdempotent(t *testing.T) {
	db := openTestDB(t)

	for range 2 {
		if err := store.Migrate(context.Background(), db); err != nil {
			t.Fatal(err)
		}
	}

	for _, table := range []string{"irma_meet_sessions", "irma_meet_participants", "sessions"} {
		var n int
		if err := db.QueryRow("SELECT COUNT(*) FROM `" + table + "`").Scan(&n); err != nil {
			t.Errorf("%s: %v", table, err)
		}
	}
}

// A database that the Laravel app filled keeps its rows.
func TestMigrateReusesExistingLaravelTables(t *testing.T) {
	db := openTestDB(t)

	for _, ddl := range []string{
		"CREATE TABLE `irma_meet_sessions` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `irma_session_id` char(48) NOT NULL, `meeting_name` varchar(255) NOT NULL, `hoster_name` varchar(255) NOT NULL, `hoster_email_address` varchar(255) NOT NULL, `start_time` datetime NOT NULL, `invitation_note` varchar(255) DEFAULT '', `bbb_session_id` varchar(255) NOT NULL, `bbb_moderator_password` varchar(255) DEFAULT NULL, `bbb_attendee_password` varchar(255) DEFAULT NULL, `created_at` timestamp NULL DEFAULT NULL, `updated_at` timestamp NULL DEFAULT NULL, `meeting_type` varchar(255) NOT NULL, PRIMARY KEY (`id`), UNIQUE KEY `irma_meet_sessions_irma_session_id_unique` (`irma_session_id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"CREATE TABLE `irma_meet_participants` (`id` bigint unsigned NOT NULL AUTO_INCREMENT, `irma_session_id` char(48) NOT NULL, `email_address` varchar(255) NOT NULL, `authentication` int NOT NULL, `created_at` timestamp NULL DEFAULT NULL, `updated_at` timestamp NULL DEFAULT NULL, PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci",
		"INSERT INTO irma_meet_sessions (irma_session_id, meeting_name, hoster_name, hoster_email_address, start_time, invitation_note, bbb_session_id, meeting_type) VALUES ('legacy1', 'Old', 'Host', 'host@example.com', '2026-10-20 10:00:00', NULL, 'bbb-legacy', 'free')",
	} {
		if _, err := db.Exec(ddl); err != nil {
			t.Fatal(err)
		}
	}

	if err := store.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}

	m, err := store.New(db).Find(context.Background(), "legacy1")
	if err != nil {
		t.Fatal(err)
	}

	if m.Name != "Old" || m.Type != meeting.Free || m.RoomID != "bbb-legacy" || m.ModeratorPassword != "" || m.InvitationNote != "" {
		t.Errorf("legacy meeting %+v", m)
	}

	if want := time.Date(2026, 10, 20, 10, 0, 0, 0, time.UTC); !m.StartTime.Equal(want) {
		t.Errorf("start time %s, want %s", m.StartTime, want)
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

func TestSetPasswordsOnlyFillsEmptyOnes(t *testing.T) {
	db, s := migrated(t)
	ctx := context.Background()

	if _, err := db.Exec("INSERT INTO irma_meet_sessions (irma_session_id, meeting_name, hoster_name, hoster_email_address, start_time, bbb_session_id, bbb_moderator_password, meeting_type) VALUES ('legacy', 'n', 'h', 'h@example.com', NOW(), 'r', 'kept', 'free')"); err != nil {
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

	s.Create(ctx, sample("old1", now.Add(-25*time.Hour)), []string{"a@example.com", "b@example.com"})
	s.Create(ctx, sample("old2", now.Add(-48*time.Hour)), nil)
	s.Create(ctx, sample("fresh", now.Add(-time.Hour)), []string{"c@example.com"})

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
