package config

import (
	"time"
)

// Config holds all configuration for the application
type Config struct {
	App        App        `mapstructure:"app"`
	Store      Store      `mapstructure:"store"`
	OtelTracer OtelTracer `mapstructure:"otel_tracer"`
}

// App config

type Port struct {
	Grpc int `mapstructure:"grpc"`
}

type App struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Env  string `mapstructure:"env"`
	Port Port   `mapstructure:"port"`
}

// Store config

type PostgresPool struct {
	MaxConns         int           `mapstructure:"max_conns"`
	MinConns         int           `mapstructure:"min_conns"`
	RetryMaxAttempts int           `mapstructure:"retry_max_attempts"`
	RetryBaseDelay   time.Duration `mapstructure:"retry_base_delay"`
	RetryMaxDelay    time.Duration `mapstructure:"retry_max_delay"`
}

type Postgres struct {
	ConnectionString string       `mapstructure:"connection_string"`
	Pool             PostgresPool `mapstructure:"pool"`
}

type Store struct {
	Postgres Postgres `mapstructure:"postgres"`
}

// Otel tracer config

type OtelTracer struct {
	Name     string `mapstructure:"name"`
	Endpoint string `mapstructure:"endpoint"`
}
