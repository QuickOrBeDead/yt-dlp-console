package ytdlp

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

type DownloadResult struct {
	DefaultTemplate string `json:"_default_template"`
}

func GetVideoData(url, password string) (*VideoData, error) {
	cmd := NewYtDlpCommand(url, password)
	cmd.AddArg("-J")
	ytDlpCmd := cmd.Execute()

	var out bytes.Buffer
	var stderr bytes.Buffer
	ytDlpCmd.Stdout = &out
	ytDlpCmd.Stderr = &stderr
	err := ytDlpCmd.Run()

	if err != nil {
		if stderr.Len() > 0 {
			fmt.Println()
			fmt.Println(stderr.String())
		} else {
			fmt.Println()
			fmt.Println("Error running yt-dlp:", err)
		}
		return nil, err
	}

	var data VideoData
	err = json.Unmarshal(out.Bytes(), &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return &data, nil
}

func DownloadVideo(url, password, format string) error {
	cmd := NewYtDlpCommand(url, password)
	cfg := appconfig.Get()
	if cfg.N > 0 {
		cmd.AddArgWithValue("-N", strconv.Itoa(cfg.N))
	}
	cmd.AddArg("--newline")
	cmd.AddArgWithValue("-f", format)
	cmd.AddArgWithValue("--progress-template", "%(progress)j")
	ytDlpCmd := cmd.Execute()

	stdout, err := ytDlpCmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating StdoutPipe:", err)
		return err
	}

	stderr, err := ytDlpCmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating StderrPipe:", err)
		return err
	}

	err = ytDlpCmd.Start()

	if err != nil {
		fmt.Println("Error running yt-dlp:", err)
		return err
	}

	downloading := false
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		for stderrScanner.Scan() {
			fmt.Println(stderrScanner.Text())
		}
	}()

	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		if json.Valid([]byte(line)) {
			var result DownloadResult
			if err := json.Unmarshal([]byte(line), &result); err == nil {
				downloading = true
				fmt.Print("\033[2K\r") // Clears the entire line and moves the cursor to the beginning
				fmt.Printf("%s", result.DefaultTemplate)
			}
		} else {
			if downloading {
				downloading = false

				fmt.Println()
			}
			fmt.Println(line)
		}
	}

	// Wait for the command to finish
	if err := ytDlpCmd.Wait(); err != nil {
		fmt.Println("Error waiting for command:", err)
		return err
	}

	return err
}
