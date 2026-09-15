package discord

import (
	"fmt"
	"log/slog"

	"freitaseric.com/job-bot/internal/discordkit"
	"github.com/bwmarrin/discordgo"
)

// Bot representa a integração da aplicação com o Discord.
type Bot struct {
	appID           string
	session         *discordgo.Session
	router          *discordkit.Router
	componentRouter *discordkit.ComponentRouter
}

// New cria e configura uma nova instância do bot.
//
// A conexão com o Discord ainda não é aberta.
// Use Start para iniciar a sessão.
func NewBot(
	token string,
	appID string,
) (*Bot, error) {
	session, err := discordgo.New(
		"Bot " + token,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create discord session: %w",
			err,
		)
	}

	bot := &Bot{
		appID:           appID,
		session:         session,
		router:          discordkit.NewRouter(),
		componentRouter: discordkit.NewComponentRouter(),
	}

	bot.registerCommands()
	if err := bot.registerComponent(); err != nil {
		slog.Error("failed to register routes", "error", err)
	}
	bot.registerHandlers()

	return bot, nil
}

// Start abre a conexão com o Discord.
//
// Após conectar, o bot entra no estado idle/inicializando até que
// o startup da aplicação seja concluído.
func (b *Bot) Start() error {
	if err := b.session.Open(); err != nil {
		return fmt.Errorf(
			"open discord session: %w",
			err,
		)
	}

	if err := b.setInitializingPresence(); err != nil {
		_ = b.session.Close()

		return fmt.Errorf(
			"set initializing presence: %w",
			err,
		)
	}

	slog.Info("discord session opened")

	return nil
}

// Close encerra a conexão com o Discord.
func (b *Bot) Close() error {
	if err := b.session.Close(); err != nil {
		return fmt.Errorf(
			"close discord session: %w",
			err,
		)
	}

	return nil
}
