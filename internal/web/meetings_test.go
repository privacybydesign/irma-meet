package web_test

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/privacybydesign/irma-meet/internal/meeting"
)

var linkPattern = regexp.MustCompile(`href="([^"]*/irma_session/join/([0-9a-f]+))"`)

func createForm(meetingType, name string) url.Values {
	return url.Values{
		"meeting_type": {meetingType},
		"meeting_name": {name},
		"agreed":       {"yes"},
	}
}

// newFreeMeeting has a host disclose and create a free meeting. It returns the
// meeting as stored.
func (e *env) newFreeMeeting(host *user, name string) meeting.Meeting {
	e.t.Helper()

	host.disclose("default", defaultAttrs("Hannah Host", "Host@Example.com"))

	resp := host.post("/irma_session/store", createForm("free", name))
	if resp.StatusCode != http.StatusOK {
		e.t.Fatalf("store: status %d: %s", resp.StatusCode, resp.Body)
	}

	for _, m := range e.store.Meetings() {
		if m.Name == name {
			return m
		}
	}

	e.t.Fatal("meeting was not stored")

	return meeting.Meeting{}
}

func TestHostCreatesFreeMeeting(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "Host@Example.com"))

	form := host.get("/irma_session/create/free")
	if form.StatusCode != http.StatusOK {
		t.Fatalf("create form: status %d", form.StatusCode)
	}

	for _, want := range []string{`value="Hannah Host"`, `value="host@example.com"`, `readonly`} {
		if !strings.Contains(form.Body, want) {
			t.Errorf("create form misses %q", want)
		}
	}

	done := host.post("/irma_session/store", createForm("free", "Weekly sync"))
	if done.StatusCode != http.StatusOK {
		t.Fatalf("store: status %d: %s", done.StatusCode, done.Body)
	}

	link := linkPattern.FindStringSubmatch(done.Body)
	if link == nil {
		t.Fatalf("success page has no meeting link: %s", done.Body)
	}

	stored := e.store.Meetings()
	if len(stored) != 1 {
		t.Fatalf("stored %d meetings, want 1", len(stored))
	}

	m := stored[0]
	if link[2] != m.ID || len(m.ID) != 32 {
		t.Errorf("link id %q, stored id %q: want the same 32 hex characters", link[2], m.ID)
	}

	if m.HostName != "Hannah Host" || m.HostEmail != "host@example.com" || m.Type != meeting.Free {
		t.Errorf("stored meeting has wrong host or type: %+v", m)
	}

	if e.irma.StartedCount() != 1 {
		t.Errorf("%d Yivi sessions started for one disclosure, want 1", e.irma.StartedCount())
	}
}

func TestHostGetsConfirmationMail(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	mails := e.smtp.To("host@example.com")
	if len(mails) != 1 {
		t.Fatalf("host got %d mails, want 1", len(mails))
	}

	mail := mails[0]
	link := e.srv.URL + "/irma_session/join/" + m.ID

	if got := mail.Header(headerSubject); got != "Yivi Meet link voor Weekly sync" {
		t.Errorf("subject %q", got)
	}

	if !strings.Contains(mail.Body(), link) {
		t.Errorf("mail does not contain the link %s:\n%s", link, mail.Body())
	}

	if !strings.Contains(mail.Body(), "Hannah Host") {
		t.Errorf("mail does not greet the host:\n%s", mail.Body())
	}
}

func TestMailsFollowLanguage(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.get("/lang/en")
	e.newFreeMeeting(host, "Weekly sync")

	mail := e.smtp.To("host@example.com")[0]

	if got := mail.Header(headerSubject); got != "Yivi Meet link for Weekly sync" {
		t.Errorf("subject %q", got)
	}

	if !strings.Contains(mail.Body(), "This mail confirms that you have made a reservation") {
		t.Errorf("mail is not English:\n%s", mail.Body())
	}
}

