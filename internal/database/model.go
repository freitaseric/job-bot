package database

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Model contém os campos comuns a todos os Active Records.
//
// Outros models podem incorporar Model usando:
//
//	type User struct {
//	    database.Model `bson:",inline"`
//	}
type Model struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`

	db *DB
}

// Bind associa o model a uma conexão de banco.
//
// Models recuperados por uma Collection devem ser vinculados
// novamente após serem decodificados do MongoDB.
func (m *Model) Bind(db *DB) {
	m.db = db
}

// Database retorna o banco associado ao Active Record.
//
// Um erro é retornado quando o model não foi criado ou carregado
// através de uma Collection válida.
func (m *Model) Database() (*DB, error) {
	if m.db == nil {
		return nil, ErrModelNotBound
	}

	return m.db, nil
}

// IsNew informa se o model ainda não foi persistido.
//
// Um ObjectID vazio significa que o documento ainda não foi inserido.
func (m *Model) IsNew() bool {
	return m.ID.IsZero()
}

// IDString retorna o ObjectID do model em formato hexadecimal.
//
// Models ainda não persistidos retornam uma string vazia.
func (m *Model) IDString() string {
	if m.ID.IsZero() {
		return ""
	}

	return m.ID.String()
}
