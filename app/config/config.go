package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Log       Log       `yaml:"log"`
	Telegram  Telegram  `yaml:"telegram"`
	OpenAI    OpenAI    `yaml:"openai"`
	Retrieval Retrieval `yaml:"retrieval"`
}

type OpenAI struct {
	APIKey         string `yaml:"api_key" validate:"required"`
	BaseURL        string `yaml:"base_url" validate:"required"`
	ChatModel      string `yaml:"chat_model" validate:"required"`
	EmbeddingModel string `yaml:"embedding_model" validate:"required"`
}

type Retrieval struct {
	TopK int `yaml:"top_k" validate:"required,min=1,max=200"`
}

type Log struct {
	// Telegram logging config
	Telegram TelegramLog `yaml:"telegram"`
}

type TelegramLog struct {
	// Chat bot token, obtain it via BotFather
	Token string `yaml:"token" example:"1234567890:ABCdefGHIjklMNopQRstUVwxyZ-123456789"`
	// Chat ID to send messages to
	ChatID string `yaml:"chat_id" example:"1001234567890"`
}

type Telegram struct {
	// Chat bot token, obtain it via BotFather
	Token string `yaml:"token" example:"1234567890:ABCdefGHIjklMNopQRstUVwxyZ-123456789"`
}

type Admin struct {
	ChatID string `yaml:"chat_id" example:"1234231" validate:"required"`
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

	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(result); err != nil {
		return nil, oops.Errorf("failed to validate config: %w", err)
	}

	return &result, nil
}
