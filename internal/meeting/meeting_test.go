package meeting_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/mailer"
	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

const (
	attrEmail     = "pbdf.pbdf.email.email"
	attrFullName  = "pbdf.gemeente.personalData.fullname"
	attrSurfEmail = "pbdf.pbdf.surfnet-2.email"
)

var t0 = time.Date(2026, 10, 21, 12, 0, 0, 0, time.UTC)

type fakeRooms struct {
	ensured []bbb.Room
	err     error
}

func (f *fakeRooms) EnsureRoom(_ context.Context, r bbb.Room) error {
	f.ensured = append(f.ensured, r)

	return f.err
}

func (f *fakeRooms) JoinURL(roomID, fullName, password string) string {
	return "join://" + roomID + "/" + fullName + "/" + password
}

type fakeMails struct{ sent []mailer.Message }

func (f *fakeMails) Send(_ context.Context, msgs ...mailer.Message) error {
	f.sent = append(f.sent, msgs...)

	return nil
}

type fixture struct {
	svc   *meeting.Service
	store *testutil.MemStore
	rooms *fakeRooms
	mails *fakeMails
	now   time.Time
}

func newFixture(t *testing.T) *fixture {
	t.Helper()

	text, err := i18n.New()
	if err != nil {
		t.Fatal(err)
	}

	f := &fixture{store: testutil.NewMemStore(), rooms: &fakeRooms{}, mails: &fakeMails{}, now: t0}
	f.svc = meeting.NewService(meeting.Options{
		Store: f.store, Rooms: f.rooms, Mails: f.mails, Text: text, BaseURL: "https://meet.example",
		NameFormat: identity.NameFormatString, Now: func() time.Time { return f.now },
	})

	return f
}

func host() identity.Verified {
	return identity.Verified{Type: identity.TypeDefault, Attrs: map[string]string{
		attrFullName: "Hannah Host", attrEmail: "Host@Example.com",
	}}
}

func TestExpiryBoundary(t *testing.T) {
	m := meeting.Meeting{StartTime: t0}

	if m.Expired(t0.Add(meeting.LinkValidity - time.Nanosecond)) {
		t.Error("expired one nanosecond early")
	}

	if !m.Expired(t0.Add(meeting.LinkValidity)) {
		t.Error("not expired after 24 hours")
	}
}

func TestHostRoleNeedsTypeAndEmail(t *testing.T) {
	free := meeting.Meeting{Type: meeting.Free, HostEmail: "host@example.com"}
	exam := meeting.Meeting{Type: meeting.Exam, HostEmail: "tanja@uni.example"}

	attrs := func(typ identity.Type, id, email string) identity.Verified {
		return identity.Verified{Type: typ, Attrs: map[string]string{id: email}}
	}

	cases := []struct {
		name string
		m    meeting.Meeting
		who  identity.Verified
		want bool
	}{
		{"host", free, attrs(identity.TypeDefault, attrEmail, "HOST@example.com"), true},
		{"other person", free, attrs(identity.TypeDefault, attrEmail, "guest@example.com"), false},
		{"nobody", free, identity.Verified{}, false},
		{"teacher", exam, attrs(identity.TypeTeacher, attrSurfEmail, "tanja@uni.example"), true},
		{"student with the teacher's email", exam, attrs(identity.TypeStudent, attrSurfEmail, "tanja@uni.example"), false},
	}

	for _, c := range cases {
		if got := c.m.IsHost(c.who); got != c.want {
			t.Errorf("%s: IsHost = %v, want %v", c.name, got, c.want)
		}
	}
}

func TestMeetingTypes(t *testing.T) {
	for name, wantHost := range map[string]identity.Type{"free": identity.TypeDefault, "exam": identity.TypeTeacher} {
		mt, ok := meeting.ParseType(name)
		if !ok || mt.HostIdentity() != wantHost {
			t.Errorf("%s: %v %v", name, mt, ok)
		}
	}

	if meeting.Exam.ParticipantIdentity() != identity.TypeStudent || meeting.Free.ParticipantIdentity() != identity.TypeDefault {
		t.Error("wrong participant identity")
	}

	for _, name := range []string{"medical_consult", "", "evil", "Free"} {
		if _, ok := meeting.ParseType(name); ok {
			t.Errorf("%q must not be a meeting type", name)
		}
	}
}

