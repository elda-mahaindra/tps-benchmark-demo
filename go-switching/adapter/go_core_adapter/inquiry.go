package go_core_adapter

import (
	"context"
	"fmt"

	pb "go-switching/adapter/go_core_adapter/pb"
	"go-switching/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// GetAccountByAccountNumberParams defines the input parameters
type GetAccountByAccountNumberParams struct {
	AccountNumber string
}

// GetAccountByAccountNumberResult defines the output result
type GetAccountByAccountNumberResult struct {
	Account  Account
	Customer Customer
}

func (adapter *Adapter) GetAccountByAccountNumber(ctx context.Context, params *GetAccountByAccountNumberParams) (*GetAccountByAccountNumberResult, error) {
	const op = "go_core_adapter.Adapter.GetAccountByAccountNumber"

	// Start span
	ctx, span := adapter.tracer.Start(ctx, op)
	defer span.End()

	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.params", fmt.Sprintf("%+v", params)),
	)

	// Initialize result
	result := &GetAccountByAccountNumberResult{}

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, adapter.logger)
	logger = logger.WithFields(logrus.Fields{
		"[op]":   op,
		"params": fmt.Sprintf("%+v", params),
	})
	logger.Info()

	// Build gRPC request
	request := &pb.GetAccountByAccountNumberRequest{
		AccountNumber: params.AccountNumber,
	}

	// Call external service
	response, err := adapter.goCoreClient.GetAccountByAccountNumber(ctx, request)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"scope": "Get account by account number",
			"err":   err.Error(),
		}).Error()

		// Set span attributes and status
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, fmt.Errorf("error sending request: %w", err)
	}

	// Map pb response to domain models
	result.Account = Account{
		AccountID:     response.Account.AccountId,
		AccountNumber: response.Account.AccountNumber,
		CustomerID:    response.Account.CustomerId,
		AccountType:   response.Account.AccountType,
		AccountStatus: response.Account.AccountStatus,
		Balance:       response.Account.Balance,
		Currency:      response.Account.Currency,
		OpenedDate:    response.Account.OpenedDate,
		ClosedDate:    response.Account.ClosedDate,
		CreatedAt:     response.Account.CreatedAt,
		UpdatedAt:     response.Account.UpdatedAt,
	}

	result.Customer = Customer{
		CustomerNumber: response.Customer.CustomerNumber,
		FullName:       response.Customer.FullName,
		IDNumber:       response.Customer.IdNumber,
		PhoneNumber:    response.Customer.PhoneNumber,
		Email:          response.Customer.Email,
		Address:        response.Customer.Address,
		DateOfBirth:    response.Customer.DateOfBirth,
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.result", fmt.Sprintf("%+v", result)),
	)
	span.SetStatus(codes.Ok, "success")

	return result, nil
}
