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

	// Map service result to protobuf response
	// Convert balance from pgtype.Numeric to string
	var balanceStr string
	if result.Account.Balance.Valid {
		balanceFloat, _ := result.Account.Balance.Float64Value()
		balanceStr = fmt.Sprintf("%.2f", balanceFloat.Float64)
	}

	// Handle optional date fields
	var closedDateStr string
	if result.Account.ClosedDate.Valid {
		closedDateStr = result.Account.ClosedDate.Time.Format("2006-01-02")
	}

	var openedDateStr string
	if result.Account.OpenedDate.Valid {
		openedDateStr = result.Account.OpenedDate.Time.Format("2006-01-02")
	}

	var createdAtStr string
	if result.Account.CreatedAt.Valid {
		createdAtStr = result.Account.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	var updatedAtStr string
	if result.Account.UpdatedAt.Valid {
		updatedAtStr = result.Account.UpdatedAt.Time.Format("2006-01-02T15:04:05Z07:00")
	}

	response.Account = &pb.AccountInfo{
		AccountId:     result.Account.AccountID,
		AccountNumber: result.Account.AccountNumber,
		CustomerId:    result.Account.CustomerID,
		AccountType:   result.Account.AccountType,
		AccountStatus: result.Account.AccountStatus,
		Balance:       balanceStr,
		Currency:      result.Account.Currency,
		OpenedDate:    openedDateStr,
		ClosedDate:    closedDateStr,
		CreatedAt:     createdAtStr,
		UpdatedAt:     updatedAtStr,
	}

	// Handle optional customer text fields
	var dateOfBirthStr string
	if result.Account.DateOfBirth.Valid {
		dateOfBirthStr = result.Account.DateOfBirth.Time.Format("2006-01-02")
	}

	response.Customer = &pb.CustomerInfo{
		CustomerNumber: result.Account.CustomerNumber,
		FullName:       result.Account.FullName,
		IdNumber:       result.Account.IDNumber,
		PhoneNumber:    result.Account.PhoneNumber.String,
		Email:          result.Account.Email.String,
		Address:        result.Account.Address.String,
		DateOfBirth:    dateOfBirthStr,
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.response", fmt.Sprintf("%+v", response)),
	)
	span.SetStatus(codes.Ok, "success")

	return response, nil
}
