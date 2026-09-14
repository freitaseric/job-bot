package database

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// ErrModelNotBound indica que um Active Record foi criado sem
// estar associado a uma conexão com o banco.
var ErrModelNotBound = errors.New("model is not bound to a database")

// DB representa a conexão da aplicação com o MongoDB.
//
// Uma única instância deve ser criada durante o startup e reutilizada
// por toda a aplicação. O mongo.Client já possui seu próprio pool
// interno de conexões.
type DB struct {
	client   *mongo.Client
	database *mongo.Database
}

// Open cria a conexão com o MongoDB e verifica se o servidor está acessível.
//
// O Ping faz a aplicação falhar durante o startup caso o MongoDB
// não possa ser alcançado.
func Open(
	ctx context.Context,
	uri string,
	databaseName string,
) (*DB, error) {
	client, err := mongo.Connect(
		options.Client().
			ApplyURI(uri).
			SetAppName("job-bot"),
	)
	if err != nil {
		return nil, fmt.Errorf("create mongodb client: %w", err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		closeCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		_ = client.Disconnect(closeCtx)

		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	return &DB{
		client:   client,
		database: client.Database(databaseName),
	}, nil
}

// Collection retorna uma collection do banco configurado.
//
// Os pacotes de models utilizam este método para executar suas
// operações de persistência.
func (db *DB) Collection(name string) *mongo.Collection {
	return db.database.Collection(name)
}

// Close encerra o client MongoDB e libera seus recursos.
func (db *DB) Close(ctx context.Context) error {
	if err := db.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("disconnect mongodb: %w", err)
	}

	return nil
}
