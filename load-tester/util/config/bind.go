package config

import (
	"github.com/spf13/viper"
)

// bindEnvironmentVariables sets up environment variable mappings for all config fields.
func bindEnvironmentVariables() {
	// App config

	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.host", "APP_HOST")
	viper.BindEnv("app.port.rest", "APP_PORT_REST")

	// External service config

	viper.BindEnv("external_service.go_gateway.name", "GO_GATEWAY_NAME")
	viper.BindEnv("external_service.go_gateway.host", "GO_GATEWAY_HOST")
	viper.BindEnv("external_service.go_gateway.port", "GO_GATEWAY_PORT")
	viper.BindEnv("external_service.py_gateway.name", "PY_GATEWAY_NAME")
	viper.BindEnv("external_service.py_gateway.host", "PY_GATEWAY_HOST")
	viper.BindEnv("external_service.py_gateway.port", "PY_GATEWAY_PORT")
}
