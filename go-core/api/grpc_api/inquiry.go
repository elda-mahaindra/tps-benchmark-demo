package grpc_api

import (
	"context"
	"fmt"

	"go-core/api/grpc_api/pb"
	"go-core/service/inquiry_service"
	apperrors "go-core/util/errors"
	"go-core/util/logging"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

func (api *Api) GetAccountByAccountNumber(ctx context.Context, request *pb.GetAccountByAccountNumberRequest) (*pb.GetAccountByAccountNumberResponse, error) {
	const op = "grpc_api.Api.GetAccountByAccountNumber"

	// Start span
	ctx, span := api.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.request", fmt.Sprintf("%+v", request)),
	)

	// Initialize response
	response := &pb.GetAccountByAccountNumberResponse{}

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, api.logger)
	logger = logger.WithFields(logrus.Fields{
		"[op]":    op,
		"request": fmt.Sprintf("%+v", request),
	})
	logger.Info()

	result, err := api.service.inquiry.GetAccountByAccountNumber(ctx, &inquiry_service.GetAccountByAccountNumberParams{
		AccountNumber: request.AccountNumber,
	})
	if err != nil {
		logger.WithError(err).Error()

		// Set span attributes and status
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, apperrors.ToGRPCError(err)
	}

	// Map domain models to protobuf response
	response.Account = &pb.AccountInfo{
		AccountId:     result.Account.AccountID,
		AccountNumber: result.Account.AccountNumber,
		CustomerId:    result.Account.CustomerID,
		AccountType:   result.Account.AccountType,
		AccountStatus: result.Account.AccountStatus,
		Balance:       result.Account.Balance,
		Currency:      result.Account.Currency,
		OpenedDate:    result.Account.OpenedDate,
		ClosedDate:    result.Account.ClosedDate,
		CreatedAt:     result.Account.CreatedAt,
		UpdatedAt:     result.Account.UpdatedAt,
	}

	response.Customer = &pb.CustomerInfo{
		CustomerNumber: result.Customer.CustomerNumber,
		FullName:       result.Customer.FullName,
		IdNumber:       result.Customer.IDNumber,
		PhoneNumber:    result.Customer.PhoneNumber,
		Email:          result.Customer.Email,
		Address:        result.Customer.Address,
		DateOfBirth:    result.Customer.DateOfBirth,
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.response", fmt.Sprintf("%+v", response)),
	)
	span.SetStatus(codes.Ok, "success")

	return response, nil
}
