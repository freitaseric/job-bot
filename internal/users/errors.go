package users

import "errors"

var (
	// ErrNotFound indica que o usuário solicitado não existe.
	ErrNotFound = errors.New("user not found")

	// ErrNotPersisted indica uma tentativa de executar uma operação
	// que exige que o usuário já exista no banco.
	ErrNotPersisted = errors.New("user is not persisted")

	// ErrDiscordIDAlreadyExists indica que já existe um usuário
	// cadastrado com o mesmo Discord ID.
	ErrDiscordIDAlreadyExists = errors.New("discord id already registered")
)
