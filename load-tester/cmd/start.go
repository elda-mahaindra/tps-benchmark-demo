package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"load-tester/api"
	"load-tester/service"
	"load-tester/util/config"
	"load-tester/util/errs"
	"load-tester/util/tracing"

	"github.com/sirupsen/logrus"
)

func start() {
	const op errs.Op = "main/start"

	// --- Init logger ---
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)
	logger.Formatter.(*logrus.TextFormatter).DisableColors = true
	logger.Formatter.(*logrus.TextFormatter).DisableTimestamp = true
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	// --- Load config ---
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "LoadConfig",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}

	// --- Init otel tracer ---
	cleanup, err := tracing.InitTracer(config.OtelTracer.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "InitTracer",
			"err":   err.Error(),
		}).Error()
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.WithFields(logrus.Fields{
				"[op]":  op,
				"scope": "CleanupTracer",
				"err":   err.Error(),
			}).Error()
		}
	}()

	tracer := tracing.GetTracer(config.OtelTracer.Name)

	logger.WithFields(logrus.Fields{
		"[op]":   op,
		"config": fmt.Sprintf("%+v", config),
	}).Infof("Starting '%s' service ...", config.App.Name)

	// init tcp clients
	pyGatewayAdapter := createPyGatewayAdapter(logger, config.ExternalService.PyGateway)
	// Ensure TCP client resources are cleaned up on exit
	defer pyGatewayAdapter.Close()

	// init grpc clients
	goGatewayAdapter, conn, err := createGoGatewayAdapter(logger, tracer, config.ExternalService.GoGateway)
	if err != nil {
		logger.WithError(err).Error(err.Error())

		os.Exit(1)
	}
	defer conn.Close()

	// init service layer
	service := service.NewService(logger, goGatewayAdapter, pyGatewayAdapter)

	// init api layer
	restApi := api.NewApi(logger, service)

	// close kafka resources
	runRestServer(logger, config.App.Port.Rest, restApi)

	// wait for ctrl + c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// block until a signal is received
	<-ch

	logger.Info("end of program...")
}
