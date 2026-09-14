package discord

import (
	"log/slog"

	"github.com/bwmarrin/discordgo"
)

// registerHandlers registra os eventos utilizados pelo bot.
func (b *Bot) registerHandlers() {
	b.session.AddHandlerOnce(b.onReady)
	b.session.AddHandler(b.onInteractionCreate)
}

// onReady é executado quando o Discord informa que a sessão está pronta.
func (b *Bot) onReady(
	session *discordgo.Session,
	ready *discordgo.Ready,
) {
	slog.Info(
		"discord bot logged in",
		"user", ready.User.String(),
	)
}

func (b *Bot) onInteractionCreate(
	_ *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "ping":
		b.handlePingCommand(i)

	case "configurar":
		b.handleConfigCommand(i)
	}
}
