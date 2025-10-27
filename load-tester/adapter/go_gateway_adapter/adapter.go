package go_gateway_adapter

import (
	pb "load-tester/adapter/go_gateway_adapter/pb"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a wrapper around the grpc client
type Adapter struct {
	serviceName string

	logger *logrus.Logger
	tracer trace.Tracer

	goGatewayClient pb.GoGatewayClient
}

// NewAdapter creates a new grpc adapter
func NewAdapter(
	serviceName string,
	logger *logrus.Logger,
	tracer trace.Tracer,
	cc *grpc.ClientConn,
) *Adapter {
	goGatewayClient := pb.NewGoGatewayClient(cc)

	return &Adapter{
		serviceName: serviceName,

		logger: logger,
		tracer: tracer,

		goGatewayClient: goGatewayClient,
	}
}
