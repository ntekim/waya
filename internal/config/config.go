package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig   `mapstructure:",squash"`
	Database DatabaseConfig `mapstructure:",squash"`
	Afriex   AfriexConfig   `mapstructure:",squash"`
	AI       AIConfig       `mapstructure:",squash"`
	Waya     WayaConfig     `mapstructure:",squash"`
}

type ServerConfig struct {
	Port         string        `mapstructure:"PORT"`
	Environment  string        `mapstructure:"ENV"`
	ReadTimeout  time.Duration `mapstructure:"READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"WRITE_TIMEOUT"`
}

type DatabaseConfig struct {
	Driver string `mapstructure:"DB_DRIVER"`
	Source string `mapstructure:"DB_SOURCE"`
}

type AfriexConfig struct {
	APIKey     string `mapstructure:"AFRIEX_API_KEY"`
	BaseURL    string `mapstructure:"AFRIEX_BASE_URL"`
	WebhookKey string `mapstructure:"AFRIEX_WEBHOOK_SECRET"`
}

type WayaConfig struct {
    APIKey string `mapstructure:"WAYA_API_KEY"`
    BETAWORKOSWebhookURL string `mapstructure:"BETAWORKOS_WEBHOOK_URL"` // New field
}

type AIConfig struct {
	OpenAIKey string `mapstructure:"OPENAI_API_KEY"`
}

// LoadConfig reads configuration from .env file or environment variables
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	// 1. Set Default Values (Good for local dev)
	v.SetDefault("PORT", "8080")
	v.SetDefault("ENV", "development")
	v.SetDefault("READ_TIMEOUT", 10*time.Second)
	v.SetDefault("WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("DB_DRIVER", "sqlite3")
	v.SetDefault("DB_SOURCE", "./waya.db")
	v.SetDefault("AFRIEX_BASE_URL", "https://staging.afx-server.com") // Mock URL for now

	// 2. Read from .env file
	v.AddConfigPath(path)
	v.SetConfigName("app") // looks for app.env
	v.SetConfigType("env")

	// 3. Read Automatic Environment Variables (e.g. Docker/Production)
	v.AutomaticEnv()

	// 4. Read the config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("⚠️  No .env file found, using system environment variables")
		} else {
			return nil, err
		}
	}

	// 5. Unmarshal into struct
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// 6. Basic Validation
	if cfg.Afriex.APIKey == "" && cfg.Server.Environment != "development" {
		return nil, errors.New("AFRIEX_API_KEY is required in production")
	}

	// --- Init Waya Config (Need this for Notifier and Auth) ---
    // You'll need to create a WayaConfig loader in internal/config
    // wayaCfg := config.WayaConfig{
    //     APIKey: viper.GetString("WAYA_API_KEY"),
    //     ClientWebhookURL: viper.GetString("CLIENT_WEBHOOK_URL"),
    // }

	return &cfg, nil
}