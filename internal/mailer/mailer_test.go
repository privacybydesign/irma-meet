package mailer_test

import (
	"context"
	"strings"
	"testing"

	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/mailer"
	"github.com/privacybydesign/irma-meet/internal/testutil"
)

func newMailer(t *testing.T, smtp *testutil.FakeSMTP) *mailer.Mailer {
	t.Helper()

	m, err := mailer.New(config.Mail{Host: smtp.Host, Port: smtp.Port, FromAddress: "noreply@yivi-meet.test", FromName: "Yivi Meet"})
	if err != nil {
		t.Fatal(err)
	}

	return m
}

// sendBody sends msg and returns the body the SMTP server received.
func sendBody(t *testing.T, msg mailer.Message) (string, error) {
	t.Helper()

	smtp := testutil.NewFakeSMTP(t)
	msg.To, msg.ReplyTo = "guest@example.com", "host@example.com"

	if err := newMailer(t, smtp).Send(context.Background(), msg); err != nil {
		return "", err
	}

	return smtp.Mails()[0].Body(), nil
}

func TestEveryTemplateRenders(t *testing.T) {
	data := mailer.Data{HosterName: "Hannah", InvitationLink: "https://meet.example/irma_session/join/abc", InvitationNote: "A note"}

	for _, locale := range []string{"nl", "en"} {
		for _, kind := range []mailer.Kind{mailer.Confirmation, mailer.Invitation} {
			for _, meetingType := range []string{"free", "exam"} {
				body, err := sendBody(t, mailer.Message{Kind: kind, MeetingType: meetingType, Locale: locale, Data: data})
				if err != nil {
					t.Errorf("%s/%s/%s: %v", locale, kind, meetingType, err)
					continue
				}

				if !strings.Contains(body, data.InvitationLink) || strings.Contains(body, "{{") {
					t.Errorf("%s/%s/%s: link missing or template markup left:\n%s", locale, kind, meetingType, body)
				}
			}
		}
	}
}

func TestMedicalConsultIsGone(t *testing.T) {
	if _, err := sendBody(t, mailer.Message{Kind: mailer.Invitation, MeetingType: "medical_consult", Locale: "en"}); err == nil {
		t.Error("the medical_consult mails were dropped")
	}
}

func TestTemplatesEscapeData(t *testing.T) {
	body, err := sendBody(t, mailer.Message{
		Kind: mailer.Invitation, MeetingType: "free", Locale: "en",
		Data: mailer.Data{HosterName: `<script>alert(1)</script>`, InvitationLink: `https://x.example/"><img src=x>`, InvitationNote: `<b>bold</b>`},
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, raw := range []string{"<script>", "<img src=x>", "<b>bold</b>"} {
		if strings.Contains(body, raw) {
			t.Errorf("unescaped %q in\n%s", raw, body)
		}
	}
}

func TestInvitationNoteIsOptional(t *testing.T) {
	without, _ := sendBody(t, mailer.Message{Kind: mailer.Invitation, MeetingType: "free", Locale: "en", Data: mailer.Data{HosterName: "Hannah"}})
	with, _ := sendBody(t, mailer.Message{Kind: mailer.Invitation, MeetingType: "free", Locale: "en", Data: mailer.Data{HosterName: "Hannah", InvitationNote: "Bring snacks"}})

	if strings.Contains(without, "*******") || !strings.Contains(with, "Bring snacks") {
		t.Error("the note block must only appear when there is a note")
	}
}

func TestSendDeliversAll(t *testing.T) {
	smtp := testutil.NewFakeSMTP(t)
	m := newMailer(t, smtp)

	err := m.Send(context.Background(),
		mailer.Message{Kind: mailer.Invitation, MeetingType: "free", Locale: "nl", To: "guest@example.com", ReplyTo: "host@example.com", Subject: "Hello é", Data: mailer.Data{HosterName: "Hannah"}},
		mailer.Message{Kind: mailer.Confirmation, MeetingType: "free", Locale: "nl", To: "host@example.com", ReplyTo: "host@example.com", Subject: "Hello", Data: mailer.Data{HosterName: "Hannah"}},
	)
	if err != nil {
		t.Fatal(err)
	}

	mails := smtp.Mails()
	if len(mails) != 2 {
		t.Fatalf("%d mails delivered, want 2", len(mails))
	}

	first := mails[0]
	if first.From != "noreply@yivi-meet.test" || first.To[0] != "guest@example.com" {
		t.Errorf("envelope %+v", first)
	}

	if got := first.Header("Subject"); got != "Hello é" {
		t.Errorf("subject %q", got)
	}

	if got := first.Header("From"); !strings.Contains(got, "Yivi Meet") {
		t.Errorf("From %q", got)
	}

	if got := first.Header("Content-Type"); !strings.HasPrefix(got, "text/html") {
		t.Errorf("Content-Type %q", got)
	}
}

func TestSubjectIsSingleLine(t *testing.T) {
	smtp := testutil.NewFakeSMTP(t)
	m := newMailer(t, smtp)

	err := m.Send(context.Background(), mailer.Message{
		Kind: mailer.Invitation, MeetingType: "free", Locale: "en", To: "guest@example.com", ReplyTo: "host@example.com",
		Subject: "Hi\r\nBcc: victim@example.com\nX-Evil: 1", Data: mailer.Data{HosterName: "Hannah"},
	})
	if err != nil {
		t.Fatal(err)
	}

	mail := smtp.Mails()[0]
	if mail.Header("Bcc") != "" || mail.Header("X-Evil") != "" {
		t.Errorf("headers were injected through the subject:\n%s", mail.Raw)
	}
}

func TestInvalidRecipient(t *testing.T) {
	smtp := testutil.NewFakeSMTP(t)
	m := newMailer(t, smtp)

	err := m.Send(context.Background(), mailer.Message{Kind: mailer.Invitation, MeetingType: "free", Locale: "en", To: "not an address", ReplyTo: "host@example.com"})
	if err == nil || len(smtp.Mails()) != 0 {
		t.Errorf("err %v, %d mails", err, len(smtp.Mails()))
	}
}

func TestUnreachableServer(t *testing.T) {
	m, err := mailer.New(config.Mail{Host: "127.0.0.1", Port: 1, FromAddress: "noreply@yivi-meet.test"})
	if err != nil {
		t.Fatal(err)
	}

	err = m.Send(context.Background(), mailer.Message{Kind: mailer.Invitation, MeetingType: "free", Locale: "en", To: "a@example.com", ReplyTo: "b@example.com"})
	if err == nil {
		t.Error("expected an error")
	}
}
