package service

import (
	"context"
	"fmt"

	"go-switching/adapter/go_core_adapter"
	"go-switching/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type GetAccountByAccountNumberParams struct {
	AccountNumber string
}

type GetAccountByAccountNumberResult struct {
	Account  Account
	Customer Customer
}

func (service *Service) GetAccountByAccountNumber(ctx context.Context, params *GetAccountByAccountNumberParams) (*GetAccountByAccountNumberResult, error) {
	const op = "service.Service.GetAccountByAccountNumberParams"

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

	// Call adapter
	adapterResult, err := service.goCoreAdapter.GetAccountByAccountNumber(ctx, &go_core_adapter.GetAccountByAccountNumberParams{
		AccountNumber: params.AccountNumber,
	})
	if err != nil {
		logger.WithFields(logrus.Fields{
			"scope": "Call go-core adapter",
			"err":   err.Error(),
		}).Error()

		// Record the error in the span
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// Map adapter models to service domain models
	result.Account = Account{
		AccountID:     adapterResult.Account.AccountID,
		AccountNumber: adapterResult.Account.AccountNumber,
		CustomerID:    adapterResult.Account.CustomerID,
		AccountType:   adapterResult.Account.AccountType,
		AccountStatus: adapterResult.Account.AccountStatus,
		Balance:       adapterResult.Account.Balance,
		Currency:      adapterResult.Account.Currency,
		OpenedDate:    adapterResult.Account.OpenedDate,
		ClosedDate:    adapterResult.Account.ClosedDate,
		CreatedAt:     adapterResult.Account.CreatedAt,
		UpdatedAt:     adapterResult.Account.UpdatedAt,
	}

	result.Customer = Customer{
		CustomerNumber: adapterResult.Customer.CustomerNumber,
		FullName:       adapterResult.Customer.FullName,
		IDNumber:       adapterResult.Customer.IDNumber,
		PhoneNumber:    adapterResult.Customer.PhoneNumber,
		Email:          adapterResult.Customer.Email,
		Address:        adapterResult.Customer.Address,
		DateOfBirth:    adapterResult.Customer.DateOfBirth,
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.result", fmt.Sprintf("%+v", result)),
	)
	span.SetStatus(codes.Ok, "success")

	return result, nil
}
