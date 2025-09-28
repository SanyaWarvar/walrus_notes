package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
)

type ContextManager struct {
}

func NewContextManager() *ContextManager {
	return &ContextManager{}
}

type txKey struct{}

func (tm *ContextManager) InjectTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func (tm *ContextManager) ExtractTx(ctx context.Context) pgx.Tx {
	if tx, ok := ctx.Value(txKey{}).(pgx.Tx); ok {
		return tx
	}
	return nil
}