func TestCreateStoresAndMails(t *testing.T) {
	f := newFixture(t)

	created, err := f.svc.Create(context.Background(), meeting.NewMeeting{
		Type: meeting.Free, Name: "Weekly sync", Note: "hello", Host: host(), Participants: []string{"a@example.com", "b@example.com"}, Locale: i18n.English,
	})
	if err != nil {
		t.Fatal(err)
	}

	m := created.Meeting
	if !strings.HasPrefix(created.Link, "https://meet.example/irma_session/join/") || !strings.HasSuffix(created.Link, m.ID) {
		t.Errorf("link %q", created.Link)
	}

	if m.HostEmail != "host@example.com" || m.HostName != "Hannah Host" || !m.StartTime.Equal(t0) {
		t.Errorf("meeting %+v", m)
	}

	if got := f.store.Participants(m.ID); len(got) != 2 {
		t.Errorf("participants %v", got)
	}

	if len(f.mails.sent) != 3 {
		t.Fatalf("%d mails, want 2 invitations and 1 confirmation", len(f.mails.sent))
	}

	if f.mails.sent[0].Kind != mailer.Invitation || f.mails.sent[1].Kind != mailer.Invitation || f.mails.sent[2].Kind != mailer.Confirmation {
		t.Error("invitations go before the confirmation")
	}

	confirmation := f.mails.sent[2]
	if confirmation.To != "host@example.com" || !strings.Contains(confirmation.Data.InvitationNote, "together with the following personal message: hello") {
		t.Errorf("confirmation %+v", confirmation)
	}

	if f.mails.sent[0].ReplyTo != "host@example.com" || !strings.Contains(f.mails.sent[0].Data.InvitationNote, "hello") {
		t.Errorf("invitation %+v", f.mails.sent[0])
	}

	if confirmation.ReplyTo != "" {
		t.Errorf("the confirmation has Reply-To %q, replies must go to the sender address", confirmation.ReplyTo)
	}
}

func TestCreateWithoutGuests(t *testing.T) {
	f := newFixture(t)

	if _, err := f.svc.Create(context.Background(), meeting.NewMeeting{Type: meeting.Free, Name: "n", Host: host(), Locale: i18n.English}); err != nil {
		t.Fatal(err)
	}

	if len(f.mails.sent) != 1 || !strings.Contains(f.mails.sent[0].Data.InvitationNote, "Please send the above meeting link") {
		t.Errorf("mails %+v", f.mails.sent)
	}
}

func TestCreateNeedsFullHost(t *testing.T) {
	f := newFixture(t)

	for _, who := range []identity.Verified{{}, {Type: identity.TypeDefault, Attrs: map[string]string{attrEmail: "a@example.com"}}} {
		_, err := f.svc.Create(context.Background(), meeting.NewMeeting{Type: meeting.Free, Name: "n", Host: who})
		if !errors.Is(err, meeting.ErrNotVerified) {
			t.Errorf("%+v: %v, want ErrNotVerified", who, err)
		}
	}

	if len(f.store.Meetings()) != 0 {
		t.Error("a meeting without a verified host was stored")
	}
}

// Meetings that predate the stored passwords get independent random ones, once.
func TestLookupBackfillsPasswords(t *testing.T) {
	f := newFixture(t)

	for _, id := range []string{"aaaa1111", "bbbb2222"} {
		f.store.Insert(meeting.Meeting{ID: id, Type: meeting.Free, Name: "Test meeting", HostName: "Host", HostEmail: "host@example.com", RoomID: "room-" + id, StartTime: t0})
	}

	a, err := f.svc.Lookup(context.Background(), "aaaa1111")
	if err != nil {
		t.Fatal(err)
	}

	b, _ := f.svc.Lookup(context.Background(), "bbbb2222")

	distinct := map[string]bool{a.ModeratorPassword: true, a.AttendeePassword: true, b.ModeratorPassword: true, b.AttendeePassword: true}
	if len(distinct) != 4 || a.ModeratorPassword == "" || len(a.ModeratorPassword) != 32 {
		t.Errorf("passwords are not random and independent: %v", distinct)
	}

	for _, p := range []string{a.ModeratorPassword, a.AttendeePassword} {
		if p == "hoster"+a.RoomID || strings.Contains(p, a.RoomID) {
			t.Errorf("password %q derives from the room id", p)
		}
	}

	again, _ := f.svc.Lookup(context.Background(), "aaaa1111")
	if again.ModeratorPassword != a.ModeratorPassword || again.AttendeePassword != a.AttendeePassword {
		t.Error("the passwords were not persisted")
	}
}

