package identity_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

const (
	attrEmail      = "pbdf.pbdf.email.email"
	attrEmailSidn  = "pbdf.sidn-pbdf.email.email"
	attrFullName   = "pbdf.gemeente.personalData.fullname"
	attrFirstName  = "pbdf.pbdf.linkedin.firstname"
	attrFamilyName = "pbdf.pbdf.linkedin.familyname"
	attrSurfFamily = "pbdf.pbdf.surfnet-2.familyname"
	attrSurfEmail  = "pbdf.pbdf.surfnet-2.email"
	attrSurfInst   = "pbdf.pbdf.surfnet-2.institute"
	attrSurfType   = "pbdf.pbdf.surfnet-2.type"
)

func verified(t identity.Type, attrs map[string]string) identity.Verified {
	return identity.Verified{Type: t, Attrs: attrs}
}

// The email is lowercased when present and "" (never an error) when absent.
func TestEmailIsLowercased(t *testing.T) {
	v := verified(identity.TypeDefault, map[string]string{attrEmail: "Alice@Example.COM"})

	if got := v.Email(); got != "alice@example.com" {
		t.Errorf("Email() = %q", got)
	}
}

func TestEmailEmptyWhenAbsent(t *testing.T) {
	for _, v := range []identity.Verified{
		verified(identity.TypeDefault, nil),
		verified(identity.TypeDefault, map[string]string{attrFullName: "Alice"}),
		verified(identity.TypeStudent, map[string]string{attrEmail: "alice@example.com"}),
		{},
	} {
		if got := v.Email(); got != "" {
			t.Errorf("%+v: Email() = %q, want empty", v, got)
		}
	}
}

func TestEmailEitherCredential(t *testing.T) {
	sidn := verified(identity.TypeDefault, map[string]string{attrEmailSidn: "sidn@example.com"})
	if got := sidn.Email(); got != "sidn@example.com" {
		t.Errorf("sidn email = %q", got)
	}

	both := verified(identity.TypeDefault, map[string]string{attrEmailSidn: "sidn@example.com", attrEmail: "pbdf@example.com"})
	if got := both.Email(); got != "pbdf@example.com" {
		t.Errorf("with both credentials the pbdf email wins, got %q", got)
	}
}

// The teacher "name" and the student "email" were values the PHP code
// iterated the wrong way, so exams answered 500.
func TestExamTypesNameAndEmail(t *testing.T) {
	for _, typ := range []identity.Type{identity.TypeTeacher, identity.TypeStudent} {
		v := verified(typ, map[string]string{attrSurfFamily: "Jansen", attrSurfEmail: "J.Jansen@uni.example", attrSurfInst: "uni.example", attrSurfType: "x"})

		if got := v.Name(); got != "Jansen" {
			t.Errorf("%s: Name() = %q", typ, got)
		}

		if got := v.Email(); got != "j.jansen@uni.example" {
			t.Errorf("%s: Email() = %q", typ, got)
		}
	}
}

func TestNameAlternatives(t *testing.T) {
	cases := map[string]struct {
		attrs map[string]string
		want  string
	}{
		"civil registry":    {map[string]string{attrFullName: "Alice Jansen"}, "Alice Jansen"},
		"linkedin":          {map[string]string{attrFirstName: "Alice", attrFamilyName: "Jansen"}, "Alice Jansen"},
		"registry wins":     {map[string]string{attrFullName: "Alice Jansen", attrFirstName: "Al", attrFamilyName: "J"}, "Alice Jansen"},
		"half of linkedin":  {map[string]string{attrFirstName: "Alice"}, ""},
		"no name":           {map[string]string{attrEmail: "alice@example.com"}, ""},
		"empty value":       {map[string]string{attrFullName: ""}, ""},
		"linkedin partial2": {map[string]string{attrFamilyName: "Jansen"}, ""},
	}

	for name, c := range cases {
		if got := verified(identity.TypeDefault, c.attrs).Name(); got != c.want {
			t.Errorf("%s: Name() = %q, want %q", name, got, c.want)
		}
	}
}

func TestBBBName(t *testing.T) {
	full := verified(identity.TypeDefault, map[string]string{attrFullName: "Alice Jansen", attrEmail: "alice@example.com"})

	if got := full.BBBName(identity.NameFormatString); got != "Alice Jansen" {
		t.Errorf("string format: %q", got)
	}

	wantJSON := `{ "name"&#58; "Alice Jansen","pbdf.gemeente.personalData.fullname"&#58; "Alice Jansen","pbdf.pbdf.email.email"&#58; "alice@example.com"}`
	if got := full.BBBName(identity.NameFormatJSON); got != wantJSON {
		t.Errorf("json format:\n got %s\nwant %s", got, wantJSON)
	}

	linkedin := verified(identity.TypeDefault, map[string]string{attrFirstName: "Alice", attrFamilyName: "Jansen", attrEmailSidn: "a@example.com"})
	wantLinkedin := `{ "name"&#58; "Alice Jansen","pbdf.pbdf.linkedin.firstname"&#58; "Alice","pbdf.pbdf.linkedin.familyname"&#58; "Jansen","pbdf.sidn-pbdf.email.email"&#58; "a@example.com"}`

	if got := linkedin.BBBName(identity.NameFormatJSON); got != wantLinkedin {
		t.Errorf("linkedin json format:\n got %s\nwant %s", got, wantLinkedin)
	}
}

