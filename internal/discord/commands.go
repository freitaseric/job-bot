package discord

import (
	"fmt"
	"log/slog"
	"time"

	"freitaseric.com/job-bot/internal/discordkit"
	"github.com/bwmarrin/discordgo"
)

// commands contém os slash commands disponíveis na aplicação.
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Verifica se o bot está funcionando",
	},
	{
		Name:        "configurar",
		Description: "Configurações do JobBot.",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "usuario",
				Description: "Configurar o JobBot para seu usuário",
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "servidor",
				Description: "Configurar o JobBot para seu servidor",
			},
		},
	},
}

// SyncCommands registra ou atualiza os slash commands globais
// da aplicação no Discord.
func (b *Bot) SyncCommands() error {
	for _, command := range commands {
		_, err := b.session.ApplicationCommandCreate(
			b.appID,
			"",
			command,
		)
		if err != nil {
			return fmt.Errorf(
				"sync command %q: %w",
				command.Name,
				err,
			)
		}

		slog.Info(
			"discord command synced",
			"command", command.Name,
		)
	}

	return nil
}

func (b *Bot) registerCommands() {
	b.router.Command(
		"ping",
		b.handlePingCommand,
	)

	b.router.Command(
		"configurar",
		b.handleConfigCommand,
	)
}

func (b *Bot) handlePingCommand(
	ctx *discordkit.Context,
) {
	if err := ctx.Defer(); err != nil {
		slog.Error(
			"failed to defer interaction",
			"error", err,
		)

		return
	}

	time.Sleep(2 * time.Second)

	if err := ctx.Edit(
		discordkit.Text("Pong depois de 2 segundos!"),
	); err != nil {
		slog.Error(
			"failed to edit interaction",
			"error", err,
		)
	}

	if err := ctx.Followup(
		discordkit.Text("Essa é uma segunda mensagem."),
	); err != nil {
		slog.Error(
			"failed to send followup",
			"error", err,
		)
	}
}

func (b *Bot) handleConfigCommand(ctx *discordkit.Context) {
	if err := ctx.Reply(
		discordkit.Text("Configuração"),
	); err != nil {
		slog.Error(
			"failed to respond to configurar",
			"error", err,
		)
	}
}