func TestHostCreatesExamMeeting(t *testing.T) {
	e := newEnv(t)
	teacher := e.newUser()

	teacher.disclose("teacher", teacherAttrs("Tanja Teacher", "tanja@uni.example"))

	form := teacher.get("/irma_session/create/exam")
	if form.StatusCode != http.StatusOK {
		t.Fatalf("exam create form: status %d, want 200 (was a 500 in the PHP app)", form.StatusCode)
	}

	if !strings.Contains(form.Body, `value="Tanja Teacher"`) || !strings.Contains(form.Body, `value="tanja@uni.example"`) {
		t.Errorf("exam form is not prefilled with the teacher:\n%s", form.Body)
	}

	post := createForm("exam", "Oral exam")
	post.Set("invitation_note", "Bring your student card")
	post.Set("participant_email_address1", "student@uni.example")
	post.Set("agreed", "true")

	done := teacher.post("/irma_session/store", post)
	if done.StatusCode != http.StatusOK {
		t.Fatalf("exam store: status %d: %s", done.StatusCode, done.Body)
	}

	invitation := e.smtp.To("student@uni.example")
	if len(invitation) != 1 {
		t.Fatalf("student got %d mails, want 1", len(invitation))
	}

	if !strings.Contains(invitation[0].Body(), "Bring your student card") || !strings.Contains(invitation[0].Body(), "Tanja Teacher") {
		t.Errorf("invitation misses the note or the teacher:\n%s", invitation[0].Body())
	}

	if got := invitation[0].Header(headerReplyTo); !strings.Contains(got, "tanja@uni.example") {
		t.Errorf("Reply-To %q, want the teacher", got)
	}

	if len(e.smtp.To("tanja@uni.example")) != 1 {
		t.Errorf("teacher did not get a confirmation")
	}

	if got := e.store.Participants(e.store.Meetings()[0].ID); len(got) != 1 || got[0] != "student@uni.example" {
		t.Errorf("stored participants %v", got)
	}
}

func TestExamTeacherAndStudentJoin(t *testing.T) {
	e := newEnv(t)
	teacher, student := e.newUser(), e.newUser()

	teacher.disclose("teacher", teacherAttrs("Tanja Teacher", "tanja@uni.example"))

	post := createForm("exam", "Oral exam")
	post.Set("invitation_note", "Good luck")
	post.Set("participant_email_address1", "student@uni.example")

	if resp := teacher.post("/irma_session/store", post); resp.StatusCode != http.StatusOK {
		t.Fatalf("store: %d %s", resp.StatusCode, resp.Body)
	}

	m := e.store.Meetings()[0]

	asHost := teacher.get("/irma_session/join_host/" + m.ID)
	if asHost.StatusCode != http.StatusFound {
		t.Fatalf("teacher join: status %d: %s", asHost.StatusCode, asHost.Body)
	}

	if !strings.Contains(asHost.location(), "password="+m.ModeratorPassword) {
		t.Errorf("teacher did not get the moderator password: %s", asHost.location())
	}

	student.disclose("student", studentAttrs("Sam Student", "student@uni.example"))

	asStudent := student.get("/irma_session/join_participant/" + m.ID)
	if asStudent.StatusCode != http.StatusFound {
		t.Fatalf("student join: status %d: %s", asStudent.StatusCode, asStudent.Body)
	}

	if !strings.Contains(asStudent.location(), "password="+m.AttendeePassword) || !strings.Contains(asStudent.location(), "fullName=Sam+Student") {
		t.Errorf("student join URL wrong: %s", asStudent.location())
	}
}

func TestStudentWithTeacherEmail(t *testing.T) {
	e := newEnv(t)
	teacher, impostor := e.newUser(), e.newUser()

	teacher.disclose("teacher", teacherAttrs("Tanja Teacher", "tanja@uni.example"))

	if resp := teacher.post("/irma_session/store", createForm("exam", "Oral exam")); resp.StatusCode != http.StatusOK {
		t.Fatalf("store: %d %s", resp.StatusCode, resp.Body)
	}

	m := e.store.Meetings()[0]

	impostor.disclose("student", studentAttrs("Tanja Teacher", "tanja@uni.example"))

	resp := impostor.get("/irma_session/join_participant/" + m.ID)
	if !strings.Contains(resp.location(), "password="+m.AttendeePassword) {
		t.Errorf("a student disclosure must only give the attendee password: %s", resp.location())
	}
}