func TestLookupErrors(t *testing.T) {
	f := newFixture(t)
	f.store.Insert(meeting.Meeting{ID: "abc", Type: meeting.Free, StartTime: t0, ModeratorPassword: "m", AttendeePassword: "a"})

	if _, err := f.svc.Lookup(context.Background(), "nope"); !errors.Is(err, meeting.ErrNotFound) {
		t.Errorf("unknown: %v", err)
	}

	f.now = t0.Add(meeting.LinkValidity)

	if _, err := f.svc.Lookup(context.Background(), "abc"); !errors.Is(err, meeting.ErrExpired) {
		t.Errorf("expired: %v", err)
	}
}

func TestJoinURLPasswordByRole(t *testing.T) {
	f := newFixture(t)
	created, _ := f.svc.Create(context.Background(), meeting.NewMeeting{Type: meeting.Free, Name: "Weekly sync", Host: host(), Locale: i18n.English})
	m := created.Meeting

	asHost, err := f.svc.JoinURL(context.Background(), m, host())
	if err != nil || asHost != "join://"+m.RoomID+"/Hannah Host/"+m.ModeratorPassword {
		t.Errorf("host: %q %v", asHost, err)
	}

	guest := identity.Verified{Type: identity.TypeDefault, Attrs: map[string]string{attrFullName: "Gerrit", attrEmail: "guest@example.com"}}

	asGuest, err := f.svc.JoinURL(context.Background(), m, guest)
	if err != nil || asGuest != "join://"+m.RoomID+"/Gerrit/"+m.AttendeePassword {
		t.Errorf("guest: %q %v", asGuest, err)
	}

	if len(f.rooms.ensured) != 2 || f.rooms.ensured[0].Name != "Weekly sync" || f.rooms.ensured[0].ModeratorPassword != m.ModeratorPassword {
		t.Errorf("rooms %+v", f.rooms.ensured)
	}
}

func TestJoinURLErrors(t *testing.T) {
	f := newFixture(t)
	created, _ := f.svc.Create(context.Background(), meeting.NewMeeting{Type: meeting.Free, Name: "n", Host: host(), Locale: i18n.English})

	incomplete := identity.Verified{Type: identity.TypeDefault, Attrs: map[string]string{attrEmail: "guest@example.com"}}
	if _, err := f.svc.JoinURL(context.Background(), created.Meeting, incomplete); !errors.Is(err, meeting.ErrNotVerified) {
		t.Errorf("incomplete identity: %v", err)
	}

	if len(f.rooms.ensured) != 0 {
		t.Error("a room was created for an unverified person")
	}

	f.rooms.err = errors.New("boom")

	if _, err := f.svc.JoinURL(context.Background(), created.Meeting, host()); !errors.Is(err, meeting.ErrRoom) {
		t.Errorf("room failure: %v", err)
	}
}

func TestCleanupDeletesExpired(t *testing.T) {
	f := newFixture(t)
	f.store.Insert(meeting.Meeting{ID: "old", StartTime: t0.Add(-25 * time.Hour)}, "a@example.com")
	f.store.Insert(meeting.Meeting{ID: "new", StartTime: t0.Add(-23 * time.Hour)}, "b@example.com")

	removed, err := f.svc.Cleanup(context.Background())
	if err != nil || removed.Meetings != 1 || removed.Participants != 1 {
		t.Errorf("removed %+v, %v", removed, err)
	}

	if got := f.store.Meetings(); len(got) != 1 || got[0].ID != "new" {
		t.Errorf("remaining %+v", got)
	}
}
