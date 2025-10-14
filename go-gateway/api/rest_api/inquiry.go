package rest_api

import (
	"fmt"
	"time"

	"go-gateway/service"
	"go-gateway/util/logging"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (api *Api) GetAccountByAccountNumber(c *fiber.Ctx) error {
	const op = "rest_api.Api.GetAccountByAccountNumber"

	// Start span
	ctx, span := api.tracer.Start(c.Context(), op)
	defer span.End()

	span.SetAttributes(
		attribute.String("api.endpoint", "/accounts"),
		attribute.String("api.method", "GET"),
	)

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, api.logger)

	accountNumber := c.Query("account_number", "")

	// Simulate a validation operation
	time.Sleep(250 * time.Millisecond)

	params := &service.GetAccountByAccountNumberParams{
		AccountNumber: accountNumber,
	}

	logger.WithFields(logrus.Fields{
		"[op]":   op,
		"params": fmt.Sprintf("%+v", params),
	}).Info()

	result, err := api.service.GetAccountByAccountNumber(ctx, params)
	if err != nil {
		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Record success attributes
	span.SetAttributes(
		attribute.String("api.output.result", fmt.Sprintf("%+v", result)),
	)
	span.SetStatus(codes.Ok, "request completed successfully")

	return c.JSON(result)
}
