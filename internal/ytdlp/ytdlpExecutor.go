package ytdlp

import (
	"bytes"
	"context"
	"io"
)

type YtDlpExecutor interface {
	Execute(ctx context.Context, cmd *YtDlpCommandArgs, cmdDesc string, stdout *bytes.Buffer, stderr *bytes.Buffer) error
	ExecuteWithStreams(ctx context.Context, cmd *YtDlpCommandArgs) (io.Reader, io.Reader, error)
}
