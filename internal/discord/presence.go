package discord

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"freitaseric.com/job-bot/internal/stats"
	"github.com/bwmarrin/discordgo"
)

// StatsProvider representa qualquer componente capaz de fornecer
// as métricas utilizadas pela presence do bot.
type StatsProvider interface {
	Snapshot(ctx context.Context) (stats.Snapshot, error)
}

// setInitializingPresence define o estado apresentado enquanto
// a aplicação ainda está inicializando seus componentes.
func (b *Bot) setInitializingPresence() error {
	return b.session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Status: "idle",

			Activities: []*discordgo.Activity{
				{
					Name:  "Custom Status",
					State: "inicializando",
					Type:  discordgo.ActivityTypeCustom,
				},
			},
		},
	)
}

// UpdatePresence busca as métricas atuais e atualiza a presence.
//
// Este método também é utilizado para a primeira atualização,
// antes do scheduler periódico ser iniciado.
func (b *Bot) UpdatePresence(
	ctx context.Context,
	provider StatsProvider,
) error {
	queryCtx, cancel := context.WithTimeout(
		ctx,
		10*time.Second,
	)
	defer cancel()

	snapshot, err := provider.Snapshot(queryCtx)
	if err != nil {
		return fmt.Errorf(
			"load presence stats: %w",
			err,
		)
	}

	state := fmt.Sprintf(
		"%d usuários cadastrados",
		snapshot.Users,
	)

	if err := b.session.UpdateStatusComplex(
		discordgo.UpdateStatusData{
			Status: "online",

			Activities: []*discordgo.Activity{
				{
					Name:  "Custom Status",
					State: state,
					Type:  discordgo.ActivityTypeCustom,
				},
			},
		},
	); err != nil {
		return fmt.Errorf(
			"update discord presence: %w",
			err,
		)
	}

	return nil
}

// StartPresenceScheduler inicia a atualização periódica da presence.
//
// O scheduler roda em uma goroutine e é encerrado automaticamente
// quando o context recebido for cancelado.
func (b *Bot) StartPresenceScheduler(
	ctx context.Context,
	provider StatsProvider,
	interval time.Duration,
) {
	go b.runPresenceScheduler(
		ctx,
		provider,
		interval,
	)
}

// runPresenceScheduler mantém o loop periódico de atualização.
func (b *Bot) runPresenceScheduler(
	ctx context.Context,
	provider StatsProvider,
	interval time.Duration,
) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	slog.Info(
		"presence scheduler started",
		"interval", interval,
	)

	for {
		select {
		case <-ctx.Done():
			slog.Info("presence scheduler stopped")
			return

		case <-ticker.C:
			if err := b.UpdatePresence(
				ctx,
				provider,
			); err != nil {
				slog.Error(
					"failed to update presence",
					"error", err,
				)
			}
		}
	}
}