func TestHostJoinsDirectly(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	resp := host.get("/irma_session/join/" + m.ID)
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("status %d, want a redirect to BigBlueButton: %s", resp.StatusCode, resp.Body)
	}

	target, err := url.Parse(resp.location())
	if err != nil || !strings.HasPrefix(resp.location(), e.bbb.BaseURL()+"api/join?") {
		t.Fatalf("redirect %q does not go to the BigBlueButton join API", resp.location())
	}

	q := target.Query()
	if q.Get("password") != m.ModeratorPassword || q.Get("fullName") != "Hannah Host" || q.Get("meetingID") != m.RoomID || q.Get("redirect") != "true" {
		t.Errorf("join parameters wrong: %v", q)
	}

	creates := e.bbb.Calls("create")
	if len(creates) != 1 {
		t.Fatalf("%d create calls, want 1", len(creates))
	}

	c := creates[0].Params
	if c.Get("name") != "Weekly sync" || c.Get("moderatorPW") != m.ModeratorPassword || c.Get("attendeePW") != m.AttendeePassword {
		t.Errorf("create parameters wrong: %v", c)
	}
}

func TestParticipantChoosesRole(t *testing.T) {
	e := newEnv(t)
	host, guest := e.newUser(), e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	landing := guest.get("/irma_session/join/" + m.ID)
	if landing.StatusCode != http.StatusOK {
		t.Fatalf("landing: status %d", landing.StatusCode)
	}

	for _, want := range []string{"Neem deel als host", "Neem deel als deelnemer", "/irma_session/join_host/" + m.ID, "/irma_session/join_participant/" + m.ID} {
		if !strings.Contains(landing.Body, want) {
			t.Errorf("landing page misses %q", want)
		}
	}

	if len(e.bbb.Calls("create")) != 0 {
		t.Fatal("the landing page must not create a BigBlueButton room")
	}

	guest.disclose("default", defaultAttrs("Gerrit Guest", "guest@example.com"))

	resp := guest.get("/irma_session/join_participant/" + m.ID)
	if resp.StatusCode != http.StatusFound {
		t.Fatalf("join: status %d: %s", resp.StatusCode, resp.Body)
	}

	if !strings.Contains(resp.location(), "password="+m.AttendeePassword) || !strings.Contains(resp.location(), "fullName=Gerrit+Guest") {
		t.Errorf("participant join URL wrong: %s", resp.location())
	}
}

func TestHostChoiceOtherEmail(t *testing.T) {
	e := newEnv(t)
	host, guest := e.newUser(), e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	guest.disclose("default", defaultAttrs("Gerrit Guest", "guest@example.com"))

	resp := guest.get("/irma_session/join_host/" + m.ID)
	if !strings.Contains(resp.location(), "password="+m.AttendeePassword) {
		t.Errorf("only the host may get the moderator password: %s", resp.location())
	}
}

func TestJoinNameJSONFormat(t *testing.T) {
	e := newEnv(t, withNameFormat("json"))
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	resp := host.get("/irma_session/join_host/" + m.ID)

	target, _ := url.Parse(resp.location())
	want := `{ "name"&#58; "Hannah Host","pbdf.gemeente.personalData.fullname"&#58; "Hannah Host","pbdf.pbdf.email.email"&#58; "Host@Example.com"}`

	if got := target.Query().Get("fullName"); got != want {
		t.Errorf("fullName\n got %s\nwant %s", got, want)
	}
}

func TestLinkExpiresAfter24h(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")

	e.clock.Advance(meeting.LinkValidity - time.Second)

	if resp := host.get("/irma_session/join/" + m.ID); resp.StatusCode != http.StatusFound {
		t.Fatalf("just before expiry: status %d, want the link to work", resp.StatusCode)
	}

	e.clock.Advance(time.Second)

	for _, path := range []string{"/irma_session/join/", "/irma_session/join_host/", "/irma_session/join_participant/"} {
		if resp := host.get(path + m.ID); resp.StatusCode != http.StatusGone {
			t.Errorf("%s after 24 hours: status %d, want 410", path, resp.StatusCode)
		}
	}
}

