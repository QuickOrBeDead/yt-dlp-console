package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"runtime/debug"
	"strings"

	"charm.land/huh/v2"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/console"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/ytdlp"
	"github.com/spf13/cobra"
)

var version = "dev"

func getVersionFromBuildInfo() string {
	// Try debug.ReadBuildInfo() first
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
		var videoUrl, password, username, accountPassword string
		var af, vf *ytdlp.VideoFormat = nil, nil

		config := appconfig.Get()
		client := ytdlp.NewYtDlpClient(ytdlp.NewYtdlpExecutor(config), config)

		err := huh.NewInput().
			Title("Video Url").
			Validate(func(s string) error {
				if len(strings.TrimSpace(s)) == 0 {
					return errors.New("video url is required")
				}
				if !isValidURL(s) {
					return errors.New("video url should be valid")
				}
				return nil
			}).
			Value(&videoUrl).
			Run()
		if err != nil {
			console.Error("Error running yt-dlp: %v", err)
			return
		}

		var auth string
		err = huh.NewSelect[string]().
			Title("Auth").
			Options(huh.NewOptions("None", "Password", "Username + Password")...).
			Value(&auth).
			Run()
		if err != nil {
			console.Error("Error running yt-dlp: %v", err)
			return
		}

		switch auth {
		case "Password":
			err = huh.NewInput().
				Title("Video Password").
				EchoMode(huh.EchoModePassword).
				Value(&password).
				Run()
			if err != nil {
				console.Error("Error running yt-dlp: %v", err)
				return
			}
		case "Username + Password":
			err = huh.NewInput().
				Title("Username").
				Value(&username).
				Run()
			if err != nil {
				console.Error("Error running yt-dlp: %v", err)
				return
			}
			err = huh.NewInput().
				Title("Password").
				EchoMode(huh.EchoModePassword).
				Value(&accountPassword).
				Run()
			if err != nil {
				console.Error("Error running yt-dlp: %v", err)
				return
			}
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		data, err := client.GetVideoData(ctx, videoUrl, password, username, accountPassword)
		if err != nil {
			console.Error("Error running yt-dlp: %v", err)
			return
		}

		formats, labels := data.GetVideoList()
		if len(labels) == 0 {
			console.Error("No video found")
			return
		}

		var videoLabel string
		err = huh.NewSelect[string]().
			Title(fmt.Sprintf("Video (%s)", data.Title)).
			Options(huh.NewOptions(labels...)...).
			Value(&videoLabel).
			Run()
		if err != nil {
			console.Error("Error running yt-dlp: %v", err)
			return
		}
		videoIdx := indexOf(labels, videoLabel)
		vf = &formats[videoIdx]
		if !vf.HasAudio() {
			formats, labels = data.GetAudioList()
			if len(labels) > 0 {
				var audioLabel string
				err = huh.NewSelect[string]().
					Title(fmt.Sprintf("Audio (%s)", data.Title)).
					Options(huh.NewOptions(labels...)...).
					Value(&audioLabel).
					Run()
				if err != nil {
					console.Error("Error running yt-dlp: %v", err)
					return
				}
				audioIdx := indexOf(labels, audioLabel)
				af = &formats[audioIdx]
			}
		}

		var format string
		if af == nil {
			format = vf.FormatID
		} else {
			format = fmt.Sprintf("%s+%s", vf.FormatID, af.FormatID)
		}

		err = client.DownloadVideo(ctx, videoUrl, password, format, username, accountPassword)
		if err != nil {
			console.Error("Error downloading video: %v", err)
			return
		}
		console.Success("Download complete!")
	},
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
