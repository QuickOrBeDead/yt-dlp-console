package ytdlp

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

type YtDlpCommand struct {
	url      string
	password string
	args     []string
}

func NewYtDlpCommand(url, password string) *YtDlpCommand {
	return &YtDlpCommand{
		url:      url,
		password: password,
		args:     make([]string, 0),
	}
}

func (c *YtDlpCommand) AddArg(arg string) {
	c.args = append(c.args, arg)
}

func (c *YtDlpCommand) AddArgWithValue(arg, value string) {
	c.args = append(c.args, arg)
	c.args = append(c.args, value)
}

func (c *YtDlpCommand) Execute() *exec.Cmd {
	cfg := appconfig.Get()

	fmt.Printf("yt-dlp command: %s %s\n", cfg.YtDlpCommand, strings.Join(c.buildArgs(true), " "))

	return exec.Command(cfg.YtDlpCommand, c.buildArgs(false)...)
}

func (c *YtDlpCommand) buildArgs(hidePassword bool) []string {
	args := make([]string, 0)

	if len(c.password) > 0 {
		args = append(args, "--video-password")

		if hidePassword {
			args = append(args, "******")
		} else {
			args = append(args, c.password)
		}
	}

	args = append(args, c.args...)
	args = append(args, c.url)

	return args
}

func (c *YtDlpCommand) ClearPassword() {
	if len(c.password) > 0 {
		c.password = ""
	}
}
