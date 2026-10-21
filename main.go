// Command irma-meet serves Yivi Meet: video meetings in BigBlueButton whose
// host and participants prove who they are with Yivi.
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/postgresstore"

	"github.com/privacybydesign/irma-meet/internal/bbb"
	"github.com/privacybydesign/irma-meet/internal/config"
	"github.com/privacybydesign/irma-meet/internal/i18n"
	"github.com/privacybydesign/irma-meet/internal/identity"
	"github.com/privacybydesign/irma-meet/internal/mailer"
	"github.com/privacybydesign/irma-meet/internal/meeting"
	"github.com/privacybydesign/irma-meet/internal/store"
	"github.com/privacybydesign/irma-meet/internal/web"
)

const (
	commandMigrate = "migrate"

	listenAddr = ":8080"

	dbWaitTimeout   = time.Minute
	cleanupInterval = time.Hour
	sessionSweep    = 5 * time.Minute
	shutdownTimeout = 10 * time.Second
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(log, os.Args[1:]); err != nil {
		log.Error("irma-meet stopped", "error", err)
		os.Exit(1)
	}
}

// run starts the server, or with "migrate" applies the migrations and returns.
func run(log *slog.Logger, args []string) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	switch {
	case len(args) == 0:
		return runServer(ctx, log)
	case len(args) == 1 && args[0] == commandMigrate:
		return runMigrate(ctx, log)
	}

	return fmt.Errorf("usage: irma-meet [%s]", commandMigrate)
}

func runMigrate(ctx context.Context, log *slog.Logger) error {
	databaseURL, err := config.LoadDatabaseURL(os.Getenv)
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	db, err := openDB(ctx, databaseURL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	if err := store.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

	log.Info("migrations applied")

	return nil
}

func openDB(ctx context.Context, databaseURL string) (*sql.DB, error) {
	db, err := store.Open(databaseURL)
	if err != nil {
		return nil, err
	}

	waitCtx, cancel := context.WithTimeout(ctx, dbWaitTimeout)
	defer cancel()

	if err := store.WaitReady(waitCtx, db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func runServer(ctx context.Context, log *slog.Logger) error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	db, err := openDB(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer func() { _ = db.Close() }()

	text, err := i18n.New()
	if err != nil {
		return err
	}

	mails, err := mailer.New(cfg.Mail)
	if err != nil {
		return err
	}

	baseURL := cfg.AppURL.String()
	meetings := meeting.NewService(meeting.Options{
		Store:      store.New(db),
		Rooms:      bbb.New(cfg.BBBBaseURL, cfg.BBBSecret, nil),
		Mails:      mails,
		Text:       text,
		BaseURL:    baseURL,
		NameFormat: cfg.NameFormat,
	})

	sessionStore := postgresstore.NewWithCleanupInterval(db, sessionSweep)
	defer sessionStore.StopCleanup()

	sessions, err := web.NewSessionManager(sessionStore, cfg.SessionKey, cfg.AppURL)
	if err != nil {
		return err
	}

	handler, err := web.New(web.Options{
		Meetings:   meetings,
		Identities: identity.NewService(cfg.IrmaServerURL, cfg.IrmaAPIToken, nil),
		Text:       text,
		Sessions:   sessions,
		AppURL:     cfg.AppURL,
		Log:        log,
	})
	if err != nil {
		return err
	}

	go meetings.RunCleanup(ctx, cleanupInterval, log)

	return serve(ctx, log, handler)
}

func serve(ctx context.Context, log *slog.Logger, handler http.Handler) error {
	srv := &http.Server{
		Addr:              listenAddr,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       2 * time.Minute,
	}

	errc := make(chan error, 1)

	go func() {
		log.Info("listening", "addr", listenAddr)
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
