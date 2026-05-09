package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2"
	"github.com/Masterminds/semver/v3"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/console"
	"github.com/spf13/cobra"
)

func getLatestVersion() (string, error) {
	cmd := exec.Command("go", "list", "-m", "github.com/QuickOrBeDead/yt-dlp-console@latest")
	out, err := cmd.CombinedOutput()

	console.Info("%s", string(out))

	if err == nil {
		fields := strings.Fields(string(out))
		if len(fields) >= 2 {
			return fields[1], nil
		}
	}

	console.Warning("go list command failed:")
	console.Info("Falling back to GitHub tags...")

	resp, err := http.Get("https://api.github.com/repos/QuickOrBeDead/yt-dlp-console/tags?per_page=1")
	if err != nil {
		return "", fmt.Errorf("failed to check latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned status: %d", resp.StatusCode)
	}

	var tags []struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return "", fmt.Errorf("failed to parse GitHub response: %w", err)
	}

	if len(tags) == 0 {
		return "", fmt.Errorf("no tags found in GitHub repository")
	}

	version := strings.TrimPrefix(tags[0].Name, "v")
	return version, nil
}

func isUpdateAvailable(current, latest string) (bool, error) {
	if current == "dev" {
		return true, nil
	}
	current = strings.TrimPrefix(current, "v")
	latest = strings.TrimPrefix(latest, "v")

	cur, err := semver.NewVersion(current)
	if err != nil {
		return false, fmt.Errorf("parse current version: %w", err)
	}
	lat, err := semver.NewVersion(latest)
	if err != nil {
		return false, fmt.Errorf("parse latest version: %w", err)
	}
	return lat.GreaterThan(cur), nil
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update yt-dlp-console to the latest version",
	RunE: func(cmd *cobra.Command, args []string) error {
		if version == "dev" || version == "(devel)" {
			console.Error("Update is not available for development builds (version: %s)", version)
			console.Info("Install yt-dlp-console using: go install github.com/QuickOrBeDead/yt-dlp-console@latest")
			return nil
		}

		console.Info("Checking for updates...")

		latest, err := getLatestVersion()
		if err != nil {
			console.Error("Failed to check for updates: %v", err)
			return err
		}

		available, err := isUpdateAvailable(version, latest)
		if err != nil {
			console.Error("Version comparison failed: %v", err)
			return err
		}

		if !available {
			console.Success("Already at latest version (%s)!", version)
			return nil
		}

		var confirm bool
		_ = runHuh(huh.NewConfirm().
			Title(fmt.Sprintf("Update available: %s → %s", version, latest)).
			Description("Update now?").
			Affirmative("Yes").
			Negative("No").
			Value(&confirm))

		if !confirm {
			console.Info("Update cancelled.")
			return nil
		}

		console.Info("Updating...")
		execCmd := exec.Command("go", "install", "github.com/QuickOrBeDead/yt-dlp-console@latest")
		execCmd.Stdout = os.Stdout
		execCmd.Stderr = os.Stderr

		if err := execCmd.Run(); err != nil {
			console.Error("Update failed: %v", err)
			return err
		}

		console.Success("Update complete! Restart your terminal if needed.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
