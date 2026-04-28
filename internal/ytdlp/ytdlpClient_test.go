package ytdlp

import (
	"context"
	"fmt"
	"testing"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

func TestYtDlpClient_GetVideoData(t *testing.T) {
	tests := []struct {
		name      string
		fake      *YtdlpFakeExecutor
		url       string
		password  string
		wantTitle string
		wantErr   bool
	}{
		{
			name: "successful video data retrieval",
			fake: &YtdlpFakeExecutor{
				VideoDataResponse: &VideoData{
					Title:   "Test Video",
					Formats: []VideoFormat{{FormatID: "22"}},
				},
			},
			url:       "https://youtube.com/watch?v=test",
			password:  "",
			wantTitle: "Test Video",
			wantErr:   false,
		},
		{
			name: "with password",
			fake: &YtdlpFakeExecutor{
				VideoDataResponse: &VideoData{
					Title:   "Private Video",
					Formats: []VideoFormat{{FormatID: "18"}},
				},
			},
			url:       "https://youtube.com/watch?v=private",
			password:  "secret",
			wantTitle: "Private Video",
			wantErr:   false,
		},
		{
			name: "executor returns error",
			fake: &YtdlpFakeExecutor{
				ExecuteErr: fmt.Errorf("yt-dlp not found"),
			},
			url:     "https://youtube.com/watch?v=test",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewYtDlpClient(tt.fake, &appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})
			ctx := context.Background()

			data, err := client.GetVideoData(ctx, tt.url, tt.password)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if !tt.wantErr && data != nil && data.Title != tt.wantTitle {
				t.Errorf("expected title %q, got %q", tt.wantTitle, data.Title)
			}
		})
	}
}

func TestYtDlpClient_DownloadVideo(t *testing.T) {
	tests := []struct {
		name    string
		fake    *YtdlpFakeExecutor
		url     string
		format  string
		wantErr bool
	}{
		{
			name: "successful download",
			fake: &YtdlpFakeExecutor{
				ProgressLines: []string{`{"_default_template": "Downloading..."}`},
			},
			url:     "https://youtube.com/watch?v=test",
			format:  "22",
			wantErr: false,
		},
		{
			name: "download with error",
			fake: &YtdlpFakeExecutor{
				StreamsErr: fmt.Errorf("download failed"),
			},
			url:     "https://youtube.com/watch?v=test",
			format:  "22",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewYtDlpClient(tt.fake, &appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})
			ctx := context.Background()

			err := client.DownloadVideo(ctx, tt.url, "", tt.format)

			if tt.wantErr && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
