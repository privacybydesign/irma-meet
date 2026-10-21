package testutil

import (
	"bufio"
	"bytes"
	"io"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// Mail is one message the fake SMTP server accepted.
type Mail struct {
	From string
	To   []string
	Raw  string
}

// Header returns a decoded header of the message.
func (m Mail) Header(name string) string {
	msg, err := mail.ReadMessage(strings.NewReader(m.Raw))
	if err != nil {
		return ""
	}

	value, _ := new(mime.WordDecoder).DecodeHeader(msg.Header.Get(name))

	return value
}

// Body returns the decoded message body.
func (m Mail) Body() string {
	msg, err := mail.ReadMessage(strings.NewReader(m.Raw))
	if err != nil {
		return ""
	}

	reader := io.Reader(msg.Body)
	if strings.EqualFold(msg.Header.Get("Content-Transfer-Encoding"), "quoted-printable") {
		reader = quotedprintable.NewReader(msg.Body)
	}

	body, _ := io.ReadAll(reader)

	return string(body)
}

// FakeSMTP is an SMTP server that stores what it receives.
type FakeSMTP struct {
	Host string
	Port int

	listener net.Listener

	mu    sync.Mutex
	mails []Mail
}

func NewFakeSMTP(t testing.TB) *FakeSMTP {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	host, portText, _ := net.SplitHostPort(listener.Addr().String())
	port, _ := strconv.Atoi(portText)

	f := &FakeSMTP{Host: host, Port: port, listener: listener}

	go f.accept()
	t.Cleanup(func() { _ = listener.Close() })

	return f
}

func (f *FakeSMTP) Mails() []Mail {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]Mail(nil), f.mails...)
}

// To returns the mails addressed to one recipient.
func (f *FakeSMTP) To(address string) []Mail {
	var out []Mail

	for _, m := range f.Mails() {
		for _, to := range m.To {
			if strings.EqualFold(to, address) {
				out = append(out, m)
			}
		}
	}

	return out
}

func (f *FakeSMTP) accept() {
	for {
		conn, err := f.listener.Accept()
		if err != nil {
			return
		}

		go f.session(conn)
	}
}

func (f *FakeSMTP) session(conn net.Conn) {
	defer func() { _ = conn.Close() }()

	rd := bufio.NewReader(conn)
	reply := func(line string) { _, _ = conn.Write([]byte(line + "\r\n")) }

	reply("220 fake.smtp ESMTP")

	var current Mail

	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			return
		}

		cmd := strings.ToUpper(strings.TrimSpace(line))

		switch {
		case strings.HasPrefix(cmd, "EHLO"), strings.HasPrefix(cmd, "HELO"):
			reply("250 fake.smtp")
		case strings.HasPrefix(cmd, "MAIL FROM:"):
			current = Mail{From: address(line)}
			reply("250 ok")
		case strings.HasPrefix(cmd, "RCPT TO:"):
			current.To = append(current.To, address(line))
			reply("250 ok")
		case cmd == "DATA":
			reply("354 go ahead")

			current.Raw = readData(rd)

			f.mu.Lock()
			f.mails = append(f.mails, current)
			f.mu.Unlock()

			reply("250 queued")
		case cmd == "QUIT":
			reply("221 bye")
			return
		default:
			reply("250 ok")
		}
	}
}

func address(line string) string {
	start, end := strings.Index(line, "<"), strings.LastIndex(line, ">")
	if start < 0 || end < start {
		return ""
	}

	return line[start+1 : end]
}

func readData(rd *bufio.Reader) string {
	var data bytes.Buffer

	for {
		line, err := rd.ReadString('\n')
		if err != nil || line == ".\r\n" {
			return data.String()
		}

		data.WriteString(strings.TrimPrefix(line, "."))
	}
}
