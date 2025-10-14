package go_core_adapter

import (
	pb "go-switching/adapter/go_core_adapter/pb"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a wrapper around the grpc client
type Adapter struct {
	serviceName string

	logger *logrus.Logger
	tracer trace.Tracer

	goCoreClient pb.GoCoreClient
}

// NewAdapter creates a new grpc adapter
func NewAdapter(
	serviceName string,
	logger *logrus.Logger,
	tracer trace.Tracer,
	cc *grpc.ClientConn,
) *Adapter {
	serviceBClient := pb.NewGoCoreClient(cc)

	return &Adapter{
		serviceName: serviceName,

		logger: logger,
		tracer: tracer,

		goCoreClient: serviceBClient,
	}
}
