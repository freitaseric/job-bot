package stats

// Snapshot representa uma fotografia das métricas atuais da aplicação.
//
// Novas métricas podem ser adicionadas aqui conforme outros models,
// como jobs e sources, forem implementados.
type Snapshot struct {
	Users int64
}
