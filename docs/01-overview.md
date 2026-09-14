# Visão geral

O Job Bot é um bot Discord escrito em Go com persistência em MongoDB.

A aplicação utiliza:

- `discordgo` para integração com o Discord;
- MongoDB Driver v2 para persistência;
- Active Record para os models;
- `context.Context` para lifecycle, timeouts e shutdown;
- goroutines e `time.Ticker` para tarefas periódicas.

## Estrutura principal

```text
job-bot/
├── cmd/
│   └── job-bot/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── discord/
│   ├── stats/
│   └── users/
├── compose.yaml
├── .env
├── .env.example
├── go.mod
└── go.sum
```

## Responsabilidades

```text
cmd/job-bot/
    Inicialização e composição das dependências.

internal/config/
    Leitura e validação das configurações.

internal/database/
    Conexão MongoDB e base dos Active Records.

internal/discord/
    Bot, commands, handlers e presence.

internal/stats/
    Agregação das métricas da aplicação.

internal/<model>/
    Active Record, queries, índices e erros do model.
```

## Fluxo de startup

```text
Config
  ↓
Discord conecta
  ↓
Presence = idle / inicializando
  ↓
MongoDB conecta e responde ao Ping
  ↓
Models e índices são inicializados
  ↓
Slash commands são sincronizados
  ↓
Stats são carregados
  ↓
Presence = online
  ↓
Schedulers são iniciados
```

O `main.go` deve continuar sendo o ponto responsável por montar as dependências.
