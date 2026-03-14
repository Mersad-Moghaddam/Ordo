package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	deliveryhttp "github.com/ordo/backend/internal/delivery/http"
	authdelivery "github.com/ordo/backend/internal/delivery/http/auth"
	collabdelivery "github.com/ordo/backend/internal/delivery/http/collab"
	taskdelivery "github.com/ordo/backend/internal/delivery/http/task"
	workspacedelivery "github.com/ordo/backend/internal/delivery/http/workspace"
	"github.com/ordo/backend/internal/infrastructure/config"
	"github.com/ordo/backend/internal/infrastructure/logging"
	"github.com/ordo/backend/internal/infrastructure/persistence/memory"
	"github.com/ordo/backend/internal/infrastructure/security"
	authusecase "github.com/ordo/backend/internal/usecase/auth"
	collabusecase "github.com/ordo/backend/internal/usecase/collab"
	taskusecase "github.com/ordo/backend/internal/usecase/task"
	workspaceusecase "github.com/ordo/backend/internal/usecase/workspace"
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
	defer synchronizeLogger(applicationLogger)

	memoryStore := memory.NewStore()
	tokenService, tokenServiceError := security.NewHMACTokenService(security.WithTokenSecret("integration-secret"))
	if tokenServiceError != nil {
		return fmt.Errorf("token service initialization failure: %w", tokenServiceError)
	}
	passwordHasher := security.NewSHA256PasswordHasher()

	authService := authusecase.NewService(memoryStore, memoryStore, memoryStore, passwordHasher, tokenService)
	workspaceService := workspaceusecase.NewService(memoryStore, memoryStore, memoryStore, memoryStore)
	taskService := taskusecase.NewService(memoryStore, memoryStore, memoryStore)
	collabService := collabusecase.NewService(memoryStore, memoryStore, memoryStore, memoryStore)

	applicationServer := deliveryhttp.NewServer(applicationConfiguration.HTTPPort, applicationLogger)
	apiVersionOne := applicationServer.Application().Group("/api/v1")
	authdelivery.NewHandler(authService).RegisterRoutes(apiVersionOne.Group("/auth"))
	workspacedelivery.NewHandler(workspaceService).RegisterRoutes(apiVersionOne)
	taskdelivery.NewHandler(taskService).RegisterRoutes(apiVersionOne)
	collabdelivery.NewHandler(collabService).RegisterRoutes(apiVersionOne)

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
