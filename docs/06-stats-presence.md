# Stats e Presence Scheduler

A presence dinâmica usa um serviço de estatísticas separado do Discord.

## Snapshot

```go
type Snapshot struct {
    Users int64
}
```

Novas métricas podem ser adicionadas:

```go
type Snapshot struct {
    Users         int64
    Sources       int64
    Opportunities int64
}
```

## Service

O serviço recebe apenas as operações necessárias.

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

## Presence

No startup:

```text
Discord conectado
    ↓
idle / inicializando
```

Depois que MongoDB, models e commands estiverem prontos:

```go
if err := bot.UpdatePresence(
    ctx,
    statsService,
); err != nil {
    return err
}
```

O bot passa para:

```text
online / N usuários cadastrados
```

## Scheduler

```go
bot.StartPresenceScheduler(
    ctx,
    statsService,
    2*time.Minute,
)
```

O scheduler roda em uma goroutine.

```text
time.Ticker
   ↓
Snapshot()
   ↓
Collection.Count()
   ↓
UpdateStatusComplex()
```

Quando o contexto principal é cancelado, a goroutine é encerrada automaticamente.

## Adicionando nova métrica

Ao criar `sources`:

1. Adicione `Sources` ao `Snapshot`.
2. Faça o `Service` receber um `SourceCounter`.
3. Chame `sourceCollection.Count(ctx)`.
4. Adicione a nova mensagem à lógica da presence.

O scheduler não precisa conhecer MongoDB diretamente.
