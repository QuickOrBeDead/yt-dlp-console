package ytdlp

import (
	"bytes"
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

func (f *YtdlpFakeExecutor) Execute(ctx context.Context, cmd *YtDlpCommandArgs, stdout *bytes.Buffer, stderr *bytes.Buffer) error {
	if f.ExecuteErr != nil {
		return f.ExecuteErr
	}
	if f.VideoDataResponse != nil {
		json.NewEncoder(stdout).Encode(f.VideoDataResponse)
	}
	return nil
}

func (f *YtdlpFakeExecutor) ExecuteWithStreams(ctx context.Context, cmd *YtDlpCommandArgs) (io.Reader, io.Reader, error) {
	if f.StreamsErr != nil {
		return nil, nil, f.StreamsErr
	}
	progress := strings.Join(f.ProgressLines, "\n")
	return strings.NewReader(progress), strings.NewReader(""), nil
}
