package ytdlp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/console"
)

type DownloadResult struct {
	DefaultTemplate string `json:"_default_template"`
}

type YtDlpClient struct {
	executor YtDlpExecutor
	config   *appconfig.Config
}

func NewYtDlpClient(executor YtDlpExecutor, config *appconfig.Config) *YtDlpClient {
	return &YtDlpClient{executor: executor, config: config}
}

func (c *YtDlpClient) GetVideoData(ctx context.Context, url, password string) (*VideoData, error) {
	cmd := NewYtDlpCommandArgs(url, password)
	cmd.AddArg("-J")

	var out bytes.Buffer
	var stderr bytes.Buffer
	err := c.executor.Execute(ctx, cmd, "Retrieving available video formats...", &out, &stderr)

	if stderr.Len() > 0 {
		console.Error("%s", stderr.String())
	}

	if err != nil {
		return nil, err
	}

	var data VideoData
	if err := json.Unmarshal(out.Bytes(), &data); err != nil {
		console.Error("Error parsing JSON: %v", err)
		return nil, err
	}

	return &data, nil
}

func (c *YtDlpClient) DownloadVideo(ctx context.Context, url, password, format string) error {
	cmd := NewYtDlpCommandArgs(url, password)

	if c.config.N > 0 {
		cmd.AddArgWithValue("-N", strconv.Itoa(c.config.N))
	}

	cmd.AddArg("--newline")
	cmd.AddArgWithValue("-f", format)
	cmd.AddArgWithValue("--progress-template", "%(progress)j")

	stdout, stderr, err := c.executor.ExecuteWithStreams(ctx, cmd)
	if err != nil {
		console.Error("Error running yt-dlp: %v", err)
		return err
	}

	downloading := false
	done := make(chan struct{})
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	go func() {
		defer close(done)
		for stderrScanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				console.Error("%s", stderrScanner.Text())
			}
		}
	}()

	for stdoutScanner.Scan() {
		line := stdoutScanner.Text()
		if json.Valid([]byte(line)) {
			var result DownloadResult
			if err := json.Unmarshal([]byte(line), &result); err == nil {
				downloading = true
				console.SuccessSameLine("\r%s\x1b[K", result.DefaultTemplate)
			} else {
				if downloading {
					downloading = false
					fmt.Fprintln(os.Stdout)
				}
				console.Info("%s", line)
			}
		} else {
			if downloading {
				downloading = false
				fmt.Fprintln(os.Stdout)
			}
			console.Info("%s", line)
		}
	}

	if err := waitForPipeClose(stdout); err != nil {
		console.Error("Error waiting for command: %v", err)
		return err
	}

	<-done

	fmt.Println()

	return nil
}

func waitForPipeClose(r io.Reader) error {
	if closer, ok := r.(io.Closer); ok {
		return closer.Close()
	}
	return nil
}
