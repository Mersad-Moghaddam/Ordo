package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	deliveryhttp "github.com/ordo/backend/internal/delivery/http"
	"github.com/ordo/backend/internal/infrastructure/config"
	"github.com/ordo/backend/internal/infrastructure/logging"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	if applicationError := runApplication(); applicationError != nil {
		_, _ = os.Stderr.WriteString(applicationError.Error())
		os.Exit(1)
	}
}

func runApplication() error {
	applicationConfiguration, configurationError := config.NewApplicationConfiguration(config.WithEnvironment())
	if configurationError != nil {
		return fmt.Errorf("configuration failure: %w", configurationError)
	}

	applicationLogger, loggerError := logging.NewZapLogger()
	if loggerError != nil {
		return fmt.Errorf("logger failure: %w", loggerError)
	}
	defer func() {
		synchronizationError := applicationLogger.Sync()
		if synchronizationError != nil {
			_, _ = os.Stderr.WriteString(synchronizationError.Error())
		}
	}()

	applicationServer := deliveryhttp.NewServer(applicationConfiguration.HTTPPort, applicationLogger)
	requestContext, cancelFunction := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunction()

	applicationErrorGroup, derivedContext := errgroup.WithContext(requestContext)
	applicationErrorGroup.Go(func() error {
		serverError := applicationServer.Start()
		if serverError != nil && !errors.Is(serverError, context.Canceled) {
			return fmt.Errorf("fiber startup failure: %w", serverError)
		}
		return nil
	})
	applicationErrorGroup.Go(func() error {
		<-derivedContext.Done()
		applicationLogger.Info("shutdown signal received", zap.String("reason", derivedContext.Err().Error()))
		return applicationServer.Shutdown(context.Background())
	})

	return applicationErrorGroup.Wait()
}

func synchronizeLogger(applicationLogger *zap.Logger) {
	synchronizationError := applicationLogger.Sync()
	if synchronizationError != nil {
		_, _ = os.Stderr.WriteString(synchronizationError.Error())
	}
}
