// Package meeting holds the domain of a Yivi Meet meeting: creating one,
// letting verified people into its BigBlueButton room, and expiring it.
package meeting

import (
	"errors"
	"strings"
	"time"

	"github.com/privacybydesign/irma-meet/internal/identity"
)

// LinkValidity is how long a meeting link works after it was created.
const LinkValidity = 24 * time.Hour

var (
	ErrNotFound    = errors.New("meeting not found")
	ErrExpired     = errors.New("meeting link expired")
	ErrNotVerified = errors.New("disclosed attributes are incomplete")
	ErrRoom        = errors.New("meeting room unavailable")
)

// Type is the kind of meeting; it decides who has to disclose what.
type Type string

const (
	Free Type = "free"
	Exam Type = "exam"
)

func ParseType(name string) (Type, bool) {
	t := Type(name)

	return t, t == Free || t == Exam
}

// HostIdentity is the disclosure the host of this meeting type has to do.
func (t Type) HostIdentity() identity.Type {
	if t == Exam {
		return identity.TypeTeacher
	}

	return identity.TypeDefault
}

// ParticipantIdentity is the disclosure everyone else has to do.
func (t Type) ParticipantIdentity() identity.Type {
	if t == Exam {
		return identity.TypeStudent
	}

	return identity.TypeDefault
}

// Meeting is a stored meeting.
type Meeting struct {
	ID                string
	Type              Type
	Name              string
	HostName          string
	HostEmail         string
	InvitationNote    string
	RoomID            string
	ModeratorPassword string
	AttendeePassword  string
	StartTime         time.Time
}

func (m Meeting) Expired(now time.Time) bool {
	return !now.Before(m.StartTime.Add(LinkValidity))
}

// IsHost reports whether who disclosed as the host type with the host's email.
func (m Meeting) IsHost(who identity.Verified) bool {
	email := who.Email()

	return who.Type == m.Type.HostIdentity() && email != "" && strings.EqualFold(email, m.HostEmail)
}
