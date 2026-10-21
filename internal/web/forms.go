package web

import (
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	maxFieldRunes   = 255
	maxParticipants = 6

	fieldMeetingName = "meeting_name"
	fieldNote        = "invitation_note"
	fieldAgreed      = "agreed"
	fieldParticipant = "participant_email_address"
)

// createForm is the submitted "create meeting" form.
type createForm struct {
	name         string
	note         string
	participants [maxParticipants]string
	agreed       string
}

func parseCreateForm(w http.ResponseWriter, r *http.Request) (createForm, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)

	if err := r.ParseForm(); err != nil {
		return createForm{}, err
	}

	f := createForm{
		name:   cleanLine(r.PostFormValue(fieldMeetingName)),
		note:   strings.TrimSpace(r.PostFormValue(fieldNote)),
		agreed: strings.ToLower(strings.TrimSpace(r.PostFormValue(fieldAgreed))),
	}

	for i := range f.participants {
		f.participants[i] = strings.TrimSpace(r.PostFormValue(fieldParticipant + string(rune('1'+i))))
	}

	return f, nil
}

// cleanLine trims s and turns control characters into spaces; names end up
// in mail subjects and BigBlueButton calls.
func cleanLine(s string) string {
	return strings.TrimSpace(strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return ' '
		}

		return r
	}, s))
}

// fieldError is one validation failure: the field it concerns and the
// translation key and label of the message.
type fieldError struct {
	field    string
	message  string
	label    string
	maxChars string
}

// validate checks the form the way the PHP application did: the meeting
// name is required, a note and a first participant go together, every
// participant address must be an email address, and the box must be checked.
func (f createForm) validate() []fieldError {
	var errs []fieldError

	if f.name == "" {
		errs = append(errs, required(fieldMeetingName, "Meeting name"))
	} else if utf8.RuneCountInString(f.name) > maxFieldRunes {
		errs = append(errs, tooLong(fieldMeetingName, "Meeting name"))
	}

	if f.note == "" && f.participants[0] != "" {
		errs = append(errs, required(fieldNote, "Invitation note"))
	}

	if f.note != "" && utf8.RuneCountInString(f.note) > maxFieldRunes {
		errs = append(errs, tooLong(fieldNote, "Invitation note"))
	}

	if f.note != "" && f.participants[0] == "" {
		errs = append(errs, required(fieldParticipant+"1", "Participant email address"))
	}

	for i, email := range f.participants {
		if email != "" && !validEmail(email) {
			errs = append(errs, fieldError{
				field:   fieldParticipant + string(rune('1'+i)),
				message: "The :attribute field must be a valid email address.",
				label:   "Participant email address",
			})
		}
	}

	if !accepted(f.agreed) {
		errs = append(errs, fieldError{field: fieldAgreed, message: "Please check the box to allow processing."})
	}

	return errs
}

func required(field, label string) fieldError {
	return fieldError{field: field, message: "The :attribute field is required.", label: label}
}

func tooLong(field, label string) fieldError {
	return fieldError{field: field, message: "The :attribute field must not be greater than :max characters.", label: label, maxChars: strconv.Itoa(maxFieldRunes)}
}

// emails returns the participant addresses that were filled in.
func (f createForm) emails() []string {
	var out []string

	for _, email := range f.participants {
		if email != "" {
			out = append(out, email)
		}
	}

	return out
}

func accepted(v string) bool {
	switch v {
	case "yes", "on", "1", "true":
		return true
	}

	return false
}

// validEmail accepts a bare address, not "Name <address>" and not a list.
func validEmail(s string) bool {
	addr, err := mail.ParseAddress(s)

	return err == nil && addr.Address == s && strings.Contains(addr.Address[strings.LastIndex(addr.Address, "@"):], ".")
}
