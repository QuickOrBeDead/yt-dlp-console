package cmd

import (
	"fmt"
	"strconv"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/console"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := appconfig.Get()
		console.Muted("Current yt-dlp command: %q\n", cfg.YtDlpCommand)

		newCmd, _ := defaultForms.Input("yt-dlp command (blank to keep)", nil)
		if newCmd != "" {
			cfg.YtDlpCommand = newCmd
		}

		nStr, _ := defaultForms.Input("Concurrent fragments -N (1..32)", func(s string) error {
			if s == "" {
				return fmt.Errorf("N cannot be empty")
			}
			n, err := strconv.Atoi(s)
			if err != nil || n < 1 || n > 32 {
				return fmt.Errorf("enter a number between 1 and 32")
			}
			return nil
		})
		n, _ := strconv.Atoi(nStr)
		cfg.N = n

		if err := appconfig.Save(cfg); err != nil {
			return err
		}
		console.Success("Saved.")
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := appconfig.Get()
		path, _ := appconfig.ConfigFilePath()

		console.Info("Config file:    %s", path)
		console.Info("yt-dlp command: %s", cfg.YtDlpCommand)
		console.Info("Concurrent -N:  %d", cfg.N)
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	rootCmd.AddCommand(configCmd)
}
