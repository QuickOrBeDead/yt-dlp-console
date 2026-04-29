package appconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndGet(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.json")

	// Override config file path for testing
	origConfig := cfg
	defer func() { cfg = origConfig }()

	cfg = nil

	// Create test config
	testCfg := &Config{
		YtDlpCommand: "/usr/local/bin/yt-dlp",
		N:            8,
	}

	// Override configFilePath by setting env or using temp
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	// Reset global config to force reload
	cfg = nil

	// Save config
	err := Save(testCfg)
	if err != nil {
		t.Fatalf("Save() error = %v", err)
	}

	// Reset and reload
	cfg = nil
	loaded := Get()

	if loaded.YtDlpCommand != testCfg.YtDlpCommand {
		t.Errorf("YtDlpCommand = %q, want %q", loaded.YtDlpCommand, testCfg.YtDlpCommand)
	}
	if loaded.N != testCfg.N {
		t.Errorf("N = %d, want %d", loaded.N, testCfg.N)
	}

	// Cleanup
	os.Remove(configPath)
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  &Config{YtDlpCommand: "yt-dlp", N: 4},
			wantErr: false,
		},
		{
			name:    "empty command",
			config:  &Config{YtDlpCommand: "", N: 4},
			wantErr: true,
		},
		{
			name:    "N too low",
			config:  &Config{YtDlpCommand: "yt-dlp", N: 0},
			wantErr: true,
		},
		{
			name:    "N too high",
			config:  &Config{YtDlpCommand: "yt-dlp", N: 33},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := defaultConfig()
	if cfg.YtDlpCommand != "yt-dlp" {
		t.Errorf("default YtDlpCommand = %q, want %q", cfg.YtDlpCommand, "yt-dlp")
	}
	if cfg.N != 4 {
		t.Errorf("default N = %d, want %d", cfg.N, 4)
	}
}
