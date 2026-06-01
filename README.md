# Yivi Meet

A demo of authenticated video meetings using [Yivi](https://yivi.app) and [BigBlueButton](https://bigbluebutton.org/). Users authenticate with their Yivi app before joining a meeting, giving the host certainty about who they are talking to.

> **Demo service** — sessions are limited to 60 minutes and 25 participants. Hosted on the BigBlueButton free tier (Blindside Networks).

## Requirements

- Docker and Docker Compose
- A BigBlueButton server (see below)

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

```bash
git clone git@github.com:privacybydesign/irma-meet.git
cd irma-meet
cp .env.docker .env
```

Generate an application key and insert it into `.env`:

```bash
make key
```

Fill in your BBB credentials in `.env` (see above), then start all services:

```bash
make up
```

The app will be available at http://localhost:8080. Mailpit (local mail catcher) runs at http://localhost:8025.

### Available make targets

| Command | Description |
|---|---|
| `make up` | Start all containers |
| `make down` | Stop all containers |
| `make key` | Generate and insert APP_KEY into .env |
| `make shell` | Open a shell in the app container |
| `make artisan` | Run an artisan command (e.g. `make artisan CMD=migrate`) |
| `make logs` | Tail container logs |

## Manual installation (without Docker)

```bash
git clone git@github.com:privacybydesign/irma-meet.git
cd irma-meet
cp .env.example .env
composer install
php artisan key:generate
php artisan migrate --seed
```

Requires PHP 8.2+ and a Postgres database.

## Frontend assets

Built with [Vite](https://vitejs.dev/) + Bootstrap 5. Build output lands in `public/build/` and is **not** committed — the Docker image's `assets` stage produces it during the image build. For local development against `php artisan serve`:

```bash
npm install
npm run dev    # Vite dev server with HMR
# — or —
npm run build  # one-shot production build into public/build/
```

`docker compose up --build` handles all of this automatically; the local stack does not need a host-side npm step.

## Translation

The application supports English and Dutch. Translations are in:

- `resources/lang/en.json` / `resources/lang/nl.json`
- `resources/views/emails/`

Framework-level strings (form-validation messages, auth scaffolding) are
served in English only — the `caouecs/laravel-lang` package that previously
supplied Dutch translations for those was dropped during the Laravel 10
upgrade and has not been re-added. Acceptable for the demo; revisit if a
fully-Dutch UX becomes a requirement.
