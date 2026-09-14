# MongoDB e Active Record

## Conexão

A conexão central fica em:

```text
internal/database/database.go
```

Uma única instância de `database.DB` deve ser criada durante o startup e reutilizada pela aplicação.

Exemplo:

```go
db, err := database.Open(
    ctx,
    cfg.MongoURI,
    cfg.MongoDatabase,
)
if err != nil {
    return err
}
```

No shutdown:

```go
closeCtx, cancel := context.WithTimeout(
    context.Background(),
    10*time.Second,
)
defer cancel()

_ = db.Close(closeCtx)
```

## Model base

Todos os Active Records incorporam:

```go
database.Model
```

Exemplo:

```go
type User struct {
    database.Model `bson:",inline"`

    DiscordID string `bson:"discord_id" json:"discord_id"`
    Username  string `bson:"username" json:"username"`
}
```

`database.Model` fornece:

- `ID`
- `CreatedAt`
- `UpdatedAt`
- vínculo interno com o banco
- `IsNew()`
- `IDString()`

## Active Record

O próprio objeto persistido é responsável pelas operações de instância:

```go
user.Save(ctx)
user.Delete(ctx)
```

A regra de `Save` é:

```text
sem ObjectID
    ↓
InsertOne

com ObjectID
    ↓
ReplaceOne
```

## Collection

Operações que não pertencem a uma instância ficam na `Collection` do model.

Exemplos:

```go
userCollection.New(...)
userCollection.Find(...)
userCollection.FindByDiscordID(...)
userCollection.Count(...)
userCollection.ExistsByDiscordID(...)
```

Isso evita banco global e mantém as dependências explícitas.
