package users

import (
	"context"
	"fmt"
	"time"

	"freitaseric.com/job-bot/internal/database"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// collectionName define onde os documentos User são armazenados.
const collectionName = "users"

// User representa um usuário cadastrado no Job Bot.
//
// O database.Model fornece ID, CreatedAt, UpdatedAt e a referência
// interna necessária para que o próprio User possa se persistir.
type User struct {
	database.Model `bson:",inline"`

	DiscordID string `bson:"discord_id" json:"discord_id"`
	Username  string `bson:"username" json:"username"`
}

// Save persiste o usuário.
//
// Se o usuário ainda não possui ID, um novo documento é criado.
// Caso ele já possua ID, o documento existente é atualizado.
func (u *User) Save(ctx context.Context) error {
	db, err := u.Database()
	if err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	collection := db.Collection(collectionName)

	now := time.Now().UTC()

	if u.IsNew() {
		u.CreatedAt = now
		u.UpdatedAt = now

		result, err := collection.InsertOne(ctx, u)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return fmt.Errorf(
					"%w: %s",
					ErrDiscordIDAlreadyExists,
					u.DiscordID,
				)
			}

			return fmt.Errorf("insert user: %w", err)
		}

		id, ok := result.InsertedID.(bson.ObjectID)
		if !ok {
			return fmt.Errorf(
				"insert user: unexpected id type %T",
				result.InsertedID,
			)
		}

		u.ID = id

		return nil
	}

	u.UpdatedAt = now

	result, err := collection.ReplaceOne(
		ctx,
		bson.M{
			"_id": u.ID,
		},
		u,
	)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf(
				"%w: %s",
				ErrDiscordIDAlreadyExists,
				u.DiscordID,
			)
		}

		return fmt.Errorf("update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrNotFound
	}

	return nil
}

// Delete remove o usuário do MongoDB.
//
// Após uma exclusão bem-sucedida, o Active Record volta a ser
// considerado um objeto não persistido.
func (u *User) Delete(ctx context.Context) error {
	if u.IsNew() {
		return ErrNotPersisted
	}

	db, err := u.Database()
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	result, err := db.
		Collection(collectionName).
		DeleteOne(
			ctx,
			bson.M{
				"_id": u.ID,
			},
		)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return ErrNotFound
	}

	u.ID = bson.NilObjectID
	u.CreatedAt = time.Time{}
	u.UpdatedAt = time.Time{}

	return nil
}
