package ytdlp

import (
	"context"
	"encoding/json"
	"io"
	"strings"
)

type YtdlpFakeExecutor struct {
	VideoDataResponse *VideoData
	ExecuteErr        error
	ProgressLines     []string
	StreamsErr        error
}

func (f *YtdlpFakeExecutor) Execute(ctx context.Context, cmd *YtDlpCommandArgs) ([]byte, error) {
	if f.ExecuteErr != nil {
		return nil, f.ExecuteErr
	}
	return json.Marshal(f.VideoDataResponse)
}

func (f *YtdlpFakeExecutor) ExecuteWithStreams(ctx context.Context, cmd *YtDlpCommandArgs) (io.Reader, io.Reader, error) {
	if f.StreamsErr != nil {
		return nil, nil, f.StreamsErr
	}
	progress := strings.Join(f.ProgressLines, "\n")
	return strings.NewReader(progress), strings.NewReader(""), nil
}
