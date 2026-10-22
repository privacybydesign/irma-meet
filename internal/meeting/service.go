package meeting

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/mailer"
)

const (
	linkIDBytes   = 16
	roomIDBytes   = 12
	passwordBytes = 16

	joinPath = "/irma_session/join/"
)

// Removed counts what a cleanup deleted.
type Removed struct {
	Meetings     int64
	Participants int64
}

// Store persists meetings.
type Store interface {
	Create(ctx context.Context, m Meeting, participantEmails []string) error
	Find(ctx context.Context, id string) (Meeting, error)
	SetPasswords(ctx context.Context, id, moderator, attendee string) error
	DeleteExpired(ctx context.Context, cutoff time.Time) (Removed, error)
}

// Rooms creates BigBlueButton rooms and builds join URLs.
type Rooms interface {
	EnsureRoom(ctx context.Context, room bbb.Room) error
	JoinURL(roomID, fullName, password string) string
}

// Mails delivers mail.
type Mails interface {
	Send(ctx context.Context, msgs ...mailer.Message) error
}

type Service struct {
	store      Store
	rooms      Rooms
	mails      Mails
	text       *i18n.Translator
	baseURL    string
	nameFormat identity.NameFormat
	now        func() time.Time
}

type Options struct {
	Store      Store
	Rooms      Rooms
	Mails      Mails
	Text       *i18n.Translator
	BaseURL    string
	NameFormat identity.NameFormat
	// Now defaults to time.Now; tests replace it.
	Now func() time.Time
}

func NewService(o Options) *Service {
	now := o.Now
	if now == nil {
		now = time.Now
	}

	return &Service{
		store:      o.Store,
		rooms:      o.Rooms,
		mails:      o.Mails,
		text:       o.Text,
		baseURL:    o.BaseURL,
		nameFormat: o.NameFormat,
		now:        now,
	}
}

// NewMeeting is the input of Create. The host identity comes from a verified
// disclosure, never from form fields.
type NewMeeting struct {
	Type         Type
	Name         string
	Note         string
	Host         identity.Verified
	Participants []string
	Locale       string
}

// Created is the result of Create.
type Created struct {
	Meeting Meeting
	Link    string
	// MailErr is set when the meeting exists but the mails were not delivered.
	MailErr error
}

func (s *Service) Create(ctx context.Context, in NewMeeting) (Created, error) {
	hostName, hostEmail := in.Host.Name(), in.Host.Email()
	if hostName == "" || hostEmail == "" {
		return Created{}, ErrNotVerified
	}

	m := Meeting{Type: in.Type, Name: in.Name, HostName: hostName, HostEmail: hostEmail, InvitationNote: in.Note, StartTime: s.now().UTC().Truncate(time.Second)}

	var err error
	if m.ID, err = randomHex(linkIDBytes); err != nil {
		return Created{}, err
	}

	if m.RoomID, err = randomHex(roomIDBytes); err != nil {
		return Created{}, err
	}

	if m.ModeratorPassword, m.AttendeePassword, err = randomPasswords(); err != nil {
		return Created{}, err
	}

	if err := s.store.Create(ctx, m, in.Participants); err != nil {
		return Created{}, err
	}

	link := s.baseURL + joinPath + m.ID

	return Created{Meeting: m, Link: link, MailErr: s.sendMails(ctx, m, link, in)}, nil
}

func (s *Service) sendMails(ctx context.Context, m Meeting, link string, in NewMeeting) error {
	t := func(key string) string { return s.text.T(in.Locale, key) }

	subject := s.text.Tr(in.Locale, "Yivi Meet link for :name", map[string]string{"name": m.Name})
	noteForHost := t("Please send the above meeting link to people that you wish to invite to the video meeting.")
	msgs := make([]mailer.Message, 0, len(in.Participants)+1)

	for _, email := range in.Participants {
		note := ""
		if in.Note != "" {
			note = t("This mail mostly contains general information, but specifically for this meeting the host has added the following personal message: ") + in.Note
		}

		msgs = append(msgs, mailer.Message{
			Kind: mailer.Invitation, MeetingType: string(m.Type), Locale: in.Locale,
			To: email, ReplyTo: m.HostEmail, Subject: subject,
			Data: mailer.Data{HosterName: m.HostName, InvitationLink: link, InvitationNote: note},
		})

		noteForHost = t("The meeting link has also been sent to the participant.")
		if in.Note != "" {
			noteForHost = t("The meeting link has also been sent to the participant, together with the following personal message: ") + in.Note
		}
	}

	msgs = append(msgs, mailer.Message{
		Kind: mailer.Confirmation, MeetingType: string(m.Type), Locale: in.Locale,
		To: m.HostEmail, Subject: subject,
		Data: mailer.Data{HosterName: m.HostName, InvitationLink: link, InvitationNote: noteForHost},
	})

	return s.mails.Send(ctx, msgs...)
}

// Lookup finds a live meeting. Meetings created before random passwords were
// stored get them now.
func (s *Service) Lookup(ctx context.Context, id string) (Meeting, error) {
	m, err := s.store.Find(ctx, id)
	if err != nil {
		return Meeting{}, err
	}

	if m.Expired(s.now()) {
		return Meeting{}, ErrExpired
	}

	if m.ModeratorPassword != "" && m.AttendeePassword != "" {
		return m, nil
	}

	moderator, attendee, err := randomPasswords()
	if err != nil {
		return Meeting{}, err
	}

	if err := s.store.SetPasswords(ctx, m.ID, moderator, attendee); err != nil {
		return Meeting{}, err
	}

	return s.store.Find(ctx, id)
}

// JoinURL makes sure the room exists and returns where to send who: the host
// gets the moderator password, anyone else the attendee password.
func (s *Service) JoinURL(ctx context.Context, m Meeting, who identity.Verified) (string, error) {
	name := who.BBBName(s.nameFormat)
	if who.Email() == "" || name == "" {
		return "", ErrNotVerified
	}

	room := bbb.Room{ID: m.RoomID, Name: m.Name, ModeratorPassword: m.ModeratorPassword, AttendeePassword: m.AttendeePassword}
	if err := s.rooms.EnsureRoom(ctx, room); err != nil {
		return "", fmt.Errorf("%w: %v", ErrRoom, err)
	}

	password := m.AttendeePassword
	if m.IsHost(who) {
		password = m.ModeratorPassword
	}

	return s.rooms.JoinURL(m.RoomID, name, password), nil
}

// Cleanup deletes meetings whose link has expired, with their participants.
func (s *Service) Cleanup(ctx context.Context) (Removed, error) {
	return s.store.DeleteExpired(ctx, s.now().Add(-LinkValidity))
}

// RunCleanup calls Cleanup every interval until ctx is done.
func (s *Service) RunCleanup(ctx context.Context, interval time.Duration, log *slog.Logger) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		removed, err := s.Cleanup(ctx)
		if err != nil {
			log.Error("cleanup of expired meetings failed", "error", err)
		} else if removed.Meetings > 0 || removed.Participants > 0 {
			log.Info("deleted expired meetings", "meetings", removed.Meetings, "participants", removed.Participants)
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}

func randomPasswords() (moderator, attendee string, err error) {
	if moderator, err = randomHex(passwordBytes); err != nil {
		return "", "", err
	}

	attendee, err = randomHex(passwordBytes)

	return moderator, attendee, err
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
