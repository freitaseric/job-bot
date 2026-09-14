package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// Config reúne as configurações necessárias para iniciar a aplicação.
//
// Os valores são carregados de variáveis de ambiente.
// Durante o desenvolvimento, um arquivo .env também pode ser utilizado.
type Config struct {
	DiscordToken string
	DiscordAppID string

	MongoURI      string
	MongoDatabase string
}

// Load carrega e valida as configurações da aplicação.
//
// O arquivo .env é opcional. Em produção, as variáveis podem ser
// fornecidas diretamente pelo ambiente.
func Load() (Config, error) {
	_ = godotenv.Load()

	cfg := Config{
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		DiscordAppID: os.Getenv("DISCORD_APP_ID"),

		MongoURI:      os.Getenv("MONGO_URI"),
		MongoDatabase: os.Getenv("MONGO_DATABASE"),
	}

	if cfg.DiscordToken == "" {
		return Config{}, errors.New("DISCORD_TOKEN is required")
	}

	if cfg.DiscordAppID == "" {
		return Config{}, errors.New("DISCORD_APP_ID is required")
	}

	if cfg.MongoURI == "" {
		return Config{}, errors.New("MONGO_URI is required")
	}

	if cfg.MongoDatabase == "" {
		return Config{}, errors.New("MONGO_DATABASE is required")
	}

	return cfg, nil
}
