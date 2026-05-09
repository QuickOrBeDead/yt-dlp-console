package ytdlp

import (
	"strings"
	"testing"
)

func TestYtDlpCommand_BuildArgs(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		password   string
		setup      func(*YtDlpCommandArgs)
		wantArgs   []string
		wantMasked string
	}{
		{
			name:     "simple url",
			url:      "https://youtube.com/watch?v=123",
			password: "",
			setup:    nil,
			wantArgs: []string{"https://youtube.com/watch?v=123"},
		},
		{
			name:     "with password",
			url:      "https://youtube.com/watch?v=123",
			password: "secret",
			setup:    nil,
			wantArgs: []string{"--video-password", "secret", "https://youtube.com/watch?v=123"},
		},
		{
			name:     "with additional args",
			url:      "https://youtube.com/watch?v=123",
			password: "",
			setup: func(c *YtDlpCommandArgs) {
				c.AddArg("-J")
				c.AddArgWithValue("-f", "22")
			},
			wantArgs: []string{"-J", "-f", "22", "https://youtube.com/watch?v=123"},
		},
		{
			name:       "masked password",
			url:        "https://youtube.com/watch?v=123",
			password:   "secret",
			setup:      nil,
			wantMasked: "--video-password ****** https://youtube.com/watch?v=123",
		},
		{
			name:     "with username and account password",
			url:      "https://youtube.com/watch?v=123",
			password: "",
			setup: func(c *YtDlpCommandArgs) {
				c.SetAccountAuth("myuser", "mypass")
			},
			wantArgs: []string{"--username", "myuser", "--password", "mypass", "https://youtube.com/watch?v=123"},
		},
		{
			name:     "with username, account password, and video password",
			url:      "https://youtube.com/watch?v=123",
			password: "videosecret",
			setup: func(c *YtDlpCommandArgs) {
				c.SetAccountAuth("myuser", "mypass")
			},
			wantArgs: []string{"--username", "myuser", "--password", "mypass", "--video-password", "videosecret", "https://youtube.com/watch?v=123"},
		},
		{
			name:       "masked username and account password",
			url:        "https://youtube.com/watch?v=123",
			password:   "",
			setup:      func(c *YtDlpCommandArgs) { c.SetAccountAuth("myuser", "mypass") },
			wantMasked: "--username myuser --password ****** https://youtube.com/watch?v=123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewYtDlpCommandArgs(tt.url, tt.password)
			if tt.setup != nil {
				tt.setup(cmd)
			}

			args := cmd.BuildArgs()
			if len(tt.wantArgs) > 0 {
				if len(args) != len(tt.wantArgs) {
					t.Errorf("BuildArgs() length = %d, want %d", len(args), len(tt.wantArgs))
				}
				for i, arg := range args {
					if i < len(tt.wantArgs) && arg != tt.wantArgs[i] {
						t.Errorf("BuildArgs()[%d] = %q, want %q", i, arg, tt.wantArgs[i])
					}
				}
			}

			if tt.wantMasked != "" {
				masked := strings.Join(cmd.BuildArgsMasked(), " ")
				if masked != tt.wantMasked {
					t.Errorf("BuildArgsMasked() = %q, want %q", masked, tt.wantMasked)
				}
			}
		})
	}
}

func TestYtDlpCommand_AddArg(t *testing.T) {
	cmd := NewYtDlpCommandArgs("https://example.com", "")
	cmd.AddArg("-J")
	cmd.AddArg("--no-playlist")

	args := cmd.BuildArgs()
	if len(args) != 3 { // 2 args + url
		t.Errorf("expected 3 args, got %d", len(args))
	}
	if args[0] != "-J" || args[1] != "--no-playlist" {
		t.Errorf("unexpected args: %v", args)
	}
}

func TestYtDlpCommand_AddArgWithValue(t *testing.T) {
	cmd := NewYtDlpCommandArgs("https://example.com", "")
	cmd.AddArgWithValue("-f", "22")
	cmd.AddArgWithValue("-N", "4")

	args := cmd.BuildArgs()
	if len(args) != 5 { // 2 pairs + url
		t.Errorf("expected 5 args, got %d", len(args))
	}
	if args[0] != "-f" || args[1] != "22" || args[2] != "-N" || args[3] != "4" {
		t.Errorf("unexpected args: %v", args)
	}
}
