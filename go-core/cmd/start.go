package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go-core/api/grpc_api"
	"go-core/service/inquiry_service"
	"go-core/store/postgres_store"
	"go-core/util/config"
	"go-core/util/tracing"

	"github.com/sirupsen/logrus"
)

func start() {
	const op = "main.start"

	// --- Init logger ---
	var logger = logrus.New()
	logger.Formatter = &logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: true,
		ForceColors:      true,
	}
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	// --- Load config ---
	config, err := config.LoadConfig(".")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "Load config",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}

	// --- Init otel tracer ---
	cleanup, err := tracing.InitTracer(config.OtelTracer.Name, config.OtelTracer.Endpoint)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "Init otel tracer",
			"err":   err.Error(),
		}).Error()

		if config.App.Env != "standalone" {
			os.Exit(1)
		}
	}
	defer func() {
		if err := cleanup(context.Background()); err != nil {
			logger.WithFields(logrus.Fields{
				"[op]":  op,
				"scope": "Cleanup otel tracer",
				"err":   err.Error(),
			}).Error()
		}
	}()

	tracer := tracing.GetTracer(config.OtelTracer.Name)

	logger.WithFields(logrus.Fields{
		"[op]":   op,
		"config": fmt.Sprintf("%+v", config),
	}).Infof("Starting '%s' service ...", config.App.Name)

	// --- Create postgres pool ---
	postgresPool, err := createPostgresPoolWithRetry(config.Store.Postgres)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"[op]":  op,
			"scope": "Create postgres pool",
			"err":   err.Error(),
		}).Error()

		os.Exit(1)
	}

	// --- Init store layers ---
	postgresStore := postgres_store.NewStore(logger, tracer, postgresPool)

	// --- Init service layers ---
	inquiryService := inquiry_service.NewService(logger, tracer,  postgresStore)

	// --- Init api layers ---
	grpcApi := grpc_api.NewApi(logger, tracer, inquiryService)

	// --- Run servers ---
	runGrpcServer(config.App.Port.Grpc, grpcApi)

	// --- Wait for signal ---
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// --- Block until signal is received ---
	<-ch

	logger.Info("end of program...")
}
