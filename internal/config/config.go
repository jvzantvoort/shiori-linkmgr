package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents application configuration
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Display  DisplayConfig  `mapstructure:"display"`
}

// DatabaseConfig holds database connection settings
type DatabaseConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Database       string `mapstructure:"database"`
	MaxConnections int    `mapstructure:"maxConnections"`
}

// DisplayConfig holds display preferences
type DisplayConfig struct {
	DefaultLimit int    `mapstructure:"defaultLimit"`
	DateFormat   string `mapstructure:"dateFormat"`
}

// Load loads configuration from file and environment
func Load(configFile string) (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.maxConnections", 10)
	v.SetDefault("display.defaultLimit", 50)
	v.SetDefault("display.dateFormat", "2006-01-02 15:04")

	// Config file
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}

		v.AddConfigPath(home)
		v.SetConfigName(".linkmgr")
		v.SetConfigType("yaml")
	}

	// Environment variables
	v.SetEnvPrefix("LINKMGR")
	v.AutomaticEnv()

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.Database.Host == "" {
		return fmt.Errorf("database host is required")
	}
	if c.Database.Port <= 0 || c.Database.Port > 65535 {
		return fmt.Errorf("database port must be between 1 and 65535")
	}
	if c.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if c.Database.Database == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// DefaultConfigPath returns the default config file path
func DefaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".linkmgr.yaml"
	}
	return filepath.Join(home, ".linkmgr.yaml")
}

// Save saves configuration to file
func (c *Config) Save(path string) error {
	v := viper.New()

	v.Set("database.host", c.Database.Host)
	v.Set("database.port", c.Database.Port)
	v.Set("database.user", c.Database.User)
	v.Set("database.password", c.Database.Password)
	v.Set("database.database", c.Database.Database)
	v.Set("database.maxConnections", c.Database.MaxConnections)
	v.Set("display.defaultLimit", c.Display.DefaultLimit)
	v.Set("display.dateFormat", c.Display.DateFormat)

	if err := v.WriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Set file permissions to 0600 (user read/write only)
	if err := os.Chmod(path, 0600); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	return nil
}
