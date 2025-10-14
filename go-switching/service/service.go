package service

import (
	"go-switching/adapter/go_core_adapter"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	logger *logrus.Logger
	tracer trace.Tracer

	goCoreAdapter *go_core_adapter.Adapter
}

func NewService(
	logger *logrus.Logger,
	tracer trace.Tracer,
	goCoreAdapter *go_core_adapter.Adapter,
) *Service {
	return &Service{
		logger: logger,
		tracer: tracer,

		goCoreAdapter: goCoreAdapter,
	}
}
