package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"freitaseric.com/job-bot/internal/config"
	"freitaseric.com/job-bot/internal/database"
	"freitaseric.com/job-bot/internal/discord"
	"freitaseric.com/job-bot/internal/stats"
	"freitaseric.com/job-bot/internal/users"
)

func main() {
	if err := run(); err != nil {
		slog.Error(
			"application stopped",
			"error", err,
		)

		os.Exit(1)
	}
}

// run executa o ciclo de vida completo da aplicação.
//
// A ordem de inicialização garante que a presence só seja marcada
// como pronta depois que Discord, MongoDB, models e commands estiverem
// corretamente inicializados.
func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	// ------------------------------------------------------------
	// Configuration
	// ------------------------------------------------------------

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf(
			"load configuration: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Discord
	// ------------------------------------------------------------

	bot, err := discord.New(
		cfg.DiscordToken,
		cfg.DiscordAppID,
	)
	if err != nil {
		return fmt.Errorf(
			"create discord bot: %w",
			err,
		)
	}

	if err := bot.Start(); err != nil {
		return fmt.Errorf(
			"start discord bot: %w",
			err,
		)
	}

	defer func() {
		if err := bot.Close(); err != nil {
			slog.Error(
				"failed to close discord bot",
				"error", err,
			)
		}
	}()

	// Neste ponto:
	//
	// Discord conectado
	// Presence: idle / "inicializando"

	// ------------------------------------------------------------
	// MongoDB
	// ------------------------------------------------------------

	db, err := database.Open(
		ctx,
		cfg.MongoURI,
		cfg.MongoDatabase,
	)
	if err != nil {
		return fmt.Errorf(
			"initialize mongodb: %w",
			err,
		)
	}

	defer func() {
		closeCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := db.Close(closeCtx); err != nil {
			slog.Error(
				"failed to close mongodb",
				"error", err,
			)
		}
	}()

	slog.Info(
		"mongodb connected",
		"database", cfg.MongoDatabase,
	)

	// ------------------------------------------------------------
	// Active Records
	// ------------------------------------------------------------

	userCollection := users.NewCollection(db)

	if err := userCollection.EnsureIndexes(ctx); err != nil {
		return fmt.Errorf(
			"initialize users collection: %w",
			err,
		)
	}

	slog.Info(
		"users collection initialized",
	)

	// ------------------------------------------------------------
	// Application services
	// ------------------------------------------------------------

	statsService := stats.NewService(
		userCollection,
	)

	// ------------------------------------------------------------
	// Discord commands
	// ------------------------------------------------------------

	if err := bot.SyncCommands(); err != nil {
		return fmt.Errorf(
			"sync discord commands: %w",
			err,
		)
	}

	slog.Info(
		"discord commands synchronized",
	)

	// ------------------------------------------------------------
	// Initial presence
	// ------------------------------------------------------------

	if err := bot.UpdatePresence(
		ctx,
		statsService,
	); err != nil {
		return fmt.Errorf(
			"set initial ready presence: %w",
			err,
		)
	}

	// ------------------------------------------------------------
	// Background workers
	// ------------------------------------------------------------

	bot.StartPresenceScheduler(
		ctx,
		statsService,
		2*time.Minute,
	)

	// ------------------------------------------------------------
	// Ready
	// ------------------------------------------------------------

	slog.Info("job-bot is ready")

	// Mantém a aplicação viva até Ctrl+C, SIGTERM ou outro
	// cancelamento do contexto principal.
	<-ctx.Done()

	slog.Info(
		"shutdown signal received",
	)

	return nil
}
