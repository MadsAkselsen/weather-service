package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RapidAPIKey       string `envconfig:"X_RAPID_API_KEY"`
	WeatherAPIKey     string `envconfig:"WEATHER_API_KEY"`
	OpenWeatherMapURL string `envconfig:"OPEN_WEATHER_MAP_URL"`
	RedisURL          string `envconfig:"REDIS_ADDR"`
}

func NewConfig() *Config {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %s", err)
	}
	return &cfg
}
