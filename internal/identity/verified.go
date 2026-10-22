package identity

import (
	"bytes"
	"encoding/json"
	"strings"
)

// NameFormat is how the participant name is handed to BigBlueButton.
type NameFormat string

const (
	// NameFormatString passes the plain display name.
	NameFormatString NameFormat = "string"
	// NameFormatJSON passes a JSON object with the name and the verified attributes.
	NameFormatJSON NameFormat = "json"
)

// colonEntity separates key and value in the BigBlueButton JSON name.
const colonEntity = "&#58; "

// Verified holds the attributes a person disclosed in a valid Yivi session.
type Verified struct {
	Type  Type              `json:"type"`
	Attrs map[string]string `json:"attrs"`
}

// Email returns the lowercased disclosed email address, or "".
func (v Verified) Email() string {
	email := ""
	for _, id := range specs[v.Type].emails {
		if value := v.Attrs[id]; value != "" {
			email = value
		}
	}

	return strings.ToLower(email)
}

// Name returns the first alternative of the name attributes that is fully disclosed.
func (v Verified) Name() string {
	for _, alternative := range specs[v.Type].names {
		if parts, ok := v.collect(alternative); ok {
			return strings.Join(parts, " ")
		}
	}

	return ""
}

// BBBName returns the name BigBlueButton shows for this person, or "" when
// the disclosure does not hold a complete identity.
func (v Verified) BBBName(format NameFormat) string {
	name := v.Name()
	if format != NameFormatJSON || name == "" {
		return name
	}

	for _, ids := range specs[v.Type].valid {
		if _, ok := v.collect(ids); !ok {
			continue
		}

		return v.jsonName(name, ids)
	}

	return ""
}

func (v Verified) collect(ids []string) ([]string, bool) {
	values := make([]string, len(ids))
	for i, id := range ids {
		values[i] = v.Attrs[id]
		if values[i] == "" {
			return nil, false
		}
	}

	return values, true
}

// jsonName renders {"name": ..., "<attribute>": ...}. The colons are written
// as the HTML entity &#58; because that is the format the BigBlueButton
// client of this service has always been fed.
func (v Verified) jsonName(name string, ids []string) string {
	var out strings.Builder

	out.WriteString("{ ")
	writeMember(&out, "name", name)

	for _, id := range ids {
		out.WriteString(",")
		writeMember(&out, id, v.Attrs[id])
	}

	out.WriteString("}")

	return out.String()
}

func writeMember(out *strings.Builder, key, value string) {
	out.WriteString(quote(key))
	out.WriteString(colonEntity)
	out.WriteString(quote(value))
}

func quote(s string) string {
	var buf bytes.Buffer

	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(s)

	return strings.TrimSuffix(buf.String(), "\n")
}
