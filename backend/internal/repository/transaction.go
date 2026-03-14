package repository

import "context"

type TransactionManager interface {
	WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error
}
