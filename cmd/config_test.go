package cmd

import (
	"os"
	"strings"
	"testing"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

func TestConfigCmd(t *testing.T) {
	origForms := defaultForms
	defer func() { defaultForms = origForms }()

	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	if err := appconfig.Save(&appconfig.Config{YtDlpCommand: "yt-dlp", N: 5}); err != nil {
		t.Fatal(err)
	}

	defaultForms = &FakeFormProvider{t: t, results: []fakeFormResult{
		{value: "my-yt-dlp"},
		{value: "10"},
	}}

	if err := configCmd.RunE(nil, nil); err != nil {
		t.Fatal(err)
	}

	cfg := appconfig.Get()
	if cfg.YtDlpCommand != "my-yt-dlp" {
		t.Errorf("expected yt-dlp command 'my-yt-dlp', got %q", cfg.YtDlpCommand)
	}
	if cfg.N != 10 {
		t.Errorf("expected N=10, got %d", cfg.N)
	}
}

func TestConfigCmdKeepYtDlpCommand(t *testing.T) {
	origForms := defaultForms
	defer func() { defaultForms = origForms }()

	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	if err := appconfig.Save(&appconfig.Config{YtDlpCommand: "existing-yt-dlp", N: 3}); err != nil {
		t.Fatal(err)
	}

	defaultForms = &FakeFormProvider{t: t, results: []fakeFormResult{
		{value: ""},
		{value: "8"},
	}}

	if err := configCmd.RunE(nil, nil); err != nil {
		t.Fatal(err)
	}

	cfg := appconfig.Get()
	if cfg.YtDlpCommand != "existing-yt-dlp" {
		t.Errorf("expected yt-dlp command 'existing-yt-dlp', got %q", cfg.YtDlpCommand)
	}
	if cfg.N != 8 {
		t.Errorf("expected N=8, got %d", cfg.N)
	}
}

func TestConfigShowCmd(t *testing.T) {
	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	err := appconfig.Save(&appconfig.Config{YtDlpCommand: "custom-yt-dlp", N: 8})
	if err != nil {
		t.Fatal(err)
	}

	output := captureOutput(func() {
		rootCmd.SetArgs([]string{"config", "show"})
		rootCmd.Execute()
	})

	if !strings.Contains(output, "Config file:") {
		t.Errorf("output missing 'Config file:', got:\n%s", output)
	}
	if !strings.Contains(output, "custom-yt-dlp") {
		t.Errorf("output missing 'custom-yt-dlp', got:\n%s", output)
	}
	if !strings.Contains(output, "8") {
		t.Errorf("output missing '8', got:\n%s", output)
	}
}
