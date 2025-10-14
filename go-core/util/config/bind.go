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

	// Store config

	viper.BindEnv("store.postgres.connection_string", "POSTGRES_CONNECTION_STRING")
	viper.BindEnv("store.postgres.pool.max_conns", "POSTGRES_MAX_CONNS")
	viper.BindEnv("store.postgres.pool.min_conns", "POSTGRES_MIN_CONNS")
	viper.BindEnv("store.postgres.pool.retry_max_attempts", "POSTGRES_POOL_RETRY_MAX_ATTEMPTS")
	viper.BindEnv("store.postgres.pool.retry_base_delay", "POSTGRES_POOL_RETRY_BASE_DELAY")
	viper.BindEnv("store.postgres.pool.retry_max_delay", "POSTGRES_POOL_RETRY_MAX_DELAY")

	// Otel tracer config

	viper.BindEnv("otel_tracer.name", "OTEL_TRACER_NAME")
	viper.BindEnv("otel_tracer.endpoint", "OTEL_TRACER_ENDPOINT")
}
