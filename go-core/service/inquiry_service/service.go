package inquiry_service

import (
	"go-core/store/postgres_store"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type store struct {
	postgres postgres_store.IStore
}

type Service struct {
	logger *logrus.Logger
	tracer trace.Tracer

	store *store
}

func NewService(
	logger *logrus.Logger,
	tracer trace.Tracer,
	postgresStore postgres_store.IStore,
) *Service {
	store := &store{
		postgres: postgresStore,
	}

	service := &Service{
		logger: logger,
		tracer: tracer,

		store: store,
	}

	return service
}
