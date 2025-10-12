package inquiry_service

import (
	"context"
	"fmt"

	"go-core/store/postgres_store/sqlc"
	"go-core/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type GetAccountByAccountNumberParams struct {
	AccountNumber string
}

type GetAccountByAccountNumberResult struct {
	Account sqlc.GetAccountByAccountNumberRow
}

func (service *Service) GetAccountByAccountNumber(ctx context.Context, params *GetAccountByAccountNumberParams) (*GetAccountByAccountNumberResult, error) {
	const op = "inquiry_service.Service.GetAccountByAccountNumber"

	// Start span
	ctx, span := service.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.params", fmt.Sprintf("%+v", params)),
	)

	// Initialize result
	result := &GetAccountByAccountNumberResult{}

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, service.logger)
	logger = logger.WithFields(logrus.Fields{
		"[op]":   op,
		"params": fmt.Sprintf("%+v", params),
	})
	logger.Info()

	// Call store layer to get account with customer details
	account, err := service.store.postgres.GetAccountByAccountNumber(ctx, params.AccountNumber)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"scope": "Get account by account number",
			"err":   err.Error(),
		}).Error()

		// Set span attributes and status
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// Set result
	result.Account = account

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.result", fmt.Sprintf("%+v", result)),
	)
	span.SetStatus(codes.Ok, "success")

	return result, nil
}
