package main

import (
	"fmt"
	"os"

	"load-tester/api"
	"load-tester/util/errs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/sirupsen/logrus"
)

func runRestServer(logger *logrus.Logger, port int, api *api.Api) {
	const op errs.Op = "main/runRestServer"

	// Init fiber app
	app := fiber.New()

	// CORS middleware configuration
	corsConfig := cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}

	app.Use(cors.New(corsConfig))

	// Endpoint definitions
	app = api.DefineEndpoints(app)

	// start the server
	err := app.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"op":    op,
			"scope": "Listen",
			"err":   err.Error(),
		}).Errorf("failed to listen at port: %v!", port)

		os.Exit(1)
	}

	logger.Info("rest server started successfully 🚀")
}
