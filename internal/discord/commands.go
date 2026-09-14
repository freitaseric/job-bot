package discord

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// commands contém os slash commands disponíveis na aplicação.
var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "Verifica se o bot está funcionando",
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
