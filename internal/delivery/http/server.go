package http

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
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
	fiberApplication.All("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))

	return &Server{fiberApplication: fiberApplication, applicationPort: applicationPort, applicationLog: applicationLog}
}

func (applicationServer *Server) Start() error {
	applicationAddress := fmt.Sprintf(":%d", applicationServer.applicationPort)
	applicationServer.applicationLog.Info("starting fiber server", zap.String("address", applicationAddress))
	return applicationServer.fiberApplication.Listen(applicationAddress)
}
