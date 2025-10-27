package api

import (
	"fmt"

	"load-tester/service"
	"load-tester/util/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type Api struct {
	logger  *logrus.Logger
	service *service.Service
}

func NewApi(logger *logrus.Logger, service *service.Service) *Api {
	return &Api{
		logger:  logger,
		service: service,
	}
}

func (api *Api) DefineEndpoints(app *fiber.App) *fiber.App {
	// Ping Routes
	ping := app.Group("/ping")
	ping.Get("/", api.Ping)

	// Load Testing Routes
	loadTest := app.Group("/test/load")
	loadTest.Post("/burst", api.loadBurst)
	loadTest.Post("/duration", api.loadDuration)
	loadTest.Post("/incremental", api.loadIncremental)
	loadTest.Post("/rps", api.loadRps)

	return app
}

// validateBaseParam validates common parameters for all test types
func (api *Api) validateBaseParam(param *service.BaseParam) error {
	if param.ServiceName == "" {
		return errs.E(errs.Validation, "service_name is required")
	}

	if param.Protocol == "" {
		return errs.E(errs.Validation, "protocol is required")
	}

	if param.Protocol != "bl2" && param.Protocol != "grpc" {
		return errs.E(errs.Validation, fmt.Sprintf("invalid protocol: %s, must be either 'bl2' or 'grpc'", param.Protocol))
	}

	// Validate payload
	if err := api.validatePayload(&param.Payload); err != nil {
		return err
	}

	return nil
}

// validatePayload validates the transaction payload
func (api *Api) validatePayload(payload *service.Payload) error {
	if payload.AccountNumber == "" {
		return errs.E(errs.Validation, "account_number is required")
	}

	return nil
}

// validateLoadBurstParam validates burst-specific parameters
func (api *Api) validateLoadBurstParam(param *service.LoadBurstParam) error {
	if err := api.validateBaseParam(&param.BaseParam); err != nil {
		return err
	}

	if param.TotalReqs <= 0 {
		return errs.E(errs.Validation, "total_reqs must be greater than 0")
	}

	return nil
}

// validateLoadRpsParam validates RPS-specific parameters
func (api *Api) validateLoadRpsParam(param *service.LoadRpsParam) error {
	if err := api.validateBaseParam(&param.BaseParam); err != nil {
		return err
	}

	if param.RPS <= 0 {
		return errs.E(errs.Validation, "rps must be greater than 0")
	}

	if param.TotalReqs <= 0 {
		return errs.E(errs.Validation, "total_reqs must be greater than 0")
	}

	return nil
}
