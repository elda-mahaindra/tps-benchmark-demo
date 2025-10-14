package grpc_api

import (
	pb "go-gateway/api/grpc_api/pb"
	"go-gateway/service"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Api struct {
	pb.UnimplementedGoGatewayServer

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
