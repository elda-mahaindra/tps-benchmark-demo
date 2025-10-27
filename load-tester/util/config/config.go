package config

// Config holds all configuration for the application
type Config struct {
	App             App             `mapstructure:"app"`
	ExternalService ExternalService `mapstructure:"external_service"`
	OtelTracer      OtelTracer      `mapstructure:"otel_tracer"`
}

// App config

type Port struct {
	Rest int `mapstructure:"rest"`
}

type App struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port Port   `mapstructure:"port"`
}

// External service config

type Service struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ExternalService struct {
	GoGateway Service `mapstructure:"go_gateway"`
	PyGateway Service `mapstructure:"py_gateway"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
