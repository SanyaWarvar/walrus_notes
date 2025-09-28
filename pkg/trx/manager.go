package trx

import (
	"context"
)

type TransactionManager interface {
	Transaction(context.Context, func(ctx context.Context) error) error
}
