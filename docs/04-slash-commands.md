# Slash commands

Os slash commands ficam no pacote:

```text
internal/discord/
```

## Declarando um comando

Exemplo simples:

```go
var commands = []*discordgo.ApplicationCommand{
    {
        Name:        "ping",
        Description: "Verifica se o bot está funcionando",
    },
}
```

## Com argumentos

```go
{
    Name:        "user",
    Description: "Consulta um usuário",
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type:        discordgo.ApplicationCommandOptionUser,
            Name:        "usuario",
            Description: "Usuário que será consultado",
            Required:    true,
        },
    },
}
```

## Sincronização

Durante o startup:

```go
if err := bot.SyncCommands(); err != nil {
    return err
}
```

A aplicação só deve mudar para a presence de pronta depois que os comandos forem sincronizados com sucesso.

## Recebendo slash commands

O evento utilizado é `InteractionCreate`.

```go
func (b *Bot) onInteractionCreate(
    s *discordgo.Session,
    i *discordgo.InteractionCreate,
) {
    if i.Type != discordgo.InteractionApplicationCommand {
        return
    }

    switch i.ApplicationCommandData().Name {
    case "ping":
        b.handlePing(s, i)

    case "stats":
        b.handleStats(s, i)
    }
}
```

## Handler de comando

```go
func (b *Bot) handlePing(
    s *discordgo.Session,
    i *discordgo.InteractionCreate,
) {
    err := s.InteractionRespond(
        i.Interaction,
        &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Pong!",
            },
        },
    )

    if err != nil {
        slog.Error(
            "failed to respond to ping",
            "error", err,
        )
    }
}
```

## Organização

Quando houver poucos comandos:

```text
discord/
├── commands.go
└── handlers.go
```

Quando crescer:

```text
discord/
├── commands.go
├── handlers.go
├── commands_jobs.go
├── commands_users.go
└── commands_admin.go
```

Evite acessar MongoDB diretamente em handlers. Prefira usar Collections ou Services recebidos pelo bot.