func TestExpiredMeetingsAreDeleted(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("teacher", teacherAttrs("Tanja Teacher", "tanja@uni.example"))

	post := createForm("exam", "Old exam")
	post.Set("invitation_note", "note")
	post.Set("participant_email_address1", "student@uni.example")
	host.post("/irma_session/store", post)

	old := e.store.Meetings()[0]

	e.clock.Advance(meeting.LinkValidity / 2)

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))
	host.post("/irma_session/store", createForm("free", "Fresh meeting"))

	e.clock.Advance(meeting.LinkValidity/2 + time.Minute)

	removed, err := e.meetings.Cleanup(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if removed.Meetings != 1 || removed.Participants != 1 {
		t.Errorf("removed %+v, want 1 meeting and 1 participant", removed)
	}

	if got := e.store.Meetings(); len(got) != 1 || got[0].Name != "Fresh meeting" {
		t.Errorf("remaining meetings %+v, want only the fresh one", got)
	}

	if got := e.store.Participants(old.ID); len(got) != 0 {
		t.Errorf("participants of the expired meeting survived: %v", got)
	}

	if resp := host.get("/irma_session/join/" + old.ID); resp.StatusCode != http.StatusNotFound {
		t.Errorf("deleted meeting: status %d, want 404", resp.StatusCode)
	}
}

func TestAttrsComeFromSession(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	for range 3 {
		if resp := host.get("/irma_session/create/free"); resp.StatusCode != http.StatusOK {
			t.Fatalf("status %d", resp.StatusCode)
		}
	}

	host.post("/irma_session/store", createForm("free", "Weekly sync"))

	if got := e.irma.ResultFetches(); got != 1 {
		t.Errorf("IRMA result fetched %d times, want once per disclosure", got)
	}

	if got := e.irma.StartedCount(); got != 1 {
		t.Errorf("%d Yivi sessions started, want 1", got)
	}
}

func TestHostIdentityNotSpoofable(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	form := createForm("free", "Weekly sync")
	form.Set("hoster_name", "Mark Rutte")
	form.Set("hoster_email_address", "pm@example.com")

	if resp := host.post("/irma_session/store", form); resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d: %s", resp.StatusCode, resp.Body)
	}

	m := e.store.Meetings()[0]
	if m.HostName != "Hannah Host" || m.HostEmail != "host@example.com" {
		t.Errorf("stored host %q <%s>, want the verified identity", m.HostName, m.HostEmail)
	}

	if len(e.smtp.To("pm@example.com")) != 0 {
		t.Error("mail was sent to the address from the form")
	}
}

func TestMeetingPasswordsRandom(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	a := e.newFreeMeeting(host, "First")
	b := e.newFreeMeeting(host, "Second")

	for name, v := range map[string]string{"a moderator": a.ModeratorPassword, "a attendee": a.AttendeePassword, "b moderator": b.ModeratorPassword, "b attendee": b.AttendeePassword} {
		if len(v) != 32 {
			t.Errorf("%s password %q is not 32 hex characters", name, v)
		}
	}

	distinct := map[string]bool{a.ModeratorPassword: true, a.AttendeePassword: true, b.ModeratorPassword: true, b.AttendeePassword: true}
	if len(distinct) != 4 {
		t.Error("passwords are shared between roles or meetings")
	}

	if a.RoomID == a.ID || strings.Contains(a.ModeratorPassword, a.ID) {
		t.Error("the BigBlueButton room id and passwords must not derive from the link id")
	}
}

func TestMailFailureKeepsMeeting(t *testing.T) {
	e := newEnv(t, withSMTPPort(unreachablePort))
	host := e.newUser()

	host.disclose("default", defaultAttrs("Hannah Host", "host@example.com"))

	resp := host.post("/irma_session/store", createForm("free", "Weekly sync"))
	if resp.StatusCode != http.StatusOK || !linkPattern.MatchString(resp.Body) {
		t.Fatalf("status %d: the host must still get the link: %s", resp.StatusCode, resp.Body)
	}

	if !strings.Contains(resp.Body, "kon niet per e-mail worden verstuurd") {
		t.Error("the success page does not warn that the mail failed")
	}
}

func TestBBBFailureShowsError(t *testing.T) {
	e := newEnv(t)
	host := e.newUser()

	m := e.newFreeMeeting(host, "Weekly sync")
	e.bbb.FailCreate = true

	resp := host.get("/irma_session/join/" + m.ID)
	if resp.StatusCode != http.StatusBadGateway || !strings.Contains(resp.Body, "We kunnen de vergadering niet aanmaken") {
		t.Errorf("status %d, body %s", resp.StatusCode, resp.Body)
	}
}
