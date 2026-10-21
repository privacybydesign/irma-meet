// Package store keeps meetings in MySQL.
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"

	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/meeting"
)

const (
	connectTimeout  = 10 * time.Second
	connMaxLifetime = 3 * time.Minute
	maxOpenConns    = 10

	// authenticationYivi is the value the participants table has always held.
	authenticationYivi = 1

	pingInterval = 2 * time.Second
)

//go:embed migrations/*.sql
var migrations embed.FS

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

// Open prepares a connection pool; it does not connect yet.
func Open(cfg config.Database) (*sql.DB, error) {
	c := mysql.NewConfig()
	c.User = cfg.User
	c.Passwd = cfg.Password
	c.Net = "tcp"
	c.Addr = net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
	c.DBName = cfg.Name
	c.ParseTime = true
	c.Loc = time.UTC
	c.Timeout = connectTimeout
	c.Collation = "utf8mb4_unicode_ci"
	c.Params = map[string]string{"time_zone": "'+00:00'"}

	connector, err := mysql.NewConnector(c)
	if err != nil {
		return nil, err
	}

	db := sql.OpenDB(connector)
	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)

	return db, nil
}

// WaitReady pings the database until it answers or ctx is done.
func WaitReady(ctx context.Context, db *sql.DB) error {
	for {
		err := db.PingContext(ctx)
		if err == nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return fmt.Errorf("database not reachable: %w", errors.Join(err, ctx.Err()))
		case <-time.After(pingInterval):
		}
	}
}

// Migrate brings the schema up to date. A table lock keeps replicas that
// start together from migrating at the same time.
func Migrate(ctx context.Context, db *sql.DB) error {
	locker, err := lock.NewMySQLTableLocker()
	if err != nil {
		return err
	}

	sub, err := fsSub()
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(goose.DialectMySQL, db, sub, goose.WithLocker(locker))
	if err != nil {
		return err
	}

	_, err = provider.Up(ctx)

	return err
}

const sessionColumns = `irma_session_id, meeting_type, meeting_name, hoster_name, hoster_email_address,
	invitation_note, bbb_session_id, bbb_moderator_password, bbb_attendee_password, start_time`

func (s *Store) Create(ctx context.Context, m meeting.Meeting, participantEmails []string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now().UTC()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO irma_meet_sessions (`+sessionColumns+`, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID, string(m.Type), m.Name, m.HostName, m.HostEmail,
		m.InvitationNote, m.RoomID, m.ModeratorPassword, m.AttendeePassword, m.StartTime, now, now)
	if err != nil {
		return err
	}

	for _, email := range participantEmails {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO irma_meet_participants (irma_session_id, email_address, authentication, created_at, updated_at)
			VALUES (?, ?, ?, ?, ?)`,
			m.ID, email, authenticationYivi, now, now)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Find(ctx context.Context, id string) (meeting.Meeting, error) {
	var (
		m                   meeting.Meeting
		meetingType         string
		note                sql.NullString
		moderator, attendee sql.NullString
	)

	err := s.db.QueryRowContext(ctx,
		`SELECT `+sessionColumns+` FROM irma_meet_sessions WHERE irma_session_id = ?`, id).
		Scan(&m.ID, &meetingType, &m.Name, &m.HostName, &m.HostEmail,
			&note, &m.RoomID, &moderator, &attendee, &m.StartTime)
	if errors.Is(err, sql.ErrNoRows) {
		return meeting.Meeting{}, meeting.ErrNotFound
	}

	if err != nil {
		return meeting.Meeting{}, err
	}

	m.Type = meeting.Type(meetingType)
	m.InvitationNote = note.String
	m.ModeratorPassword = moderator.String
	m.AttendeePassword = attendee.String

	return m, nil
}

// SetPasswords fills in passwords that are still empty; existing ones win.
func (s *Store) SetPasswords(ctx context.Context, id, moderator, attendee string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE irma_meet_sessions
		SET bbb_moderator_password = IF(COALESCE(bbb_moderator_password, '') = '', ?, bbb_moderator_password),
			bbb_attendee_password = IF(COALESCE(bbb_attendee_password, '') = '', ?, bbb_attendee_password),
			updated_at = ?
		WHERE irma_session_id = ?`,
		moderator, attendee, time.Now().UTC(), id)

	return err
}

// DeleteExpired removes meetings that started before cutoff and every
// participant row that no longer belongs to a meeting.
func (s *Store) DeleteExpired(ctx context.Context, cutoff time.Time) (meeting.Removed, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return meeting.Removed{}, err
	}
	defer tx.Rollback()

	meetings, err := tx.ExecContext(ctx, `DELETE FROM irma_meet_sessions WHERE start_time < ?`, cutoff.UTC())
	if err != nil {
		return meeting.Removed{}, err
	}

	participants, err := tx.ExecContext(ctx,
		`DELETE p FROM irma_meet_participants p
		LEFT JOIN irma_meet_sessions s ON s.irma_session_id = p.irma_session_id
		WHERE s.id IS NULL`)
	if err != nil {
		return meeting.Removed{}, err
	}

	var removed meeting.Removed
	if removed.Meetings, err = meetings.RowsAffected(); err != nil {
		return meeting.Removed{}, err
	}

	if removed.Participants, err = participants.RowsAffected(); err != nil {
		return meeting.Removed{}, err
	}

	return removed, tx.Commit()
}

func fsSub() (fs.FS, error) {
	return fs.Sub(migrations, "migrations")
}
