package users

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// EnsureIndexes cria os índices necessários para a collection users.
//
// A operação é segura para ser executada durante o startup.
// MongoDB reutiliza o índice existente quando sua configuração
// já corresponde à definição solicitada.
func (c *Collection) EnsureIndexes(
	ctx context.Context,
) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{
					Key:   "discord_id",
					Value: 1,
				},
			},
			Options: options.Index().
				SetName("discord_id_unique").
				SetUnique(true),
		},
	}

	_, err := c.db.
		Collection(collectionName).
		Indexes().
		CreateMany(ctx, indexes)

	if err != nil {
		return fmt.Errorf(
			"create users indexes: %w",
			err,
		)
	}

	return nil
}
