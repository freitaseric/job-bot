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
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		b.router.Handle(s, i)

	case discordgo.InteractionMessageComponent,
		discordgo.InteractionModalSubmit:

		if err := b.componentRouter.Dispatch(s, i); err != nil {
			slog.Error(
				"failed to handle component interaction",
				"error", err,
			)
		}
	}
}
