package cmd

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/ytdlp"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var url, password string
		var af, vf *ytdlp.VideoFormat = nil, nil

		prompt := promptui.Prompt{
			Label: "Video Url",
			Validate: func(s string) error {
				if len(strings.TrimSpace(s)) == 0 {
					return errors.New("video url is required")
				}

				if !isValidURL(s) {
					return errors.New("video url should be valid")
				}

				return nil
			},
		}
	url:
		url, err := prompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(-1)
			return
		}

		if err != nil {
			fmt.Println(err)
			goto url
		}

		if len(strings.TrimSpace(url)) == 0 {
			goto url
		}

		authPrompt := promptui.Select{
			Label: "Auth",
			Items: []string{"None", "Password", "Username + Password"},
		}

		_, result, err := authPrompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(-1)
			return
		}

		switch result {
		case "Password":
			prompt = promptui.Prompt{
				Label: "Video Password",
				Mask:  '*',
			}

		password:
			password, err = prompt.Run()

			if err == promptui.ErrInterrupt {
				os.Exit(-1)
				return
			}

			if err != nil {
				fmt.Println(err)
				goto password
			}
		}

		s := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithSuffix(" Retrieving available video formats..."))
		s.Start()
		data, err := ytdlp.GetVideoData(url, password)
		s.Stop()

		if err != nil {
			fmt.Println("Error running yt-dlp:", err)
			return
		}

	choose:
		formats, labels := data.GetVideoList()

		videoPrompt := promptui.Select{
			Label: fmt.Sprintf("Video (%s)", data.Title),
			Items: labels,
		}

		i, _, err := videoPrompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(-1)
			return
		}

		if err != nil {
			fmt.Println(err)
			goto choose
		}

		vf = &formats[i]
		if !vf.HasAudio() {
			formats, labels = data.GetAudioList()
			audioPrompt := promptui.Select{
				Label: fmt.Sprintf("Audio (%s)", data.Title),
				Items: labels,
			}

			i, _, err = audioPrompt.Run()
			if err == promptui.ErrInterrupt {
				os.Exit(-1)
				return
			}

			af = &formats[i]
		}

		var format string
		if af == nil {
			format = vf.FormatID
		} else {
			format = fmt.Sprintf("%s+%s", vf.FormatID, af.FormatID)
		}

		ytdlp.DownloadVideo(url, password, format)
	},
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
