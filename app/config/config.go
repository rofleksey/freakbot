package config

import (
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

type Config struct {
	// Service name for telemetry and logs
	ServiceName string    `yaml:"service_name" env:"SERVICE_NAME" example:"freakbot" validate:"required"`
	Sentry      Sentry    `yaml:"sentry" envPrefix:"SENTRY_"`
	Log         Log       `yaml:"log" envPrefix:"LOG_"`
	Telemetry   Telemetry `yaml:"telemetry" envPrefix:"TELEMETRY_"`
	Telegram    Telegram  `yaml:"telegram" envPrefix:"TELEGRAM_"`
}

type Sentry struct {
	DSN string `yaml:"dsn" env:"DSN" example:"https://a1b2c3d4e5f6g7h8a1b2c3d4e5f6g7h8@o123456.ingest.sentry.io/1234567"`
}

type Log struct {
	// Telegram logging config
	Telegram TelegramLog `yaml:"telegram" envPrefix:"TELEGRAM_"`
}

type TelegramLog struct {
	// Chat bot token, obtain it via BotFather
	Token string `yaml:"token" env:"TOKEN" example:"1234567890:ABCdefGHIjklMNopQRstUVwxyZ-123456789"`
	// Chat ID to send messages to
	ChatID string `yaml:"chat_id" env:"CHAT_ID" example:"1001234567890"`
}

type Telemetry struct {
	// Whether to enable opentelemetry logs/metrics/traces export
	Enabled bool `yaml:"enabled" env:"ENABLED" example:"false"`
}

type Telegram struct {
	// Chat bot token, obtain it via BotFather
	Token string `yaml:"token" env:"TOKEN" example:"1234567890:ABCdefGHIjklMNopQRstUVwxyZ-123456789"`
}

type Admin struct {
	ChatID string `yaml:"chat_id" env:"CHAT_ID" example:"1234231" validate:"required"`
}

func Load(configPath string) (*Config, error) {
	var result Config

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, oops.Errorf("failed to read config file: %w", err)
	}
	if err := yaml.Unmarshal(data, &result); err != nil {
		return nil, oops.Errorf("failed to parse YAML config: %w", err)
	}

	if err := env.ParseWithOptions(&result, env.Options{ //nolint:exhaustruct
		Prefix: "FREAKBOT_",
	}); err != nil {
		return nil, oops.Errorf("failed to parse environment variables: %w", err)
	}

	if result.ServiceName == "" {
		result.ServiceName = "freakbot"
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(result); err != nil {
		return nil, oops.Errorf("failed to validate config: %w", err)
	}

	return &result, nil
}
