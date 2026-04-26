package ytdlp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

type DownloadResult struct {
	DefaultTemplate string `json:"_default_template"`
}

func GetVideoData(ctx context.Context, url, password string) (*VideoData, error) {
	cmd := NewYtDlpCommand(url, password)
	cmd.AddArg("-J")
	ytDlpCmd := cmd.Execute(ctx)

	var out bytes.Buffer
	var stderr bytes.Buffer
	ytDlpCmd.Stdout = &out
	ytDlpCmd.Stderr = &stderr
	err := ytDlpCmd.Run()

	cmd.ClearPassword()

	if err != nil {
		if stderr.Len() > 0 {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, stderr.String())
		} else {
			fmt.Fprintln(os.Stderr)
			fmt.Fprintln(os.Stderr, "Error running yt-dlp:", err)
		}
		return nil, err
	}

	var data VideoData
	err = json.Unmarshal(out.Bytes(), &data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing JSON:", err)
		return nil, err
	}

	return &data, nil
}

func DownloadVideo(ctx context.Context, url, password, format string) error {
	cmd := NewYtDlpCommand(url, password)
	cfg := appconfig.Get()
	if cfg.N > 0 {
		cmd.AddArgWithValue("-N", strconv.Itoa(cfg.N))
	}
	cmd.AddArg("--newline")
	cmd.AddArgWithValue("-f", format)
	cmd.AddArgWithValue("--progress-template", "%(progress)j")
	ytDlpCmd := cmd.Execute(ctx)

	cmd.ClearPassword()

	stdout, err := ytDlpCmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe:", err)
		return err
	}

	stderr, err := ytDlpCmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StderrPipe:", err)
		return err
	}

	err = ytDlpCmd.Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error running yt-dlp:", err)
		return err
	}

	downloading := false
	done := make(chan struct{})
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func(done chan struct{}) {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				if stderrScanner.Scan() {
					fmt.Fprintln(os.Stderr, stderrScanner.Text())
				} else {
					return
				}
			}
		}
	}(done)

	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		if json.Valid([]byte(line)) {
			var result DownloadResult
			if err := json.Unmarshal([]byte(line), &result); err == nil {
				downloading = true
				fmt.Fprintf(os.Stderr, "\r%s\x1b[K", result.DefaultTemplate)
			}
		} else {
			if downloading {
				downloading = false
				fmt.Fprintln(os.Stderr)
			}
			fmt.Fprintln(os.Stderr, line)
		}
	}

	if err := ytDlpCmd.Wait(); err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for command:", err)
		return err
	}

	<-done

	fmt.Println()

	return err
}
