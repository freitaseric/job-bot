# Handlers de eventos

Todos os handlers do Discord devem ser registrados em um ponto central.

Exemplo:

```go
func (b *Bot) registerHandlers() {
    b.session.AddHandlerOnce(b.onReady)

    b.session.AddHandler(b.onGuildCreate)
    b.session.AddHandler(b.onMessageCreate)
    b.session.AddHandler(b.onInteractionCreate)
}
```

## Ready

```go
func (b *Bot) onReady(
    session *discordgo.Session,
    ready *discordgo.Ready,
) {
    slog.Info(
        "discord bot logged in",
        "user", ready.User.String(),
    )
}
```

## GuildCreate

```go
func (b *Bot) onGuildCreate(
    s *discordgo.Session,
    event *discordgo.GuildCreate,
) {
    slog.Info(
        "bot joined guild",
        "guild", event.Name,
        "id", event.ID,
    )
}
```

## MessageCreate

```go
func (b *Bot) onMessageCreate(
    s *discordgo.Session,
    event *discordgo.MessageCreate,
) {
    if event.Author.Bot {
        return
    }

    // Processamento do evento.
}
```

## Direção recomendada

```text
Evento Discord
      ↓
Handler
      ↓
Service / Collection
      ↓
Active Record
      ↓
MongoDB
```

Handlers devem ficar focados em entrada e saída do Discord, sem concentrar regras de persistência ou lógica de negócio.
