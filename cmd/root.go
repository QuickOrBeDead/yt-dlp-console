package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/console"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/ytdlp"
	"github.com/spf13/cobra"
)

var version = "dev"

var newYtdlpExecutor = func(config *appconfig.Config) ytdlp.YtDlpExecutor {
	return ytdlp.NewYtdlpExecutor(config)
}

func getVersionFromBuildInfo() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			return info.Main.Version
		}
	}

	return ""
}

var rootCmd = &cobra.Command{
	Use:     "yt-dlp-console",
	Short:   "Interactive CLI for downloading videos using yt-dlp",
	Long:    "An interactive command-line tool for selecting and downloading videos using yt-dlp with format selection.",
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		config := appconfig.Get()
		client := ytdlp.NewYtDlpClient(newYtdlpExecutor(config), config)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		if err := runDownloadFlow(ctx, client, defaultForms); err != nil {
			console.Error("%v", err)
			return
		}
		console.Success("Download complete!")
	},
}

func runDownloadFlow(ctx context.Context, client *ytdlp.YtDlpClient, forms FormProvider) error {
	videoUrl, err := forms.Input("Video Url", func(s string) error {
		if len(strings.TrimSpace(s)) == 0 {
			return errors.New("video url is required")
		}
		if !isValidURL(s) {
			return errors.New("video url should be valid")
		}
		return nil
	})
	if err != nil {
		return err
	}

	auth, err := forms.Select("Auth", []string{"None", "Password", "Username + Password"})
	if err != nil {
		return err
	}

	var password, username, accountPassword string
	switch auth {
	case "Password":
		password, err = forms.InputPassword("Video Password")
		if err != nil {
			return err
		}
	case "Username + Password":
		username, err = forms.Input("Username", nil)
		if err != nil {
			return err
		}
		accountPassword, err = forms.InputPassword("Password")
		if err != nil {
			return err
		}
	}

	data, err := client.GetVideoData(ctx, videoUrl, password, username, accountPassword)
	if err != nil {
		return err
	}

	formats, labels := data.GetVideoList()
	if len(labels) == 0 {
		return errors.New("no video found")
	}

	videoLabel, err := forms.Select(fmt.Sprintf("Video (%s)", data.Title), labels)
	if err != nil {
		return err
	}
	videoIdx := indexOf(labels, videoLabel)
	vf := formats[videoIdx]

	var af *ytdlp.VideoFormat
	if !vf.HasAudio() {
		formats, labels = data.GetAudioList()
		if len(labels) > 0 {
			audioLabel, err := forms.Select(fmt.Sprintf("Audio (%s)", data.Title), labels)
			if err != nil {
				return err
			}
			audioIdx := indexOf(labels, audioLabel)
			af = &formats[audioIdx]
		}
	}

	formatStr := vf.FormatID
	if af != nil {
		formatStr = fmt.Sprintf("%s+%s", vf.FormatID, af.FormatID)
	}

	return client.DownloadVideo(ctx, videoUrl, password, formatStr, username, accountPassword)
}

func init() {
	if v := getVersionFromBuildInfo(); v != "" {
		version = v
	}

	rootCmd.Version = version
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func isValidURL(rawURL string) bool {
	parsedURL, err := url.ParseRequestURI(rawURL)
	return err == nil && parsedURL.Scheme != "" && parsedURL.Host != ""
}

func indexOf(slice []string, value string) int {
	for i, v := range slice {
		if v == value {
			return i
		}
	}
	return 0
}
