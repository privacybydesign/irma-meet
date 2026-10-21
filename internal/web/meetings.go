package web

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/meeting"
)

const (
	createPath       = "/irma_session/create/"
	joinPagePath     = "/irma_session/join/"
	joinHostPath     = "/irma_session/join_host/"
	joinParticipPath = "/irma_session/join_participant/"
)

// meetingIDPattern is the shape of the ids Create hands out.
var meetingIDPattern = regexp.MustCompile(`^[0-9a-f]{1,48}$`)

// createView is the data of the create form, filled or not.
type createView struct {
	Name        string
	Note        string
	Participant string
	HostName    string
	HostEmail   string
	Errors      []string
	Invalid     map[string]bool
}

type createdView struct {
	Name       string
	Link       string
	MailFailed bool
}

type joinView struct {
	HostLabel           string
	HostStartURL        string
	HostURL             string
	ParticipantLabel    string
	ParticipantStartURL string
	ParticipantURL      string
}

func (a *app) createForm(w http.ResponseWriter, r *http.Request) {
	mt, ok := meeting.ParseType(r.PathValue("meetingType"))
	if !ok {
		a.notFound(w, r)
		return
	}

	host, ok := a.identityFor(w, r, mt.HostIdentity(), r.URL.Path)
	if !ok {
		return
	}

	a.renderCreate(w, r, mt, http.StatusOK, createView{HostName: host.Name(), HostEmail: host.Email()})
}

func (a *app) renderCreate(w http.ResponseWriter, r *http.Request, mt meeting.Type, status int, v createView) {
	page := pageCreateFree
	if mt == meeting.Exam {
		page = pageCreateExam
	}

	a.page(w, r, page, status, v)
}

func (a *app) store(w http.ResponseWriter, r *http.Request) {
	form, err := parseCreateForm(w, r)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mt, ok := meeting.ParseType(r.PostFormValue(fieldMeetingType))
	if !ok {
		a.notFound(w, r)
		return
	}

	host, ok := a.identityFor(w, r, mt.HostIdentity(), createPath+string(mt))
	if !ok {
		return
	}

	if errs := form.validate(); len(errs) > 0 {
		a.renderCreate(w, r, mt, http.StatusUnprocessableEntity, a.invalidView(r, form, host, errs))
		return
	}

	created, err := a.meetings.Create(r.Context(), meeting.NewMeeting{
		Type:         mt,
		Name:         form.name,
		Note:         form.note,
		Host:         host,
		Participants: form.emails(),
		Locale:       a.lang(r),
	})
	if err != nil {
		a.serverError(w, r, "creating meeting failed", err)
		return
	}

	if created.MailErr != nil {
		a.log.Error("sending meeting mails failed", "error", created.MailErr)
	}

	a.page(w, r, pageCreated, http.StatusOK, createdView{Name: created.Meeting.Name, Link: created.Link, MailFailed: created.MailErr != nil})
}

// invalidView refills the form with what was submitted and the messages.
func (a *app) invalidView(r *http.Request, f createForm, host identity.Verified, errs []fieldError) createView {
	v := createView{
		Name: f.name, Note: f.note, Participant: f.participants[0],
		HostName: host.Name(), HostEmail: host.Email(),
		Invalid: map[string]bool{},
	}

	lang := a.lang(r)

	for _, e := range errs {
		v.Invalid[e.field] = true

		replace := map[string]string{"max": e.maxChars}
		if e.label != "" {
			replace["attribute"] = strings.ToLower(a.text.T(lang, e.label))
		}

		v.Errors = append(v.Errors, a.text.Tr(lang, e.message, replace))
	}

	return v
}

// join is the landing page of a meeting link. Anyone may open it; the room is
// only created once someone has disclosed through join_host or join_participant.
func (a *app) join(w http.ResponseWriter, r *http.Request) {
	m, ok := a.loadMeeting(w, r)
	if !ok {
		return
	}

	if v, ok, err := a.identityOf(r); err == nil && ok && m.IsHost(v) {
		a.enter(w, r, m, v)
		return
	}

	host, participant := m.Type.HostIdentity(), m.Type.ParticipantIdentity()

	a.page(w, r, pageJoin, http.StatusOK, joinView{
		HostLabel:           host.HostLabel(),
		HostStartURL:        startPath(host),
		HostURL:             joinHostPath + m.ID,
		ParticipantLabel:    participant.ParticipantLabel(),
		ParticipantStartURL: startPath(participant),
		ParticipantURL:      joinParticipPath + m.ID,
	})
}

func (a *app) joinHost(w http.ResponseWriter, r *http.Request) {
	a.joinAs(w, r, joinHostPath, meeting.Type.HostIdentity)
}

func (a *app) joinParticipant(w http.ResponseWriter, r *http.Request) {
	a.joinAs(w, r, joinParticipPath, meeting.Type.ParticipantIdentity)
}

// joinAs lets in whoever disclosed the identity role asks for. Which password
// they get depends on who they turn out to be, not on the route they used.
func (a *app) joinAs(w http.ResponseWriter, r *http.Request, route string, role func(meeting.Type) identity.Type) {
	m, ok := a.loadMeeting(w, r)
	if !ok {
		return
	}

	who, ok := a.identityFor(w, r, role(m.Type), route+m.ID)
	if !ok {
		return
	}

	a.enter(w, r, m, who)
}

func (a *app) enter(w http.ResponseWriter, r *http.Request, m meeting.Meeting, who identity.Verified) {
	target, err := a.meetings.JoinURL(r.Context(), m, who)

	switch {
	case errors.Is(err, meeting.ErrNotVerified):
		a.fail(w, r, http.StatusForbidden, "Your attributes could not be verified by IRMA.")
	case errors.Is(err, meeting.ErrRoom):
		a.log.Error("creating meeting room failed", "error", err)
		a.fail(w, r, http.StatusBadGateway, "Can't create room! please contact our administrator.")
	case err != nil:
		a.serverError(w, r, "joining meeting failed", err)
	default:
		http.Redirect(w, r, target, http.StatusFound)
	}
}

func (a *app) loadMeeting(w http.ResponseWriter, r *http.Request) (meeting.Meeting, bool) {
	id := r.PathValue("irmaSessionId")
	if !meetingIDPattern.MatchString(id) {
		a.notFound(w, r)
		return meeting.Meeting{}, false
	}

	m, err := a.meetings.Lookup(r.Context(), id)

	switch {
	case errors.Is(err, meeting.ErrNotFound):
		a.notFound(w, r)
	case errors.Is(err, meeting.ErrExpired):
		a.fail(w, r, http.StatusGone, "This meeting link has expired.")
	case err != nil:
		a.serverError(w, r, "looking up meeting failed", err)
	default:
		return m, true
	}

	return meeting.Meeting{}, false
}
