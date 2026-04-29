package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

func TestYtDlpClient_GetVideoData(t *testing.T) {
	tests := []struct {
		name       string
		fake       *YtdlpFakeExecutor
		url        string
		password   string
		wantTitle  string
		wantErr    bool
		wantStdOut string
		wantStdErr string
	}{
		{
			name: "successful video data retrieval",
			fake: &YtdlpFakeExecutor{
				VideoDataResponse: &VideoData{
					Title:   "Test Video",
					Formats: []VideoFormat{{FormatID: "22"}},
				},
			},
			url:        "https://youtube.com/watch?v=test",
			password:   "",
			wantTitle:  "Test Video",
			wantErr:    false,
			wantStdOut: "",
			wantStdErr: "",
		},
		{
			name: "with password",
			fake: &YtdlpFakeExecutor{
				VideoDataResponse: &VideoData{
					Title:   "Private Video",
					Formats: []VideoFormat{{FormatID: "18"}},
				},
			},
			url:        "https://youtube.com/watch?v=private",
			password:   "secret",
			wantTitle:  "Private Video",
			wantErr:    false,
			wantStdOut: "",
			wantStdErr: "",
		},
		{
			name: "executor returns error",
			fake: &YtdlpFakeExecutor{
				ExecuteErr: fmt.Errorf("yt-dlp not found"),
			},
			url:        "https://youtube.com/watch?v=test",
			wantErr:    true,
			wantStdOut: "",
			wantStdErr: "Error running yt-dlp: yt-dlp not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStderr := os.Stderr
			oldStdout := os.Stdout
			rErr, wErr, _ := os.Pipe()
			rOut, wOut, _ := os.Pipe()
			os.Stderr = wErr
			os.Stdout = wOut

			client := NewYtDlpClient(tt.fake, &appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})
			ctx := context.Background()

			data, err := client.GetVideoData(ctx, tt.url, tt.password)

			wErr.Close()
			wOut.Close()
			os.Stderr = oldStderr
			os.Stdout = oldStdout
			capturedErr, _ := io.ReadAll(rErr)
			capturedOut, _ := io.ReadAll(rOut)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && data != nil && data.Title != tt.wantTitle {
				t.Errorf("expected title %q, got %q", tt.wantTitle, data.Title)
			}
			if string(capturedErr) != tt.wantStdErr {
				t.Errorf("stderr: got %q, want %q", string(capturedErr), tt.wantStdErr)
			}
			if string(capturedOut) != tt.wantStdOut {
				t.Errorf("stdout: got %q, want %q", string(capturedOut), tt.wantStdOut)
			}
		})
	}
}

func TestYtDlpClient_DownloadVideo(t *testing.T) {
	tests := []struct {
		name       string
		fake       *YtdlpFakeExecutor
		url        string
		format     string
		wantErr    bool
		wantStdOut string
		wantStdErr string
	}{
		{
			name: "successful download",
			fake: &YtdlpFakeExecutor{
				ProgressLines: []string{
					"Preparing for download",
					`{"_default_template": "Downloading 50%..."}`,
					`{"_default_template": 1}`,
					"Warning",
					`{"_default_template": "Downloading 99%..."}`,
					"Finished!",
				},
			},
			url:        "https://youtube.com/watch?v=test",
			format:     "22",
			wantErr:    false,
			wantStdOut: "Preparing for download\n\rDownloading 50%...\x1b[K\n{\"_default_template\": 1}\nWarning\n\rDownloading 99%...\x1b[K\nFinished!\n\n",
			wantStdErr: "",
		},
		{
			name: "download with error",
			fake: &YtdlpFakeExecutor{
				StreamsErr: fmt.Errorf("download failed"),
			},
			url:        "https://youtube.com/watch?v=test",
			format:     "22",
			wantErr:    true,
			wantStdOut: "",
			wantStdErr: "Error running yt-dlp: download failed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStderr := os.Stderr
			oldStdout := os.Stdout
			rErr, wErr, _ := os.Pipe()
			rOut, wOut, _ := os.Pipe()
			os.Stderr = wErr
			os.Stdout = wOut

			client := NewYtDlpClient(tt.fake, &appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})
			ctx := context.Background()

			err := client.DownloadVideo(ctx, tt.url, "", tt.format)

			wErr.Close()
			wOut.Close()
			os.Stderr = oldStderr
			os.Stdout = oldStdout
			capturedErr, _ := io.ReadAll(rErr)
			capturedOut, _ := io.ReadAll(rOut)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if string(capturedErr) != tt.wantStdErr {
				t.Errorf("stderr: got %q, want %q", string(capturedErr), tt.wantStdErr)
			}
			if string(capturedOut) != tt.wantStdOut {
				t.Errorf("stdout: got %q, want %q", string(capturedOut), tt.wantStdOut)
			}
		})
	}
}
