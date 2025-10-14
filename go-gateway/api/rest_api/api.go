package rest_api

import (
	"go-gateway/api/rest_api/rest_middleware"
	"go-gateway/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Api struct {
	logger *logrus.Logger
	tracer trace.Tracer

	service *service.Service
}

func NewApi(
	logger *logrus.Logger,
	tracer trace.Tracer,
	service *service.Service,
) *Api {
	return &Api{
		logger: logger,
		tracer: tracer,

		service: service,
	}
}

func (api *Api) SetupRoutes(app *fiber.App) *fiber.App {
	// Error handler middleware
	app.Use(rest_middleware.ErrorHandler())

	// Account Routes
	account := app.Group("/accounts")
	account.Get("/", api.GetAccountByAccountNumber)

	return app
}
