package users

import (
	"context"
	"errors"
	"fmt"

	"freitaseric.com/job-bot/internal/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Collection fornece as operações de consulta e criação de Users.
//
// O Active Record é responsável por Save e Delete.
// A Collection é responsável por localizar e criar Active Records.
type Collection struct {
	db *database.DB
}

// NewCollection cria o ponto de acesso aos Users armazenados no MongoDB.
func NewCollection(db *database.DB) *Collection {
	return &Collection{
		db: db,
	}
}

// New cria um novo User associado ao banco.
//
// O usuário ainda não é salvo automaticamente.
// Chame user.Save(ctx) para persistir o documento.
func (c *Collection) New(
	discordID string,
	username string,
) *User {
	user := &User{
		DiscordID: discordID,
		Username:  username,
	}

	user.Bind(c.db)

	return user
}

// Find procura um usuário pelo ObjectID hexadecimal.
func (c *Collection) Find(
	ctx context.Context,
	id string,
) (*User, error) {
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	return c.findOne(
		ctx,
		bson.M{
			"_id": objectID,
		},
	)
}

// FindByDiscordID procura um usuário através do seu ID do Discord.
func (c *Collection) FindByDiscordID(
	ctx context.Context,
	discordID string,
) (*User, error) {
	return c.findOne(
		ctx,
		bson.M{
			"discord_id": discordID,
		},
	)
}

// Count retorna a quantidade total de usuários cadastrados.
func (c *Collection) Count(
	ctx context.Context,
) (int64, error) {
	count, err := c.db.
		Collection(collectionName).
		CountDocuments(
			ctx,
			bson.M{},
		)

	if err != nil {
		return 0, fmt.Errorf("count users: %w", err)
	}

	return count, nil
}

// ExistsByDiscordID informa se um usuário do Discord já está cadastrado.
func (c *Collection) ExistsByDiscordID(
	ctx context.Context,
	discordID string,
) (bool, error) {
	count, err := c.db.
		Collection(collectionName).
		CountDocuments(
			ctx,
			bson.M{
				"discord_id": discordID,
			},
		)

	if err != nil {
		return false, fmt.Errorf(
			"check discord user existence: %w",
			err,
		)
	}

	return count > 0, nil
}

// findOne centraliza a leitura de um único User.
//
// Todo usuário carregado do MongoDB precisa receber novamente a
// referência do banco, pois campos privados não são armazenados no BSON.
func (c *Collection) findOne(
	ctx context.Context,
	filter any,
) (*User, error) {
	var user User

	err := c.db.
		Collection(collectionName).
		FindOne(ctx, filter).
		Decode(&user)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("find user: %w", err)
	}

	user.Bind(c.db)

	return &user, nil
}
