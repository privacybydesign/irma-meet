// Package store keeps meetings in Postgres.
package store

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"

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

// Open prepares a connection pool for a postgresql:// URL; it does not
// connect yet.
func Open(databaseURL string) (*sql.DB, error) {
	cfg, err := pgx.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL: %w", err)
	}

	if cfg.ConnectTimeout == 0 {
		cfg.ConnectTimeout = connectTimeout
	}

	db := stdlib.OpenDB(*cfg)
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

// Migrate brings the schema up to date. It needs the right to create tables,
// which the server's own database user does not have. An advisory lock keeps
// concurrent runs from migrating at the same time.
func Migrate(ctx context.Context, db *sql.DB) error {
	locker, err := lock.NewPostgresSessionLocker()
	if err != nil {
		return err
	}

	sub, err := fsSub()
	if err != nil {
		return err
	}

	provider, err := goose.NewProvider(goose.DialectPostgres, db, sub, goose.WithSessionLocker(locker))
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
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx,
		`INSERT INTO irma_meet_sessions (`+sessionColumns+`)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		m.ID, string(m.Type), m.Name, m.HostName, m.HostEmail,
		m.InvitationNote, m.RoomID, m.ModeratorPassword, m.AttendeePassword, m.StartTime)
	if err != nil {
		return err
	}

	for _, email := range participantEmails {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO irma_meet_participants (irma_session_id, email_address, authentication)
			VALUES ($1, $2, $3)`,
			m.ID, email, authenticationYivi)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *Store) Find(ctx context.Context, id string) (meeting.Meeting, error) {
	var (
		m           meeting.Meeting
		meetingType string
	)

	err := s.db.QueryRowContext(ctx,
		`SELECT `+sessionColumns+` FROM irma_meet_sessions WHERE irma_session_id = $1`, id).
		Scan(&m.ID, &meetingType, &m.Name, &m.HostName, &m.HostEmail,
			&m.InvitationNote, &m.RoomID, &m.ModeratorPassword, &m.AttendeePassword, &m.StartTime)
	if errors.Is(err, sql.ErrNoRows) {
		return meeting.Meeting{}, meeting.ErrNotFound
	}

	if err != nil {
		return meeting.Meeting{}, err
	}

	m.Type = meeting.Type(meetingType)
	m.StartTime = m.StartTime.UTC()

	return m, nil
}

// SetPasswords fills in passwords that are still empty; existing ones win.
func (s *Store) SetPasswords(ctx context.Context, id, moderator, attendee string) error {
	_, err := s.db.ExecContext(ctx,
		`UPDATE irma_meet_sessions
		SET bbb_moderator_password = COALESCE(NULLIF(bbb_moderator_password, ''), $1),
			bbb_attendee_password = COALESCE(NULLIF(bbb_attendee_password, ''), $2),
			updated_at = now()
		WHERE irma_session_id = $3`,
		moderator, attendee, id)

	return err
}

// DeleteExpired removes meetings that started before cutoff and every
// participant row that no longer belongs to a meeting.
func (s *Store) DeleteExpired(ctx context.Context, cutoff time.Time) (meeting.Removed, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return meeting.Removed{}, err
	}
	defer func() { _ = tx.Rollback() }()

	meetings, err := tx.ExecContext(ctx, `DELETE FROM irma_meet_sessions WHERE start_time < $1`, cutoff)
	if err != nil {
		return meeting.Removed{}, err
	}

	participants, err := tx.ExecContext(ctx,
		`DELETE FROM irma_meet_participants p
		WHERE NOT EXISTS (SELECT 1 FROM irma_meet_sessions s WHERE s.irma_session_id = p.irma_session_id)`)
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
