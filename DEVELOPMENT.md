# Guia de desenvolvimento

Este documento reúne as convenções de código e os padrões usados no Job Bot.

O objetivo é manter o projeto simples, previsível e fácil de expandir.

## Estrutura principal

```text
job-bot/
├── cmd/
│   └── job-bot/
│       └── main.go
│
├── internal/
│   ├── config/
│   ├── database/
│   ├── discord/
│   ├── stats/
│   └── users/
│
├── docs/
├── compose.yaml
├── .env.example
└── go.mod
```

## Responsabilidades

```text
cmd/job-bot/
    Inicialização e composição das dependências.

internal/config/
    Leitura e validação das configurações.

internal/database/
    Infraestrutura MongoDB e base dos Active Records.

internal/discord/
    Bot, commands, handlers e presence.

internal/stats/
    Agregação de métricas.

internal/<model>/
    Active Record, queries, índices e erros.
```

## Regra de arquitetura

A direção preferida é:

```text
Discord
   ↓
Service / Collection
   ↓
Active Record
   ↓
MongoDB
```

Evite acessar MongoDB diretamente dentro de handlers do Discord.

---

# Criando slash commands

Os comandos ficam no pacote:

```text
internal/discord/
```

Exemplo:

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

Os comandos devem ser sincronizados durante o startup:

```go
if err := bot.SyncCommands(); err != nil {
	return err
}
```

---

# Tratando slash commands

Use `InteractionCreate`.

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

Handler:

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

Quando os commands crescerem:

```text
discord/
├── commands.go
├── handlers.go
├── commands_jobs.go
├── commands_users.go
└── commands_admin.go
```

---

# Criando handlers de eventos

Registre handlers em um único ponto:

```go
func (b *Bot) registerHandlers() {
	b.session.AddHandlerOnce(b.onReady)

	b.session.AddHandler(b.onGuildCreate)
	b.session.AddHandler(b.onMessageCreate)
	b.session.AddHandler(b.onInteractionCreate)
}
```

Exemplo:

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

Handlers devem ser pequenos e focados em entrada/saída do Discord.

---

# Criando um novo model

Use `users` como referência.

Para um model `Source`:

```text
internal/sources/
├── source.go
├── collection.go
├── indexes.go
└── errors.go
```

## Model

```go
type Source struct {
	database.Model `bson:",inline"`

	Name    string `bson:"name" json:"name"`
	URL     string `bson:"url" json:"url"`
	Enabled bool   `bson:"enabled" json:"enabled"`
}
```

O model contém:

- estado;
- `Save`;
- `Delete`;
- comportamentos próprios da entidade.

## Collection

```go
type Collection struct {
	db *database.DB
}

func NewCollection(db *database.DB) *Collection {
	return &Collection{
		db: db,
	}
}
```

A Collection contém:

```text
New
Find
FindBy...
Count
Exists...
queries específicas
```

Exemplo:

```go
func (c *Collection) New(
	name string,
	url string,
) *Source {
	source := &Source{
		Name:    name,
		URL:     url,
		Enabled: true,
	}

	source.Bind(c.db)

	return source
}
```

---

# Active Record

Operações de instância:

```go
user.Save(ctx)
user.Delete(ctx)
```

Consultas:

```go
userCollection.Find(ctx, id)
userCollection.FindByDiscordID(ctx, discordID)
userCollection.Count(ctx)
```

Regra de `Save`:

```text
sem ObjectID
    ↓
InsertOne

com ObjectID
    ↓
ReplaceOne
```

Evite banco global. Os models devem ser vinculados a uma instância de `database.DB`.

---

# Índices

Cada model deve declarar seus índices em:

```text
indexes.go
```

Exemplo:

```go
index := mongo.IndexModel{
	Keys: bson.D{
		{
			Key:   "url",
			Value: 1,
		},
	},
	Options: options.Index().
		SetName("url_unique").
		SetUnique(true),
}
```

Os índices devem ser garantidos durante o startup.

---

# Erros

Erros públicos ficam em:

```text
errors.go
```

Exemplo:

```go
var (
	ErrNotFound = errors.New("source not found")
	ErrURLAlreadyExists = errors.New("source url already registered")
)
```

Ao propagar erros:

```go
return fmt.Errorf(
	"count users: %w",
	err,
)
```

Para testar erro conhecido:

```go
if errors.Is(err, users.ErrNotFound) {
	// ...
}
```

---

# Stats

Stats devem agregar dados de models sem acoplar o Discord ao MongoDB.

Exemplo:

```go
type UserCounter interface {
	Count(ctx context.Context) (int64, error)
}
```

`users.Collection` satisfaz essa interface automaticamente.

```go
statsService := stats.NewService(
	userCollection,
)
```

Ao adicionar novas métricas:

```go
type Snapshot struct {
	Users         int64
	Sources       int64
	Opportunities int64
}
```

---

# Presence Scheduler

Durante o startup:

```text
Discord conectado
    ↓
idle / inicializando
```

Depois que as dependências estiverem prontas:

```go
if err := bot.UpdatePresence(
	ctx,
	statsService,
); err != nil {
	return err
}
```

Scheduler:

```go
bot.StartPresenceScheduler(
	ctx,
	statsService,
	2*time.Minute,
)
```

O scheduler deve rodar em goroutine e encerrar quando o contexto principal for cancelado.

---

# Context

Toda operação de I/O deve receber `context.Context`.

```go
user.Save(ctx)
userCollection.Find(ctx, id)
userCollection.Count(ctx)
```

Use timeout em operações de background:

```go
queryCtx, cancel := context.WithTimeout(
	ctx,
	10*time.Second,
)
defer cancel()
```

---

# Logging

Use `slog` com atributos estruturados.

```go
slog.Info(
	"user registered",
	"discord_id", user.DiscordID,
)
```

Evite:

```go
slog.Info("user %s registered", username)
```

---

# Main

O `main` deve:

- carregar configuração;
- iniciar Discord;
- conectar MongoDB;
- criar Collections;
- garantir índices;
- montar Services;
- sincronizar commands;
- atualizar a primeira presence;
- iniciar schedulers;
- coordenar shutdown.

Evite regra de negócio no `main`.

---

# Comandos úteis

Formatar:

```bash
go fmt ./...
```

Testar:

```bash
go test ./...
```

Verificar:

```bash
go vet ./...
```

Executar:

```bash
go run ./cmd/job-bot
```

MongoDB local:

```bash
docker compose up -d
```

---

# Convenções gerais

Prefira:

- dependências explícitas;
- packages pequenos;
- interfaces definidas perto de quem as consome;
- `context.Context` em I/O;
- erros com contexto;
- handlers pequenos;
- abstrações só quando forem necessárias.

Evite:

- DB global;
- singletons implícitos;
- `init()` para conectar serviços;
- queries MongoDB espalhadas pelo projeto;
- lógica de negócio no Discord;
- abstrações grandes sem uso real.
