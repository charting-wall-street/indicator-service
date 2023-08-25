package config

import (
	"flag"
)

type Config struct {
	candleServiceURL string
	port             string
	allowedOrigins   string
}

func (c *Config) CandleServiceURL() string {
	return c.candleServiceURL
}

func (c *Config) Port() string {
	return c.port
}

func (c *Config) AllowedOrigins() string {
	return c.allowedOrigins
}

var serviceConfig = &Config{
	candleServiceURL: "http://localhost:9702",
	port:             "9703",
	allowedOrigins:   "*",
}

func ServiceConfig() *Config {
	return serviceConfig
}

func LoadConfig() {
	// Parse flags
	confCandleServiceUrl := flag.String("candles-url", "http://localhost:9702", "path to the candle service")
	confPort := flag.String("port", "9703", "port from which to run the service")
	confAllowedOrigins := flag.String("origins", "*", "cors origins")
	flag.Parse()

	// Set all config variables
	serviceConfig.candleServiceURL = *confCandleServiceUrl
	serviceConfig.port = *confPort
	serviceConfig.allowedOrigins = *confAllowedOrigins
}
