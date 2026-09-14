# Convenções do projeto

## Dependências

Dependências devem ser explícitas.

Prefira:

```go
userCollection := users.NewCollection(db)
statsService := stats.NewService(userCollection)
```

Evite:

```text
DB global
singletons implícitos
init() para conectar serviços
```

## Context

Toda operação de I/O deve receber `context.Context`.

```go
user.Save(ctx)
userCollection.Find(ctx, id)
userCollection.Count(ctx)
```

Use timeout quando uma operação de background não pode ficar bloqueada indefinidamente.

## Erros

Adicione contexto ao propagar erros:

```go
return fmt.Errorf(
    "count users: %w",
    err,
)
```

Use erros públicos do domínio quando o chamador precisar distinguir situações específicas.

```go
errors.Is(err, users.ErrNotFound)
```

## Logging

Use `slog` com atributos estruturados:

```go
slog.Info(
    "user registered",
    "discord_id", user.DiscordID,
)
```

Não use formatação estilo `Printf`:

```go
// Evitar:
slog.Info("user %s registered", username)
```

## Active Record

Instância:

```go
user.Save(ctx)
user.Delete(ctx)
```

Collection:

```go
userCollection.New(...)
userCollection.Find(...)
userCollection.Count(...)
```

## MongoDB

Evite acessar collections diretamente fora do pacote do model.

Prefira:

```go
userCollection.FindByDiscordID(ctx, id)
```

em vez de espalhar:

```go
db.Collection("users").FindOne(...)
```

## Main

O `main` deve:

- carregar configuração;
- inicializar Discord;
- conectar MongoDB;
- criar Collections;
- garantir índices;
- montar Services;
- sincronizar comandos;
- iniciar schedulers;
- coordenar shutdown.

Não coloque regra de negócio no `main`.

## Organização

```text
database/
    infraestrutura compartilhada

discord/
    integração Discord

stats/
    agregação de métricas

users/, jobs/, sources/
    Active Records

cmd/job-bot/
    composição da aplicação
```

## Regra prática

Ao criar uma nova funcionalidade, pergunte:

```text
É persistência de uma entidade?
    → pacote do model

É integração com Discord?
    → discord/

É composição de informações de vários models?
    → service específico

É inicialização?
    → main.go
```
