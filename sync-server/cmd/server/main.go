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

	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/auth"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/config"
	providercrypto "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/crypto"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/db"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/device"
	serverhttp "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/http"
	"github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/providers"
	settingssync "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/sync"
	syncws "github.com/KilimcininKorOglu/MiniClaude/sync-server/internal/ws"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		slog.Error("sync-server failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()
	pool, err := db.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	command := "serve"
	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	switch command {
	case "serve":
		return serve(cfg, pool)
	case "migrate":
		if len(os.Args) < 3 {
			return fmt.Errorf("usage: sync-server migrate <up|down>")
		}
		switch os.Args[2] {
		case "up":
			return db.MigrateUp(ctx, pool, cfg.MigrationsDir)
		case "down":
			return db.MigrateDown(ctx, pool, cfg.MigrationsDir)
		default:
			return fmt.Errorf("unknown migration direction %q", os.Args[2])
		}
	default:
		return fmt.Errorf("unknown command %q", command)
	}
}

func serve(cfg config.Config, pool *pgxpool.Pool) error {
	tokenManager, err := auth.NewTokenManager(cfg.JWTSecret, cfg.SessionTTL)
	if err != nil {
		return err
	}
	secretBox, err := providercrypto.NewSecretBox(cfg.ProviderSecretMasterKey)
	if err != nil {
		return err
	}
	providerStore := providers.NewStore(pool, secretBox)
	clientTokenManager, err := device.NewClientTokenManager(cfg.JWTSecret, 30*24*time.Hour)
	if err != nil {
		return err
	}
	authService := auth.NewService(auth.NewStore(pool), tokenManager)
	deviceService := device.NewService(device.NewStore(pool), clientTokenManager, cfg.AppURL)
	syncService := settingssync.NewService(settingssync.NewStore(pool, providerStore))
	hub := syncws.NewHub()

	server := &http.Server{
		Addr:              cfg.Addr,
		Handler:           serverhttp.New(pool, authService, deviceService, syncService, providerStore, hub, auth.CookieConfig{Secure: cfg.CookieSecure}),
		ReadHeaderTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("sync-server listening", "addr", cfg.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case signal := <-signals:
		slog.Info("sync-server shutting down", "signal", signal.String())
		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()
		return server.Shutdown(ctx)
	}
}
