package appconfig

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	YtDlpCommand string `json:"ytDlpCommand"`
	N            int    `json:"N"`
}

var cfg *Config

func defaultConfig() *Config {
	return &Config{YtDlpCommand: "yt-dlp", N: 4}
}

func configFilePath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "yt-dlp-console")
	if err := os.MkdirAll(appDir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "config.json"), nil
}

func load() *Config {
	if cfg != nil {
		return cfg
	}
	cfg = defaultConfig()
	p, err := configFilePath()
	if err != nil {
		return cfg
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return cfg
	}

	var c Config
	if err := json.Unmarshal(b, &c); err != nil {
		log.Printf("Warning: could not parse config file: %v\n", err)
		return cfg
	}

	if c.YtDlpCommand != "" {
		cfg.YtDlpCommand = c.YtDlpCommand
	}
	if c.N > 0 {
		cfg.N = c.N
	}
	return cfg
}

func Get() *Config {
	return load()
}

func Save(c *Config) error {
	if c == nil {
		return fmt.Errorf("config is nil")
	}
	if err := validate(c); err != nil {
		return err
	}
	p, err := configFilePath()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(p, b, 0o644); err != nil {
		return err
	}
	cfg = c
	return nil
}

func validate(c *Config) error {
	if c.YtDlpCommand == "" {
		return fmt.Errorf("yt-dlp command cannot be empty")
	}
	if c.N < 1 || c.N > 32 {
		return fmt.Errorf("N must be between 1 and 32")
	}
	return nil
}
