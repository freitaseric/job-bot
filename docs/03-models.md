# Criando novos models

Use `users` como referência para qualquer novo model.

## Estrutura padrão

Para um model chamado `Source`:

```text
internal/sources/
├── source.go
├── collection.go
├── indexes.go
└── errors.go
```

## 1. Model

```go
package sources

import "job-bot/internal/database"

const collectionName = "sources"

type Source struct {
    database.Model `bson:",inline"`

    Name    string `bson:"name" json:"name"`
    URL     string `bson:"url" json:"url"`
    Enabled bool   `bson:"enabled" json:"enabled"`
}
```

O arquivo do model deve conter principalmente:

- campos;
- `Save`;
- `Delete`;
- comportamentos específicos daquela entidade.

## 2. Collection

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

- `New`
- `Find`
- queries específicas
- `Count`
- `Exists`
- operações em conjunto

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

## 3. Índices

Declare os índices necessários em `indexes.go`.

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

Os índices devem ser aplicados durante o startup.

## 4. Erros

Erros públicos do model ficam em `errors.go`.

Exemplo:

```go
var (
    ErrNotFound = errors.New("source not found")
    ErrURLAlreadyExists = errors.New("source url already registered")
)
```

## 5. Registrar no main

```go
sourceCollection := sources.NewCollection(db)

if err := sourceCollection.EnsureIndexes(ctx); err != nil {
    return err
}
```

## Convenção

```text
model.go
    estado + Save/Delete

collection.go
    New/Find/Count/queries

indexes.go
    índices MongoDB

errors.go
    erros públicos do domínio
```
