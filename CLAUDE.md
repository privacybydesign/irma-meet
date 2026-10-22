# irma-meet repo notes

Go service (`net/http`, `html/template`, Postgres) that lets people create and join
BigBlueButton meetings after a Yivi disclosure. It replaced a Laravel app.

## Build & test

- `go vet ./... && go test ./...`. The Go version is set in `go.mod`.
- Most tests need no network, Postgres or Docker: they use the fakes in
  `internal/testutil` (IRMA server, BigBlueButton, SMTP, in-memory store).
- `internal/store`, `TestEndToEndAgainstPostgres` and `TestMigrateCommand` run
  against a real Postgres only when `TEST_DATABASE_URL` is set (a server URL,
  e.g. `postgresql://postgres:secret@127.0.0.1:5433/postgres?sslmode=disable`;
  each test makes its own database); otherwise they skip. Run them when you
  change SQL or a migration. `Dockerfile.test` (CI) starts its own Postgres.

## Layout

- `main.go` wires config, store, services and the HTTP server. `irma-meet
  migrate` applies the migrations and exits; the server never migrates.
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
- The server's database user has read/write rights only. Keep DDL out of the
  server path; schema changes are goose migrations in `internal/store/migrations`
  and run through `irma-meet migrate`.
- `MAIL_ENCRYPTION` is `starttls` (default), `tls` or `none`. The empty value
  of `config.MailEncryption` means STARTTLS, so tests against the fake SMTP
  server (no TLS) set `config.MailNone`.
- `DATABASE_URL` holds the password: never put it, or an error that quotes it,
  into a log or a config error message.
