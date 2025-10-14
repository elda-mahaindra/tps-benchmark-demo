package config

import (
	"github.com/spf13/viper"
)

// bindEnvironmentVariables sets up environment variable mappings for all config fields.
func bindEnvironmentVariables() {
	// App config

	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.host", "APP_HOST")
	viper.BindEnv("app.env", "APP_ENV")
	viper.BindEnv("app.port.grpc", "APP_PORT_GRPC")
	viper.BindEnv("app.port.rest", "APP_PORT_REST")

	// External service config

	viper.BindEnv("external_service.go_switching.name", "GO_SWITCHING_NAME")
	viper.BindEnv("external_service.go_switching.host", "GO_SWITCHING_HOST")
	viper.BindEnv("external_service.go_switching.port", "GO_SWITCHING_PORT")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
