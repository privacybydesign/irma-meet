package testutil

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/url"
	"os"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/store"
)

// Postgres returns a connection to a new, empty database, see PostgresURL.
func Postgres(t testing.TB) *sql.DB {
	t.Helper()

	db, err := store.Open(PostgresURL(t))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { _ = db.Close() })

	return db
}

// PostgresURL creates a new, empty database on the server named by
// TEST_DATABASE_URL, for example
// postgresql://postgres:secret@127.0.0.1:5432/postgres, and returns its URL.
// The user needs the right to create and drop databases. The test is skipped
// when the variable is not set; the database is dropped when the test ends.
func PostgresURL(t testing.TB) string {
	t.Helper()

	raw := os.Getenv("TEST_DATABASE_URL")
	if raw == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatal(err)
	}

	name := "irmameet_test_" + hex.EncodeToString(suffix)

	admin, err := store.Open(raw)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := admin.Exec(`CREATE DATABASE "` + name + `"`); err != nil {
		t.Fatal(err)
	}

	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}

	u.Path = "/" + name

	t.Cleanup(func() {
		_, _ = admin.ExecContext(context.Background(), `DROP DATABASE "`+name+`" WITH (FORCE)`)
		_ = admin.Close()
	})

	return u.String()
}
