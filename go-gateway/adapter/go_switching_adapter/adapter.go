package go_switching_adapter

import (
	pb "go-gateway/adapter/go_switching_adapter/pb"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Adapter is a wrapper around the grpc client
type Adapter struct {
	serviceName string

	logger *logrus.Logger
	tracer trace.Tracer

	goSwitchingClient pb.GoSwitchingClient
}

// NewAdapter creates a new grpc adapter
func NewAdapter(
	serviceName string,
	logger *logrus.Logger,
	tracer trace.Tracer,
	cc *grpc.ClientConn,
) *Adapter {
	serviceBClient := pb.NewGoSwitchingClient(cc)

	return &Adapter{
		serviceName: serviceName,

		logger: logger,
		tracer: tracer,

		goSwitchingClient: serviceBClient,
	}
}
