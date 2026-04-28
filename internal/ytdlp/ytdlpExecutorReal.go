package ytdlp

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

type YtdlpExecutorReal struct {
	config *appconfig.Config
}

func NewYtdlpExecutor(config *appconfig.Config) YtdlpExecutorReal {
	return YtdlpExecutorReal{config: config}
}

func (r YtdlpExecutorReal) Execute(ctx context.Context, cmd *YtDlpCommandArgs) ([]byte, error) {
	fmt.Printf("yt-dlp command: %s %s\n", r.config.YtDlpCommand, strings.Join(cmd.BuildArgsMasked(), " "))
	execCmd := exec.CommandContext(ctx, r.config.YtDlpCommand, cmd.BuildArgs()...)

	cmd.ClearPassword()

	return execCmd.Output()
}

func (r YtdlpExecutorReal) ExecuteWithStreams(ctx context.Context, cmd *YtDlpCommandArgs) (io.Reader, io.Reader, error) {
	execCmd := exec.CommandContext(ctx, r.config.YtDlpCommand, cmd.BuildArgs()...)

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	fmt.Printf("yt-dlp command: %s %s\n", r.config.YtDlpCommand, strings.Join(cmd.BuildArgsMasked(), " "))

	err = execCmd.Start()

	cmd.ClearPassword()

	return stdout, stderr, err
}
