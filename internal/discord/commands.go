package discord

import (
	"encoding/json"
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

func (b *Bot) handlePingCommand(i *discordgo.InteractionCreate) {
	err := b.session.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pong!",
			},
		},
	)

	if err != nil {
		slog.Error("failed to respond to ping", "error", err)
	}
}

func (b *Bot) handleConfigCommand(i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options

	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		slog.Error("failed to marshal options", "error", err)
		return
	}

	err = b.session.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf(
					"```json\n%s\n```",
					string(data),
				),
			},
		},
	)

	if err != nil {
		slog.Error("failed to respond to config", "error", err)
	}
}
