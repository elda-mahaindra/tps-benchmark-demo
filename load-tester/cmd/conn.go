package main

import (
	"fmt"

	"load-tester/adapter/go_gateway_adapter"
	"load-tester/adapter/py_gateway_adapter"
	"load-tester/util/config"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func createPyGatewayAdapter(logger *logrus.Logger, serviceConfig config.Service) *py_gateway_adapter.Adapter {
	pyGatewayAdapter := py_gateway_adapter.NewAdapter(logger, fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port))

	return pyGatewayAdapter
}

func createGoGatewayAdapter(logger *logrus.Logger, tracer trace.Tracer, serviceConfig config.Service) (*go_gateway_adapter.Adapter, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to %s grpc server: %w", serviceConfig.Name, err)
	}

	grpcClient := go_gateway_adapter.NewAdapter(serviceConfig.Name, logger, tracer, conn)

	return grpcClient, conn, nil
}
