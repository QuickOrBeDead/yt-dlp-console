package cmd

import (
	"fmt"
	"strconv"

	"charm.land/huh/v2"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure settings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := appconfig.Get()
		fmt.Printf("Current yt-dlp command: %q\n", cfg.YtDlpCommand)

		huh.NewInput().
			Title("yt-dlp command (blank to keep)").
			Value(&cfg.YtDlpCommand).
			Run()

		nStr := strconv.Itoa(cfg.N)
		huh.NewInput().
			Title("Concurrent fragments -N (1..32)").
			Validate(func(s string) error {
				if s == "" {
					return fmt.Errorf("N cannot be empty")
				}
				n, err := strconv.Atoi(s)
				if err != nil || n < 1 || n > 32 {
					return fmt.Errorf("enter a number between 1 and 32")
				}
				return nil
			}).
			Value(&nStr).
			Run()
		n, _ := strconv.Atoi(nStr)
		cfg.N = n

		if err := appconfig.Save(cfg); err != nil {
			return err
		}
		fmt.Println("Saved.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
