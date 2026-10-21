# Yivi Meet

A demo of authenticated video meetings using [Yivi](https://yivi.app) and [BigBlueButton](https://bigbluebutton.org/). Users authenticate with their Yivi app before joining a meeting, giving the host certainty about who they are talking to.

> **Demo service:** sessions are limited to 60 minutes and 25 participants. Hosted on the BigBlueButton free tier (Blindside Networks).

The service is a single Go binary (`net/http`, `html/template`) with Postgres for storage. It listens on port 8080 and answers `200` on `/`.

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
| `DATABASE_URL` | Postgres connection string, `postgresql://user:password@host:5432/dbname?connect_timeout=5`. The server needs read and write access only. Add `sslmode=require` or `verify-full` to refuse a connection without TLS; the default tries TLS and falls back to plain. |
| `MAIL_HOST`, `MAIL_PORT`, `MAIL_USERNAME`, `MAIL_PASSWORD` | SMTP server. Scaleway Transactional Email is `smtp.tem.scw.cloud` on port 587. |
| `MAIL_ENCRYPTION` | `starttls` (default, refuses to send without STARTTLS), `tls` (TLS on connect) or `none` (no TLS, for a local mail catcher). |
| `MAIL_FROM_ADDRESS`, `MAIL_FROM_NAME` | Sender of the meeting mails. `MAIL_FROM_NAME` is optional. |
| `IRMA_SERVER_URL`, `IRMA_API_TOKEN` | IRMA server and the requestor token it expects in the `Authorization` header. |
| `IRMA_MEET_NAME_FORMAT` | `string` (the display name) or `json` (a JSON object with the name and the disclosed attributes) as the BigBlueButton name. |
| `BBB_SERVER_BASE_URL`, `BBB_SECRET` | BigBlueButton server, ending in `/bigbluebutton/`, and its shared secret. Calls are signed with SHA-256. |

`.env.example` lists them all.

## Database migrations

The server never changes the schema. Run the migrations before the first start and after every update that adds one:

```bash
irma-meet migrate
```

The command reads only `DATABASE_URL`, applies the goose migrations in `internal/store/migrations` and exits. Run it with a database user that may create tables; the server itself only needs read and write access to those tables. In Kubernetes, run it as an init container with the owner user and give the server container the read/write user. Concurrent runs wait for each other.

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

`make up` also runs `irma-meet migrate` once before the app starts. The app will be available at http://localhost:8080. Mailpit (local mail catcher) runs at http://localhost:8025 and the IRMA server at http://localhost:8088.

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

Requires Go 1.26 and a Postgres database. Create the schema first, then start the server:

```bash
go run . migrate
go run .
```

## Tests

```bash
go vet ./...
go test ./...
```

The tests run the whole service against fake IRMA, BigBlueButton and SMTP servers (`internal/testutil`) and an in-memory store. The SQL, the migrations and the Postgres session store are tested against a real Postgres when `TEST_DATABASE_URL` is set. It points at a server, not a database: each test creates and drops its own database, so the user needs the right to create and drop databases.

```bash
docker run -d --name irmameet-test -e POSTGRES_PASSWORD=secret -p 5433:5432 postgres:17
TEST_DATABASE_URL='postgresql://postgres:secret@127.0.0.1:5433/postgres?sslmode=disable' go test ./...
```

`Dockerfile.test`, which CI uses, installs Postgres and starts a throwaway server, so CI runs these tests too.

## Static files

Only `internal/web/static` is embedded and served; nothing else from the repository is reachable over HTTP. The Bootstrap 5.3.8 bundle (`css/bootstrap.min.css`, `js/bootstrap.bundle.min.js`) is the unmodified `dist/` build from the `bootstrap` npm package, minus the source map comment. `js/irma.js` is the Yivi web frontend. Update them by replacing the files.

## Translation

The interface and mails are available in Dutch (default) and English. The translations are in `internal/i18n/locales/{nl,en}.json`, keyed by the English source text; an empty value means "use the key". The mail templates are in `internal/mailer/templates/{nl,en}`.

Text that contains markup is only rendered as HTML when it comes from a translation file (`.HTML` in a template). Everything else, including meeting names and notes, is escaped by `html/template`.
