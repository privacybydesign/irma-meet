// Package identity defines what attributes a Yivi disclosure asks for and
// talks to the IRMA server to run and verify those disclosures.
package identity

import "encoding/json"

// Type names a disclosure: the set of attributes a person has to reveal.
type Type string

const (
	// TypeDefault is an email address plus a name from the civil registry or LinkedIn.
	TypeDefault Type = "default"
	// TypeTeacher is a SURFconext employee card.
	TypeTeacher Type = "teacher"
	// TypeStudent is a SURFconext student card.
	TypeStudent Type = "student"
)

const disclosureContext = "https://irma.app/ld/request/disclosure/v2"

const (
	attrEmail        = "pbdf.pbdf.email.email"
	attrEmailSidn    = "pbdf.sidn-pbdf.email.email"
	attrFullName     = "pbdf.gemeente.personalData.fullname"
	attrFirstName    = "pbdf.pbdf.linkedin.firstname"
	attrFamilyName   = "pbdf.pbdf.linkedin.familyname"
	attrSurfInstit   = "pbdf.pbdf.surfnet-2.institute"
	attrSurfType     = "pbdf.pbdf.surfnet-2.type"
	attrSurfFamily   = "pbdf.pbdf.surfnet-2.familyname"
	attrSurfEmail    = "pbdf.pbdf.surfnet-2.email"
	surfTypeEmployee = "employee"
	surfTypeStudent  = "student"

	// The labels are translation keys.
	labelJoinTeacher = "Join as teacher"
	labelJoinStudent = "Join as student"
)

// attributeRequest is one attribute in a disclosure request. It marshals to a
// bare attribute id, or to {type, value} when the value is constrained.
type attributeRequest struct {
	id    string
	value string
}

func (a attributeRequest) MarshalJSON() ([]byte, error) {
	if a.value == "" {
		return json.Marshal(a.id)
	}

	return json.Marshal(struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}{a.id, a.value})
}

type disclosureRequest struct {
	Context  string                 `json:"@context"`
	Disclose [][][]attributeRequest `json:"disclose"`
}

type spec struct {
	disclose [][][]attributeRequest
	// valid lists the attribute sets that each prove an identity; the first
	// complete one is what the BBB JSON name format reports.
	valid [][]string
	// names lists alternative attribute sets that make up the display name.
	names [][]string
	// emails lists the attributes that can carry the email address.
	emails []string
	// hostLabel and participantLabel are the translation keys of the buttons
	// that start this disclosure as the host or as a participant.
	hostLabel        string
	participantLabel string
}

func plain(ids ...string) []attributeRequest {
	reqs := make([]attributeRequest, len(ids))
	for i, id := range ids {
		reqs[i] = attributeRequest{id: id}
	}

	return reqs
}

var specs = map[Type]spec{
	TypeDefault: {
		disclose: [][][]attributeRequest{
			{plain(attrEmail), plain(attrEmailSidn)},
			{plain(attrFullName), plain(attrFirstName, attrFamilyName)},
		},
		valid: [][]string{
			{attrFullName, attrEmail},
			{attrFirstName, attrFamilyName, attrEmail},
			{attrFullName, attrEmailSidn},
			{attrFirstName, attrFamilyName, attrEmailSidn},
		},
		names:            [][]string{{attrFullName}, {attrFirstName, attrFamilyName}},
		emails:           []string{attrEmailSidn, attrEmail},
		hostLabel:        "Join as host",
		participantLabel: "Join as participant",
	},
	TypeTeacher: {
		disclose: [][][]attributeRequest{
			{surfRequest(surfTypeEmployee)},
		},
		valid:            [][]string{{attrSurfInstit, attrSurfType, attrSurfFamily, attrSurfEmail}},
		names:            [][]string{{attrSurfFamily}},
		emails:           []string{attrSurfEmail},
		hostLabel:        labelJoinTeacher,
		participantLabel: labelJoinTeacher,
	},
	TypeStudent: {
		disclose: [][][]attributeRequest{
			{surfRequest(surfTypeStudent)},
		},
		valid:            [][]string{{attrSurfInstit, attrSurfType, attrSurfFamily, attrSurfEmail}},
		names:            [][]string{{attrSurfFamily}},
		emails:           []string{attrSurfEmail},
		hostLabel:        labelJoinStudent,
		participantLabel: labelJoinStudent,
	},
}

func surfRequest(cardType string) []attributeRequest {
	return []attributeRequest{
		{id: attrSurfInstit},
		{id: attrSurfType, value: cardType},
		{id: attrSurfFamily},
		{id: attrSurfEmail},
	}
}

// Parse returns the Type for a name received from a client.
func Parse(name string) (Type, bool) {
	t := Type(name)
	_, ok := specs[t]

	return t, ok
}

// HostLabel is the translation key of the button that starts this disclosure
// for the host of a meeting.
func (t Type) HostLabel() string {
	return specs[t].hostLabel
}

// ParticipantLabel is the translation key of the button that starts this
// disclosure for a participant.
func (t Type) ParticipantLabel() string {
	return specs[t].participantLabel
}

func (t Type) request() disclosureRequest {
	return disclosureRequest{Context: disclosureContext, Disclose: specs[t].disclose}
}
