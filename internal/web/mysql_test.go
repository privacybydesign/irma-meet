package web_test

import (
	"bytes"
	"context"
	"database/sql"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/go-sql-driver/mysql"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/mailer"
	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/store"
	"github.com/privacybydesign/irma-meet/internal/testutil"
	"github.com/privacybydesign/irma-meet/internal/web"
)

// TestEndToEndAgainstMySQL runs a meeting through the real store and the
// MySQL session store. It needs TEST_MYSQL_DSN, a database that may be wiped.
func TestEndToEndAgainstMySQL(t *testing.T) {
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

	if err := store.Migrate(context.Background(), db); err != nil {
		t.Fatal(err)
	}

	irma := testutil.NewFakeIRMA(t, irmaToken)
	fakeBBB := testutil.NewFakeBBB(t, bbbSecret)
	smtp := testutil.NewFakeSMTP(t)

	var handler http.Handler
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { handler.ServeHTTP(w, r) }))
	t.Cleanup(srv.Close)

	text, _ := i18n.New()
	mails, _ := mailer.New(config.Mail{Host: smtp.Host, Port: smtp.Port, FromAddress: "noreply@yivi-meet.test"})

	meetings := meeting.NewService(meeting.Options{
		Store: store.New(db), Rooms: bbb.New(fakeBBB.BaseURL(), bbbSecret, nil), Mails: mails, Text: text,
		BaseURL: srv.URL, NameFormat: identity.NameFormatString,
	})

	sessionStore := mysqlstore.New(db)
	defer sessionStore.StopCleanup()

	key := bytes.Repeat([]byte{7}, config.SessionKeyLen)

	sessions, err := web.NewSessionManager(sessionStore, key, false)
	if err != nil {
		t.Fatal(err)
	}

	appURL, _ := url.Parse(srv.URL)

	handler, err = web.New(web.Options{
		Meetings: meetings, Identities: identity.NewService(irma.URL(), irmaToken, nil), Text: text,
		Sessions: sessions, AppURL: appURL, Log: slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err != nil {
		t.Fatal(err)
	}

	e := &env{t: t, srv: srv, irma: irma, bbb: fakeBBB, smtp: smtp}
	host, guest := e.newUser(), e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "Host@Example.com"))

	if resp := host.post("/irma_session/store", createForm("free", "Weekly sync")); resp.StatusCode != http.StatusOK {
		t.Fatalf("store: status %d: %s", resp.StatusCode, resp.Body)
	}

	var id string
	if err := db.QueryRow("SELECT irma_session_id FROM irma_meet_sessions WHERE meeting_name = 'Weekly sync'").Scan(&id); err != nil {
		t.Fatal(err)
	}

	if resp := host.get("/irma_session/join/" + id); resp.StatusCode != http.StatusFound {
		t.Errorf("host join: status %d: %s", resp.StatusCode, resp.Body)
	}

	guest.disclose("default", defaultAttrs("Gerrit Guest", "guest@example.com"))

	if resp := guest.get("/irma_session/join_participant/" + id); resp.StatusCode != http.StatusFound {
		t.Errorf("guest join: status %d: %s", resp.StatusCode, resp.Body)
	}

	var stored int
	if err := db.QueryRow("SELECT COUNT(*) FROM sessions").Scan(&stored); err != nil || stored < 2 {
		t.Fatalf("%d sessions stored, %v", stored, err)
	}

	rows, err := db.Query("SELECT data FROM sessions")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		if err := rows.Scan(&data); err != nil {
			t.Fatal(err)
		}

		for _, secret := range []string{"Hannah Host", "host@example.com", "Gerrit", "guest@example.com", "pbdf.pbdf.email"} {
			if strings.Contains(strings.ToLower(string(data)), strings.ToLower(secret)) {
				t.Errorf("a stored session contains %q in plain text", secret)
			}
		}
	}
}
