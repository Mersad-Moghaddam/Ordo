package mysql

import (
	"context"
	"database/sql"
	"fmt"
)

type TransactionManager struct {
	databaseConnection *sql.DB
}

func NewTransactionManager(databaseConnection *sql.DB) *TransactionManager {
	return &TransactionManager{databaseConnection: databaseConnection}
}

func (transactionManager *TransactionManager) WithTransaction(requestContext context.Context, transactionWorkload func(transactionContext context.Context) error) error {
	databaseTransaction, transactionError := transactionManager.databaseConnection.BeginTx(requestContext, nil)
	if transactionError != nil {
		return fmt.Errorf("begin transaction failure: %w", transactionError)
	}

	transactionContext := context.WithValue(requestContext, transactionContextKey{}, databaseTransaction)
	workloadError := transactionWorkload(transactionContext)
	if workloadError != nil {
		rollbackError := databaseTransaction.Rollback()
		if rollbackError != nil {
			return fmt.Errorf("rollback failure after workload error: %w", rollbackError)
		}
		return workloadError
	}

	if commitError := databaseTransaction.Commit(); commitError != nil {
		return fmt.Errorf("commit transaction failure: %w", commitError)
	}
	return nil
}

type transactionContextKey struct{}

func TransactionFromContext(requestContext context.Context) (*sql.Tx, bool) {
	databaseTransaction, hasTransaction := requestContext.Value(transactionContextKey{}).(*sql.Tx)
	return databaseTransaction, hasTransaction
}
