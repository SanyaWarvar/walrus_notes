package postgres

import (
	"context"
	"fmt"
)

type PGTransactionManager struct {
	db         *Pool
	ctxManager *ContextManager
}

func NewPGTransactionManager(db *Pool, ctxManager *ContextManager) *PGTransactionManager {
	return &PGTransactionManager{
		db:         db,
		ctxManager: ctxManager,
	}
}

func (tm *PGTransactionManager) Transaction(ctx context.Context, tFunc func(ctx context.Context) error) error {
	tx := tm.ctxManager.ExtractTx(ctx)
	if tx != nil {
		return tFunc(ctx)
	}
	newTx, err := tm.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("error to begin transaction : %w", err)
	}
	tx = newTx

	if err := tFunc(tm.ctxManager.InjectTx(ctx, tx)); err != nil {
		if errRollback := tx.Rollback(ctx); errRollback != nil {
			return fmt.Errorf("error to rollback transaction: %w", errRollback)
		}
		return err
	}
	// if no error, commit
	if errCommit := tx.Commit(ctx); errCommit != nil {
		return fmt.Errorf("error to commit transaction: %w", errCommit)
	}
	return nil
}
