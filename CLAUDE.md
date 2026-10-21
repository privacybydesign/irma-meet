# irma-meet repo notes

Go service (`net/http`, `html/template`, MySQL) that lets people create and join
BigBlueButton meetings after a Yivi disclosure. It replaced a Laravel app.

## Build & test

- `go vet ./... && go test ./...`. The Go version is set in `go.mod`.
- Tests need no network, MySQL or Docker: they use the fakes in
  `internal/testutil` (IRMA server, BigBlueButton, SMTP, in-memory store).
- `internal/store` and `TestEndToEndAgainstMySQL` run against a real MySQL
  only when `TEST_MYSQL_DSN` is set (a server DSN, e.g.
  `root:secret@tcp(127.0.0.1:3307)/`; each test makes its own database);
  otherwise they skip. Run them when you change SQL or a migration.

## Layout

- `main.go` wires config, store, services and the HTTP server.
- `internal/web` HTTP handlers, templates (`templates/`), embedded `static/`.
  Handlers call the `Meetings` and `Identities` services, never the store, the
  IRMA server or BigBlueButton directly.
- `internal/meeting` the domain: create, expiry, join, cleanup.
- `internal/identity` disclosure types and the IRMA server client.
- `internal/bbb`, `internal/mailer`, `internal/store`, `internal/i18n`, `internal/config`.

## Gotchas

- Only `internal/web/static` is public. Do not widen the `go:embed` pattern;
  `TestOnlyStaticFilesAreServed` guards it.
- Output is escaped by `html/template`. `.HTML`/`.HTMLf` in a template mark a
  *translation string* as trusted markup; never pass user input to `.HTML`.
  `.HTMLf` escapes its arguments.
- A session holds the disclosed attributes (`verified`) once, so handlers do
  not ask the IRMA server again. `identityOf` in `internal/web/auth.go` is the
  only place that verifies a pending Yivi session.
- CSRF protection is `http.CrossOriginProtection`; it only checks unsafe
  methods, so state-changing endpoints must not be GET.
- Translation keys are the English text. A key used in a template or handler
  needs an entry in `nl.json`; `TestTemplateKeysTranslated` checks it.
- `MAIL_ENCRYPTION` is `starttls` (default), `tls` or `none`. The empty value
  of `config.MailEncryption` means STARTTLS, so tests against the fake SMTP
  server (no TLS) set `config.MailNone`.
- The meeting tables (`irma_meet_sessions`, `irma_meet_participants`) keep the
  Laravel schema so existing databases work; migration 1 is `IF NOT EXISTS`.
