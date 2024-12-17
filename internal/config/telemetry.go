package config

import "fmt"

type TelemetryConfig struct {
	ServiceName    string `mapstructure:"service_name"`
	ServiceVersion string `mapstructure:"service_version"`
	Environment    string `mapstructure:"environment"`
	OTLPEndpoint   string `mapstructure:"otlp_endpoint"`
}

func (c *TelemetryConfig) Validate() error {
	if c.ServiceName == "" {
		return fmt.Errorf("service name is required")
	}
	if c.ServiceVersion == "" {
		return fmt.Errorf("service version is required")
	}
	if c.Environment == "" {
		return fmt.Errorf("environment is required")
	}
	if c.OTLPEndpoint == "" {
		return fmt.Errorf("OTLP endpoint is required")
	}
	return nil
}
