package service

import (
	"go-gateway/adapter/go_switching_adapter"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Service struct {
	logger *logrus.Logger
	tracer trace.Tracer

	goSwitchingAdapter *go_switching_adapter.Adapter
}

func NewService(
	logger *logrus.Logger,
	tracer trace.Tracer,
	goSwitchingAdapter *go_switching_adapter.Adapter,
) *Service {
	return &Service{
		logger: logger,
		tracer: tracer,

		goSwitchingAdapter: goSwitchingAdapter,
	}
}
