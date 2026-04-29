package ytdlp

import (
	"context"
	"io"
)

type YtDlpExecutor interface {
	Execute(ctx context.Context, cmd *YtDlpCommandArgs) ([]byte, error)
	ExecuteWithStreams(ctx context.Context, cmd *YtDlpCommandArgs) (io.Reader, io.Reader, error)
}
