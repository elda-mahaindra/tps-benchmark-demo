package postgres_store

import (
	"context"
	"sync"

	"go-core/store/postgres_store/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type IStore interface {
	sqlc.Querier

	WithTx(ctx context.Context, args *WithTxArgs) (*WithTxData, error)
	WithTxOptions(ctx context.Context, args *WithTxOptionsArgs) (*WithTxOptionsData, error)
}

type Store struct {
	*sqlc.Queries

	logger *logrus.Logger
	tracer trace.Tracer
	mutex  sync.Mutex

	pool *pgxpool.Pool
}

func NewStore(logger *logrus.Logger, tracer trace.Tracer, pool *pgxpool.Pool) IStore {
	return &Store{
		Queries: sqlc.New(pool),

		logger: logger,
		tracer: tracer,

		pool: pool,
	}
}
