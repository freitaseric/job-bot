package stats

import (
	"context"
	"fmt"
)

// UserCounter define somente a operação necessária para obter
// a quantidade de usuários.
//
// users.Collection satisfaz esta interface automaticamente porque
// já possui o método:
//
//	Count(ctx context.Context) (int64, error)
type UserCounter interface {
	Count(ctx context.Context) (int64, error)
}

// Service reúne as métricas utilizadas pela aplicação.
type Service struct {
	users UserCounter
}

// NewService cria o serviço de estatísticas.
func NewService(users UserCounter) *Service {
	return &Service{
		users: users,
	}
}

// Snapshot consulta as métricas atuais.
//
// Conforme novos models forem criados, suas contagens podem ser
// adicionadas aqui sem alterar o scheduler do Discord.
func (s *Service) Snapshot(
	ctx context.Context,
) (Snapshot, error) {
	userCount, err := s.users.Count(ctx)
	if err != nil {
		return Snapshot{}, fmt.Errorf(
			"count users: %w",
			err,
		)
	}

	return Snapshot{
		Users: userCount,
	}, nil
}
