package http

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	fiberApplication *fiber.App
	applicationPort  int
	applicationLog   *zap.Logger
}

func NewServer(applicationPort int, applicationLog *zap.Logger) *Server {
	fiberApplication := fiber.New(fiber.Config{ReadTimeout: 10 * time.Second, WriteTimeout: 10 * time.Second})
	fiberApplication.Use(recover.New())
	fiberApplication.Get("/health", func(fiberContext *fiber.Ctx) error {
		return fiberContext.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok"})
	})
	fiberApplication.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	return &Server{fiberApplication: fiberApplication, applicationPort: applicationPort, applicationLog: applicationLog}
}

func (applicationServer *Server) Start() error {
	applicationAddress := fmt.Sprintf(":%d", applicationServer.applicationPort)
	applicationServer.applicationLog.Info("starting fiber server", zap.String("address", applicationAddress))
	return applicationServer.fiberApplication.Listen(applicationAddress)
}

func (applicationServer *Server) Shutdown(requestContext context.Context) error {
	shutdownCompletion := make(chan error, 1)
	go func() {
		shutdownCompletion <- applicationServer.fiberApplication.Shutdown()
	}()
	select {
	case <-requestContext.Done():
		return requestContext.Err()
	case shutdownError := <-shutdownCompletion:
		return shutdownError
	}
}
