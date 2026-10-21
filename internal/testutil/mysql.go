package testutil

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"os"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

// MySQL returns a connection to a new, empty database on the server named by
// TEST_MYSQL_DSN, for example root:secret@tcp(127.0.0.1:3307)/. The user
// needs the right to create and drop databases. The test is skipped when the
// variable is not set; the database is dropped when the test ends.
func MySQL(t testing.TB) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("TEST_MYSQL_DSN is not set")
	}

	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}

	suffix := make([]byte, 6)
	if _, err := rand.Read(suffix); err != nil {
		t.Fatal(err)
	}

	name := "irmameet_test_" + hex.EncodeToString(suffix)

	admin, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		t.Fatal(err)
	}

	if _, err := admin.Exec("CREATE DATABASE `" + name + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci"); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		_, _ = admin.Exec("DROP DATABASE `" + name + "`")
		_ = admin.Close()
	})

	cfg.DBName = name
	cfg.ParseTime = true
	cfg.Loc = time.UTC
	cfg.Params = map[string]string{"time_zone": "'+00:00'"}

	connector, err := mysql.NewConnector(cfg)
	if err != nil {
		t.Fatal(err)
	}

	db := sql.OpenDB(connector)
	t.Cleanup(func() { _ = db.Close() })

	return db
}
