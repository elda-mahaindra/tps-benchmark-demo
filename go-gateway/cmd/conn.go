package main

import (
	"fmt"

	"go-gateway/adapter/go_switching_adapter"
	"go-gateway/util/config"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createGoSwitchingAdapter(config config.GoSwitching, logger *logrus.Logger, tracer trace.Tracer) (*go_switching_adapter.Adapter, error) {
	address := fmt.Sprintf("%s:%d", config.Host, config.Port)

	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler(
			otelgrpc.WithPropagators(propagation.TraceContext{}),
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("error connecting to %s grpc server: %w", config.Name, err)
	}

	grpcAdapter := go_switching_adapter.NewAdapter(config.Name, logger, tracer, conn)

	return grpcAdapter, nil
}
