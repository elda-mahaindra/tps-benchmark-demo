package config

// Config holds all configuration for the application
type Config struct {
	App             App             `mapstructure:"app"`
	ExternalService ExternalService `mapstructure:"external_service"`
	OtelTracer      OtelTracer      `mapstructure:"otel_tracer"`
}

// App config

type Port struct {
	Grpc int `mapstructure:"grpc"`
	Rest int `mapstructure:"rest"`
}

type App struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Env  string `mapstructure:"env"`
	Port Port   `mapstructure:"port"`
}

// External service config

type GoSwitching struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type ExternalService struct {
	GoCore GoSwitching `mapstructure:"go_switching"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
