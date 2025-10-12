package postgres_store

import (
	"context"
	"errors"
	"fmt"

	"go-core/store/postgres_store/sqlc"
	apperrors "go-core/util/errors"
	"go-core/util/logging"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// TxOptions represents available transaction options
type TxOptions struct {
	Isolation  pgx.TxIsoLevel
	AccessMode pgx.TxAccessMode
	Deferrable bool
	ReadOnly   bool
}

// DefaultTxOptions returns default transaction options (Serializable, ReadWrite, NotDeferrable)
func DefaultTxOptions() TxOptions {
	return TxOptions{
		Isolation:  pgx.Serializable,
		AccessMode: pgx.ReadWrite,
		Deferrable: false,
		ReadOnly:   false,
	}
}

// ReadOnlyTxOptions returns options optimized for read-only transactions
func ReadOnlyTxOptions() TxOptions {
	return TxOptions{
		Isolation:  pgx.RepeatableRead,
		AccessMode: pgx.ReadOnly,
		Deferrable: true,
		ReadOnly:   true,
	}
}

type WithTxOptionsArgs struct {
	Opts TxOptions
	Fn   func(*sqlc.Queries) error
}

type WithTxOptionsData struct {
}

// WithTxOptions executes a function within a database transaction with custom options
func (store *Store) WithTxOptions(ctx context.Context, args *WithTxOptionsArgs) (*WithTxOptionsData, error) {
	const op = "postgres_store.Store.WithTxOptions"

	// Start span
	ctx, span := store.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.args", fmt.Sprintf("%+v", args)),
	)

	// Initialize data
	data := &WithTxOptionsData{}

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, store.logger)
	logger = logger.WithFields(logrus.Fields{
		"[op]": op,
		"args": fmt.Sprintf("%+v", args),
	})
	logger.Info()

	store.mutex.Lock()
	defer store.mutex.Unlock()

	deferrable := pgx.NotDeferrable
	if args.Opts.Deferrable {
		deferrable = pgx.Deferrable
	}

	txOptions := &pgx.TxOptions{
		IsoLevel:       args.Opts.Isolation,
		AccessMode:     args.Opts.AccessMode,
		DeferrableMode: deferrable,
	}

	tx, err := store.pool.BeginTx(ctx, *txOptions)
	if err != nil {
		appErr := apperrors.Wrap(apperrors.ErrorCodeUnavailable, err, "failed to begin transaction")

		logger.WithError(appErr).Error()

		// Set span attributes and status
		span.RecordError(appErr)
		span.SetStatus(codes.Error, appErr.Error())

		return nil, appErr
	}

	q := sqlc.New(tx)
	err = args.Fn(q)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			if errors.Is(rbErr, pgx.ErrTxClosed) {
				// Transaction already closed, return original error
				return nil, err
			}

			// Handle PostgreSQL-specific rollback errors
			var pgErr *pgconn.PgError
			if errors.As(rbErr, &pgErr) {
				// PostgreSQL error during rollback - this is a serious server issue
				appErr := apperrors.Newf(apperrors.ErrorCodeInternal,
					"PG Error Code: %s, PG Error Message: %s, RollbackTx() error, Original error: %s",
					pgErr.Code, pgErr.Message, err.Error())

				logger.WithError(appErr).Error()

				// Set span attributes and status
				span.RecordError(appErr)
				span.SetStatus(codes.Error, appErr.Error())

				return nil, appErr
			}

			// Generic rollback error
			appErr := apperrors.Wrapf(apperrors.ErrorCodeInternal, rbErr,
				"failed to rollback transaction, original error: %s", err.Error())

			logger.WithError(appErr).Error()

			// Set span attributes and status
			span.RecordError(appErr)
			span.SetStatus(codes.Error, appErr.Error())

			return nil, appErr
		}

		// Return the original error from the function (let the caller handle semantic mapping)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		// Handle commit failures - could be concurrency conflicts or connection issues
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Check for concurrency/serialization errors
			switch pgErr.Code {
			case "40001": // serialization_failure
				appErr := apperrors.Wrap(apperrors.ErrorCodeConcurrencyConflict, err, "transaction serialization failure")

				logger.WithError(appErr).Error()

				// Set span attributes and status
				span.RecordError(appErr)
				span.SetStatus(codes.Error, appErr.Error())

				return nil, appErr
			case "40P01": // deadlock_detected
				appErr := apperrors.Wrap(apperrors.ErrorCodeConcurrencyConflict, err, "transaction deadlock detected")

				logger.WithError(appErr).Error()

				// Set span attributes and status
				span.RecordError(appErr)
				span.SetStatus(codes.Error, appErr.Error())

				return nil, appErr
			default:
				// Other PostgreSQL errors during commit
				appErr := apperrors.Wrapf(apperrors.ErrorCodeInternal, err,
					"PG Error Code: %s, PG Error Message: %s", pgErr.Code, pgErr.Message)

				logger.WithError(appErr).Error()

				// Set span attributes and status
				span.RecordError(appErr)
				span.SetStatus(codes.Error, appErr.Error())

				return nil, appErr
			}
		}

		// Generic commit failure
		appErr := apperrors.Wrap(apperrors.ErrorCodeInternal, err, "failed to commit transaction")

		logger.WithError(appErr).Error()

		// Set span attributes and status
		span.RecordError(appErr)
		span.SetStatus(codes.Error, appErr.Error())

		return nil, appErr
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.data", fmt.Sprintf("%+v", data)),
	)
	span.SetStatus(codes.Ok, "success")

	return data, nil
}

type WithTxArgs struct {
	Fn func(*sqlc.Queries) error
}

type WithTxData struct {
}

// WithTx executes a function within a database transaction with default options
func (store *Store) WithTx(ctx context.Context, args *WithTxArgs) (*WithTxData, error) {
	const op = "postgres_store.Store.WithTx"

	// Start span
	ctx, span := store.tracer.Start(ctx, op)
	defer span.End()

	// Set span attributes
	span.SetAttributes(
		attribute.String("operation", op),
		attribute.String("input.args", fmt.Sprintf("%+v", args)),
	)

	// Initialize data
	data := &WithTxData{}

	// Get logger with trace id
	logger := logging.LogWithTrace(ctx, store.logger)
	logger = logger.WithFields(logrus.Fields{
		"[op]": op,
		"args": fmt.Sprintf("%+v", args),
	})
	logger.Info()

	_, err := store.WithTxOptions(ctx, &WithTxOptionsArgs{
		Opts: DefaultTxOptions(),
		Fn:   args.Fn,
	})

	if err != nil {
		logger.WithError(err).Error()

		// Set span attributes and status
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		return nil, err
	}

	// Set span attributes and status
	span.SetAttributes(
		attribute.String("output.data", fmt.Sprintf("%+v", data)),
	)
	span.SetStatus(codes.Ok, "success")

	return data, nil
}
