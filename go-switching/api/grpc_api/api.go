package grpc_api

import (
	pb "go-switching/api/grpc_api/pb"
	"go-switching/service"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Api struct {
	pb.UnimplementedGoSwitchingServer

	logger *logrus.Logger
	tracer trace.Tracer

	service *service.Service
}

func NewApi(
	logger *logrus.Logger,
	tracer trace.Tracer,
	service *service.Service,
) *Api {
	return &Api{
		logger: logger,
		tracer: tracer,

		service: service,
	}
}
