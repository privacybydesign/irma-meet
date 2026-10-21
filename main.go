// Command irma-meet serves Yivi Meet: video meetings in BigBlueButton whose
// host and participants prove who they are with Yivi.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alexedwards/scs/mysqlstore"

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
	listenAddr = ":8080"

	dbWaitTimeout   = time.Minute
	cleanupInterval = time.Hour
	sessionSweep    = 5 * time.Minute
	shutdownTimeout = 10 * time.Second
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := run(log); err != nil {
		log.Error("irma-meet stopped", "error", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger) error {
	cfg, err := config.Load(os.Getenv)
	if err != nil {
		return fmt.Errorf("configuration: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := store.Open(cfg.Database)
	if err != nil {
		return err
	}
	defer db.Close()

	waitCtx, cancel := context.WithTimeout(ctx, dbWaitTimeout)
	defer cancel()

	if err := store.WaitReady(waitCtx, db); err != nil {
		return err
	}

	if err := store.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrations: %w", err)
	}

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

	sessionStore := mysqlstore.NewWithCleanupInterval(db, sessionSweep)
	defer sessionStore.StopCleanup()

	sessions, err := web.NewSessionManager(sessionStore, cfg.SessionKey, cfg.AppURL.Scheme == "https")
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
