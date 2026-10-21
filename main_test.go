package main

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/store"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

func quietLog() *slog.Logger {
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func TestUnknownCommand(t *testing.T) {
	for _, args := range [][]string{{"serve"}, {"migrate", "now"}} {
		err := run(quietLog(), args)
		if err == nil || !strings.Contains(err.Error(), "usage") {
			t.Errorf("args %v: %v, want a usage error", args, err)
		}
	}
}

func TestMigrateNeedsDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	err := run(quietLog(), []string{commandMigrate})
	if err == nil || !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Errorf("got %v, want an error naming DATABASE_URL", err)
	}
}

// The migrate command reads DATABASE_URL and nothing else, so an init
// container does not need the other secrets.
func TestMigrateCommand(t *testing.T) {
	databaseURL := testutil.PostgresURL(t)
	t.Setenv("DATABASE_URL", databaseURL)

	for range 2 {
		if err := run(quietLog(), []string{commandMigrate}); err != nil {
			t.Fatal(err)
		}
	}

	db, err := store.Open(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = db.Close() }()

	for _, table := range []string{"irma_meet_sessions", "irma_meet_participants", "sessions"} {
		var n int
		if err := db.QueryRowContext(context.Background(), `SELECT COUNT(*) FROM `+table).Scan(&n); err != nil {
			t.Errorf("%s: %v", table, err)
		}
	}
}
