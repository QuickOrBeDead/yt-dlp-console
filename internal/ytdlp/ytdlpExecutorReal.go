package ytdlp

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"charm.land/huh/v2/spinner"
	"github.com/QuickOrBeDead/yt-dlp-console/internal/appconfig"
)

type YtdlpExecutorReal struct {
	config *appconfig.Config
}

func NewYtdlpExecutor(config *appconfig.Config) YtdlpExecutorReal {
	return YtdlpExecutorReal{config: config}
}

func (r YtdlpExecutorReal) Execute(ctx context.Context, cmd *YtDlpCommandArgs, cmdDesc string, stdout *bytes.Buffer, stderr *bytes.Buffer) error {
	execCmd := exec.CommandContext(ctx, r.config.YtDlpCommand, cmd.BuildArgs()...)
	execCmd.Stdout = stdout
	execCmd.Stderr = stderr

	fmt.Fprintf(os.Stdout, "yt-dlp command: %s %s\n", r.config.YtDlpCommand, strings.Join(cmd.BuildArgsMasked(), " "))

	err := spinner.New().
		Title(cmdDesc).
		ActionWithErr(func(ctx context.Context) error {
			err := execCmd.Run()
			return err
		}).
		Run()

	cmd.ClearPassword()

	return err
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

	fmt.Fprintf(os.Stdout, "yt-dlp command: %s %s\n", r.config.YtDlpCommand, strings.Join(cmd.BuildArgsMasked(), " "))

	err = execCmd.Start()

	cmd.ClearPassword()

	return stdout, stderr, err
}