func TestBBBNameJSONEscapes(t *testing.T) {
	v := verified(identity.TypeDefault, map[string]string{attrFullName: `Al "Quote" <b>\`, attrEmail: "a@example.com"})

	got := v.BBBName(identity.NameFormatJSON)
	if !strings.Contains(got, `"Al \"Quote\" <b>\\"`) {
		t.Errorf("quotes and backslashes must be escaped: %s", got)
	}

	if strings.Contains(got, `"name": `) {
		t.Errorf("the key/value separator must be the &#58; entity: %s", got)
	}
}

func TestBBBNameNeedsFullIdentity(t *testing.T) {
	onlyName := verified(identity.TypeDefault, map[string]string{attrFullName: "Alice Jansen"})

	if got := onlyName.BBBName(identity.NameFormatJSON); got != "" {
		t.Errorf("a name without email must not make a JSON identity: %q", got)
	}

	if got := (identity.Verified{}).BBBName(identity.NameFormatString); got != "" {
		t.Errorf("empty identity: %q", got)
	}
}

func TestParseKeptTypesOnly(t *testing.T) {
	for _, name := range []string{"default", "teacher", "student"} {
		if _, ok := identity.Parse(name); !ok {
			t.Errorf("%s should parse", name)
		}
	}

	for _, name := range []string{"medical", "patient", "", "DEFAULT", "default "} {
		if _, ok := identity.Parse(name); ok {
			t.Errorf("%q must not parse", name)
		}
	}
}

func TestStartPostsOneRequest(t *testing.T) {
	irma := testutil.NewFakeIRMA(t, "secret-token")
	svc := identity.NewService(irma.URL(), "secret-token", nil)

	session, err := svc.Start(context.Background(), identity.TypeDefault)
	if err != nil {
		t.Fatal(err)
	}

	if irma.StartedCount() != 1 || session.Token == "" || !strings.Contains(string(session.Ptr), `"u"`) {
		t.Errorf("started %d, session %+v", irma.StartedCount(), session)
	}

	want := `{"@context":"https://irma.app/ld/request/disclosure/v2","disclose":[[["pbdf.pbdf.email.email"],["pbdf.sidn-pbdf.email.email"]],[["pbdf.gemeente.personalData.fullname"],["pbdf.pbdf.linkedin.firstname","pbdf.pbdf.linkedin.familyname"]]]}`
	if got := string(irma.Requests()[0]); got != want {
		t.Errorf("request\n got %s\nwant %s", got, want)
	}

	if _, err := identity.NewService(irma.URL(), "wrong", nil).Start(context.Background(), identity.TypeDefault); err == nil {
		t.Error("a wrong requestor token must fail")
	}
}

func TestVerifyValidDisclosure(t *testing.T) {
	irma := testutil.NewFakeIRMA(t, "tok")
	svc := identity.NewService(irma.URL(), "tok", nil)

	session, _ := svc.Start(context.Background(), identity.TypeDefault)
	irma.CompleteLatest(map[string]string{attrFullName: "Alice Jansen", attrEmail: "Alice@Example.com"})

	v, err := svc.Verify(context.Background(), session.Token, identity.TypeDefault)
	if err != nil {
		t.Fatal(err)
	}

	if v.Type != identity.TypeDefault || v.Name() != "Alice Jansen" || v.Email() != "alice@example.com" {
		t.Errorf("verified %+v", v)
	}
}

func TestVerifyRejectsInvalid(t *testing.T) {
	irma := testutil.NewFakeIRMA(t, "tok")
	svc := identity.NewService(irma.URL(), "tok", nil)

	running, _ := svc.Start(context.Background(), identity.TypeDefault)
	if _, err := svc.Verify(context.Background(), running.Token, identity.TypeDefault); !errors.Is(err, identity.ErrPending) {
		t.Errorf("running session: %v, want ErrPending", err)
	}

	invalid, _ := svc.Start(context.Background(), identity.TypeDefault)
	irma.FailLatest("INVALID_TIMESTAMP")

	if _, err := svc.Verify(context.Background(), invalid.Token, identity.TypeDefault); !errors.Is(err, identity.ErrNotVerified) {
		t.Errorf("invalid proof: %v, want ErrNotVerified", err)
	}

	for _, token := range []string{"", "../session", "a b", strings.Repeat("a", 65), "x/result"} {
		if _, err := svc.Verify(context.Background(), token, identity.TypeDefault); !errors.Is(err, identity.ErrNotVerified) {
			t.Errorf("token %q: %v, want ErrNotVerified without a request", token, err)
		}
	}
}

func TestVerifyDropsAbsentAttrs(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"status":"DONE","proofStatus":"VALID","disclosed":[[
			{"id":"` + attrEmail + `","rawvalue":"a@example.com","status":"PRESENT"},
			{"id":"` + attrFullName + `","rawvalue":"Injected","status":"NULL"}]]}`))
	}))
	defer server.Close()

	v, err := identity.NewService(server.URL, "tok", nil).Verify(context.Background(), "abc", identity.TypeDefault)
	if err != nil {
		t.Fatal(err)
	}

	if v.Name() != "" || v.Email() != "a@example.com" {
		t.Errorf("attributes %v", v.Attrs)
	}
}

func TestVerifyServerFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := identity.NewService(server.URL, "tok", nil).Verify(context.Background(), "abc", identity.TypeDefault)
	if err == nil || errors.Is(err, identity.ErrNotVerified) || errors.Is(err, identity.ErrPending) {
		t.Errorf("a failing IRMA server is neither 'not verified' nor 'pending': %v", err)
	}
}
