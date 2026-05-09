package cmd

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/ytdlp"
)

func captureOutput(f func()) string {
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = oldStdout

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

// FakeFormProvider returns pre-programmed responses for each form interaction.
// Calls are consumed sequentially from results.
type FakeFormProvider struct {
	t       *testing.T
	results []fakeFormResult
	pos     int
}

type fakeFormResult struct {
	value string
	err   error
}

func (f *FakeFormProvider) Input(title string, validate func(string) error) (string, error) {
	r := f.next()
	if validate != nil && r.err == nil {
		r.err = validate(r.value)
	}
	return r.value, r.err
}

func (f *FakeFormProvider) InputPassword(title string) (string, error) {
	r := f.next()
	return r.value, r.err
}

func (f *FakeFormProvider) Select(title string, options []string) (string, error) {
	r := f.next()
	return r.value, r.err
}

func (f *FakeFormProvider) Confirm(title, description string) (bool, error) {
	r := f.next()
	if r.err != nil {
		return false, r.err
	}
	return r.value == "true", nil
}

func (f *FakeFormProvider) next() fakeFormResult {
	if f.pos >= len(f.results) {
		f.t.Fatal("unexpected FormProvider call")
	}
	r := f.results[f.pos]
	f.pos++
	return r
}

func TestIsValidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{"valid https", "https://youtube.com/watch?v=test", true},
		{"valid http", "http://example.com/video", true},
		{"valid with params", "https://example.com/video?id=123&t=30s", true},
		{"valid with auth", "https://user:pass@example.com/video", true},
		{"empty", "", false},
		{"no scheme", "not-a-url", false},
		{"only scheme", "https://", false},
		{"relative path", "/relative/path", false},
		{"spaces only", "   ", false},
		{"ftp with host", "ftp://files.example.com/video", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidURL(tt.url); got != tt.want {
				t.Errorf("isValidURL(%q) = %v, want %v", tt.url, got, tt.want)
			}
		})
	}
}

func TestIndexOf(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		value string
		want  int
	}{
		{"first element", []string{"a", "b", "c"}, "a", 0},
		{"middle element", []string{"a", "b", "c"}, "b", 1},
		{"last element", []string{"a", "b", "c"}, "c", 2},
		{"not found", []string{"a", "b", "c"}, "d", 0},
		{"empty slice", []string{}, "a", 0},
		{"duplicates", []string{"a", "b", "a"}, "a", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := indexOf(tt.slice, tt.value); got != tt.want {
				t.Errorf("indexOf(%v, %q) = %d, want %d", tt.slice, tt.value, got, tt.want)
			}
		})
	}
}

