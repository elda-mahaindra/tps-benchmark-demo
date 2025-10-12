package grpc_api

import (
	"go-core/api/grpc_api/pb"
	"go-core/service/inquiry_service"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type service struct {
	inquiry *inquiry_service.Service
}

type Api struct {
	pb.UnimplementedGoCoreServer

	logger *logrus.Logger
	tracer trace.Tracer

	service *service
}

func NewApi(
	logger *logrus.Logger,
	tracer trace.Tracer,
	inquiryService *inquiry_service.Service,
) *Api {
	service := &service{
		inquiry: inquiryService,
	}

	return &Api{
		logger: logger,
		tracer: tracer,

		service: service,
	}
}
