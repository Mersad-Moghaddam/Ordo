package mysql

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestTransactionManagerWithTransaction(testingSuite *testing.T) {
	testCases := []struct {
		testName         string
		workloadError    error
		expectationError bool
	}{
		{testName: "commit success", workloadError: nil, expectationError: false},
		{testName: "rollback success", workloadError: errors.New("workload failure"), expectationError: true},
	}

	for _, testCase := range testCases {
		testingSuite.Run(testCase.testName, func(testingContext *testing.T) {
			databaseConnection, sqlMock, mockError := sqlmock.New()
			if mockError != nil {
				testingContext.Fatalf("sqlmock create failure: %v", mockError)
			}
			defer databaseConnection.Close()

			transactionManager := NewTransactionManager(databaseConnection)
			sqlMock.ExpectBegin()
			if testCase.workloadError == nil {
				sqlMock.ExpectCommit()
			} else {
				sqlMock.ExpectRollback()
			}

			executionError := transactionManager.WithTransaction(context.Background(), func(transactionContext context.Context) error {
				_, hasTransaction := TransactionFromContext(transactionContext)
				if !hasTransaction {
					testingContext.Fatalf("transaction missing from context")
				}
				return testCase.workloadError
			})
			if testCase.expectationError && executionError == nil {
				testingContext.Fatalf("expected error but got nil")
			}
			if !testCase.expectationError && executionError != nil {
				testingContext.Fatalf("unexpected execution error: %v", executionError)
			}
			if verificationError := sqlMock.ExpectationsWereMet(); verificationError != nil {
				testingContext.Fatalf("sqlmock expectations failure: %v", verificationError)
			}
		})
	}
}
