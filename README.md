# Yivi Meet

A demo of authenticated video meetings using [Yivi](https://yivi.app) and [BigBlueButton](https://bigbluebutton.org/). Users authenticate with their Yivi app before joining a meeting, giving the host certainty about who they are talking to.

> **Demo service:** sessions are limited to 60 minutes and 25 participants. Hosted on the BigBlueButton free tier (Blindside Networks).

The service is a single Go binary (`net/http`, `html/template`) with MySQL for storage. It listens on port 8080 and answers `200` on `/`.

## How it works

1. The host discloses an email address and a name with Yivi (or a SURFconext employee card for an oral exam) and fills in a meeting name.
2. The service stores the meeting under a random link id, with random BigBlueButton moderator and attendee passwords, and mails the host the link `/irma_session/join/{id}`.
3. Anyone with the link chooses host or participant and discloses with Yivi. The room is created on first use and the browser is redirected to BigBlueButton with the verified name.
4. A link works for 24 hours. A background job deletes expired meetings and their invited participants every hour.

The attributes someone disclosed are kept in their server-side session (encrypted with `SESSION_KEY`), so the IRMA server is asked for a result once per disclosure.

## Configuration

All configuration is read from the environment.

| Variable | Description |
|---|---|
| `APP_URL` | Public origin, e.g. `https://meet.example.org`. Used for the links in mails and for the CSRF check. |
| `SESSION_KEY` | 32 random bytes, base64 or hex (`openssl rand -base64 32`). Encrypts the stored sessions. Every replica must use the same key; rotating it logs everyone out. |
| `DB_HOST`, `DB_PORT`, `DB_DATABASE`, `DB_USERNAME`, `DB_PASSWORD` | MySQL connection. The user needs rights to create tables, because the service migrates the schema on start. |
| `MAIL_HOST`, `MAIL_PORT`, `MAIL_USERNAME`, `MAIL_PASSWORD` | SMTP server. |
| `MAIL_ENCRYPTION` | `tls` (STARTTLS), `ssl` (TLS on connect) or `null` (plain, STARTTLS when offered). |
| `MAIL_FROM_ADDRESS`, `MAIL_FROM_NAME` | Sender of the meeting mails. |
| `IRMA_SERVER_URL`, `IRMA_API_TOKEN` | IRMA server and the requestor token it expects in the `Authorization` header. |
| `IRMA_MEET_NAME_FORMAT` | `string` (the display name) or `json` (a JSON object with the name and the disclosed attributes) as the BigBlueButton name. |
| `BBB_SERVER_BASE_URL`, `BBB_SECRET` | BigBlueButton server, ending in `/bigbluebutton/`, and its shared secret. Calls are signed with SHA-256. |

`.env.example` lists them all.

## BigBlueButton setup

This project requires a BigBlueButton server. For local development you can use the free tier offered by Blindside Networks:

1. Register at https://registration-portal.blindsidenetworks.com/
2. After registration you will receive a server base URL and a shared secret
3. Add these to your `.env`:
   ```
   BBB_SERVER_BASE_URL=https://<your-server>/bigbluebutton/
   BBB_SECRET=<your-secret>
   ```

## Local development

Requires Docker and Docker Compose.

```bash
git clone git@github.com:privacybydesign/irma-meet.git
cd irma-meet
cp .env.docker .env
```

Fill in `SESSION_KEY` (`make key` prints one), `IRMA_API_TOKEN` and your BBB credentials in `.env`, then start all services:

```bash
make up
```

The app will be available at http://localhost:8080. Mailpit (local mail catcher) runs at http://localhost:8025 and the IRMA server at http://localhost:8088.

| Command | Description |
|---|---|
| `make up` | Build and start all containers |
| `make down` | Stop all containers |
| `make logs` | Tail the app log |
| `make test` | Run `go vet` and the tests |
| `make key` | Print a new `SESSION_KEY` |

## Running without Docker

```bash
set -a; . ./.env; set +a   # or set the variables yourself
go run .
```

Requires Go 1.26 and a MySQL database. The schema is created by goose migrations in `internal/store/migrations`, which run on start. The two meeting tables match the ones the PHP version created, so an existing database keeps working.

## Tests

```bash
go vet ./...
go test ./...
```

The tests run the whole service against fake IRMA, BigBlueButton and SMTP servers (`internal/testutil`) and an in-memory store. The SQL and the MySQL session store are tested against a real MySQL when `TEST_MYSQL_DSN` is set. It points at a server, not a database: each test creates and drops its own database, so the user needs `CREATE` and `DROP` rights.

```bash
docker run -d --name irmameet-test -e MYSQL_ROOT_PASSWORD=secret -p 3307:3306 mysql:8.0
TEST_MYSQL_DSN='root:secret@tcp(127.0.0.1:3307)/' go test ./...
```

## Static files

Only `internal/web/static` is embedded and served; nothing else from the repository is reachable over HTTP. The Bootstrap 5.3.8 bundle (`css/bootstrap.min.css`, `js/bootstrap.bundle.min.js`) is the unmodified `dist/` build from the `bootstrap` npm package, minus the source map comment. `js/irma.js` is the Yivi web frontend. Update them by replacing the files.

## Translation

The interface and mails are available in Dutch (default) and English. The translations are in `internal/i18n/locales/{nl,en}.json`, keyed by the English source text; an empty value means "use the key". The mail templates are in `internal/mailer/templates/{nl,en}`.

Text that contains markup is only rendered as HTML when it comes from a translation file (`.HTML` in a template). Everything else, including meeting names and notes, is escaped by `html/template`.
