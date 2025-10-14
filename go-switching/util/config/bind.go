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

	// External service config

	viper.BindEnv("external_service.go_core.name", "GO_CORE_NAME")
	viper.BindEnv("external_service.go_core.host", "GO_CORE_HOST")
	viper.BindEnv("external_service.go_core.port", "GO_CORE_PORT")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