func TestRunDownloadFlow(t *testing.T) {
	hasAudioFmt := ytdlp.VideoFormat{
		FormatID: "22", VCodec: "avc1", ACodec: "mp4a", Ext: "mp4", Height: 720,
	}
	videoOnlyFmt := ytdlp.VideoFormat{
		FormatID: "137", VCodec: "avc1", ACodec: "none", Ext: "mp4", Height: 1080,
	}
	audioOnlyFmt := ytdlp.VideoFormat{
		FormatID: "140", VCodec: "none", ACodec: "mp4a", Ext: "m4a",
	}

	hasAudioData := &ytdlp.VideoData{Title: "Test", Formats: []ytdlp.VideoFormat{hasAudioFmt}}
	_, hasAudioLabels := hasAudioData.GetVideoList()

	videoOnlyWithAudioData := &ytdlp.VideoData{
		Title:   "Test",
		Formats: []ytdlp.VideoFormat{videoOnlyFmt, audioOnlyFmt},
	}
	_, videoLabels := videoOnlyWithAudioData.GetVideoList()
	_, audioLabels := videoOnlyWithAudioData.GetAudioList()

	videoOnlyNoAudioData := &ytdlp.VideoData{
		Title:   "Test",
		Formats: []ytdlp.VideoFormat{videoOnlyFmt},
	}
	_, videoOnlyLabels := videoOnlyNoAudioData.GetVideoList()

	noFormatsData := &ytdlp.VideoData{Title: "Test", Formats: []ytdlp.VideoFormat{}}

	tests := []struct {
		name    string
		forms   *FakeFormProvider
		exec    *ytdlp.YtdlpFakeExecutor
		wantErr string
	}{
		{
			name: "success no auth video has audio",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
				{value: hasAudioLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: hasAudioData,
				ProgressLines:     []string{`{"_default_template": "100%"}`},
			},
		},
		{
			name: "success video only with audio selection",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
				{value: videoLabels[0]},
				{value: audioLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: videoOnlyWithAudioData,
				ProgressLines:     []string{`{"_default_template": "100%"}`},
			},
		},
		{
			name: "success video only no audio available",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
				{value: videoOnlyLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: videoOnlyNoAudioData,
				ProgressLines:     []string{`{"_default_template": "100%"}`},
			},
		},
		{
			name: "success password auth",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "Password"},
				{value: "secret123"},
				{value: hasAudioLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: hasAudioData,
				ProgressLines:     []string{`{"_default_template": "100%"}`},
			},
		},
		{
			name: "success username password auth",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "Username + Password"},
				{value: "myuser"},
				{value: "mypass"},
				{value: hasAudioLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: hasAudioData,
				ProgressLines:     []string{`{"_default_template": "100%"}`},
			},
		},
		{
			name: "get video data error",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				ExecuteErr: fmt.Errorf("yt-dlp not found"),
			},
			wantErr: "yt-dlp not found",
		},
		{
			name: "no video formats found",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: noFormatsData,
			},
			wantErr: "no video found",
		},
		{
			name: "download error",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{value: "https://youtube.com/watch?v=test"},
				{value: "None"},
				{value: hasAudioLabels[0]},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{
				VideoDataResponse: hasAudioData,
				StreamsErr:        fmt.Errorf("download failed"),
			},
			wantErr: "download failed",
		},
		{
			name: "form input error",
			forms: &FakeFormProvider{t: t, results: []fakeFormResult{
				{err: fmt.Errorf("user cancelled")},
			}},
			exec: &ytdlp.YtdlpFakeExecutor{},
			wantErr: "user cancelled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := ytdlp.NewYtDlpClient(tt.exec, &appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})
			err := runDownloadFlow(context.Background(), client, tt.forms)

			if tt.wantErr != "" {
				if err == nil {
					t.Errorf("expected error %q, got nil", tt.wantErr)
				} else if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("expected error containing %q, got %q", tt.wantErr, err.Error())
				}
			} else if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRootCmd_Success(t *testing.T) {
	origForms := defaultForms
	defer func() { defaultForms = origForms }()
	origExecutor := newYtdlpExecutor
	defer func() { newYtdlpExecutor = origExecutor }()

	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	appconfig.Save(&appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})

	hasAudioFmt := ytdlp.VideoFormat{
		FormatID: "22", VCodec: "avc1", ACodec: "mp4a", Ext: "mp4", Height: 720,
	}
	hasAudioData := &ytdlp.VideoData{Title: "Test", Formats: []ytdlp.VideoFormat{hasAudioFmt}}
	_, labels := hasAudioData.GetVideoList()

	newYtdlpExecutor = func(config *appconfig.Config) ytdlp.YtDlpExecutor {
		return &ytdlp.YtdlpFakeExecutor{
			VideoDataResponse: hasAudioData,
			ProgressLines:     []string{`{"_default_template": "100%"}`},
		}
	}

	defaultForms = &FakeFormProvider{t: t, results: []fakeFormResult{
		{value: "https://youtube.com/watch?v=test"},
		{value: "None"},
		{value: labels[0]},
	}}

	output := captureOutput(func() {
		rootCmd.SetArgs([]string{})
		rootCmd.Execute()
	})

	if !strings.Contains(output, "Download complete!") {
		t.Errorf("expected success message, got:\n%s", output)
	}
}

func TestRootCmd_FormError(t *testing.T) {
	origForms := defaultForms
	defer func() { defaultForms = origForms }()
	origExecutor := newYtdlpExecutor
	defer func() { newYtdlpExecutor = origExecutor }()

	tmpDir := t.TempDir()
	os.Setenv("XDG_CONFIG_HOME", tmpDir)
	defer os.Unsetenv("XDG_CONFIG_HOME")

	appconfig.Save(&appconfig.Config{YtDlpCommand: "yt-dlp", N: 5})

	newYtdlpExecutor = func(config *appconfig.Config) ytdlp.YtDlpExecutor {
		return &ytdlp.YtdlpFakeExecutor{}
	}

	defaultForms = &FakeFormProvider{t: t, results: []fakeFormResult{
		{err: fmt.Errorf("user cancelled")},
	}}

	output := captureOutput(func() {
		rootCmd.SetArgs([]string{})
		rootCmd.Execute()
	})

	if !strings.Contains(output, "user cancelled") {
		t.Errorf("expected error message, got:\n%s", output)
	}
}
