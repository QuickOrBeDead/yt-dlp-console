package ytdlp

import (
	"strings"
	"testing"
)

func TestVideoFormat_HasAudio(t *testing.T) {
	tests := []struct {
		name   string
		codec  string
		expect bool
	}{
		{"has audio", "aac", true},
		{"no audio", "none", false},
		{"empty codec", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := VideoFormat{ACodec: tt.codec}
			if got := f.HasAudio(); got != tt.expect {
				t.Errorf("HasAudio() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestVideoFormat_HasVideo(t *testing.T) {
	tests := []struct {
		name   string
		codec  string
		expect bool
	}{
		{"has video", "h264", true},
		{"no video", "none", false},
		{"empty codec", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := VideoFormat{VCodec: tt.codec}
			if got := f.HasVideo(); got != tt.expect {
				t.Errorf("HasVideo() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestVideoFormat_getAudioCodec(t *testing.T) {
	tests := []struct {
		name   string
		codec  string
		expect string
	}{
		{"mp4a.40.2", "mp4a.40.2", "aac"},
		{"mp4a.40.5", "mp4a.40.5", "aac"},
		{"opus", "opus", "opus"},
		{"vorbis", "vorbis", "vorbis"},
		{"mp3", "mp3", "mp3"},
		{"aac", "aac", "aac"},
		{"flac", "flac", "flac"},
		{"wav", "wav", "wav"},
		{"none", "none", ""},
		{"empty", "", "unknown"},
		{"unknown codec", "unknown_codec", "unknown_codec"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := VideoFormat{ACodec: tt.codec}
			if got := f.getAudioCodec(); got != tt.expect {
				t.Errorf("getAudioCodec() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestVideoFormat_getVideoCodec(t *testing.T) {
	tests := []struct {
		name   string
		codec  string
		expect string
	}{
		{"avc1.640015", "avc1.640015", "h264"},
		{"avc1.64001f", "avc1.64001f", "h264"},
		{"avc1.640032", "avc1.640032", "h265"},
		{"vp9", "vp9", "vp9"},
		{"av01", "av01", "av1"},
		{"none", "none", ""},
		{"empty", "", "unknown"},
		{"unknown codec", "unknown_codec", "unknown_codec"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := VideoFormat{VCodec: tt.codec}
			if got := f.getVideoCodec(); got != tt.expect {
				t.Errorf("getVideoCodec() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name   string
		bytes  uint64
		expect string
	}{
		{"zero", 0, "0B"},
		{"bytes", 512, "512B"},
		{"kibibytes", 1024, "1.00KiB"},
		{"mibibytes", 1048576, "1.00MiB"},
		{"gibibytes", 1073741824, "1.00GiB"},
		{"large", 1536, "1.50KiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formatBytes(tt.bytes); got != tt.expect {
				t.Errorf("formatBytes() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestFileSizeFromTotalBitRate(t *testing.T) {
	tests := []struct {
		name     string
		bitrate  float64
		duration uint64
		expect   uint64
	}{
		{"normal", 1000, 60, uint64(60 * 1000 * 125)},
		{"zero bitrate", 0, 60, 0},
		{"zero duration", 1000, 0, 0},
		{"both zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileSizeFromTotalBitRate(tt.bitrate, tt.duration); got != tt.expect {
				t.Errorf("fileSizeFromTotalBitRate() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestVideoFormat_getSize(t *testing.T) {
	video := &VideoData{Duration: 60}

	tests := []struct {
		name   string
		format VideoFormat
		expect string
	}{
		{"filesize set", VideoFormat{FileSize: 1048576}, "1.00MiB"},
		{"filesize approx", VideoFormat{FileSizeApprox: 2097152}, "~2.00MiB"},
		{"from bitrate", VideoFormat{TotalBitRate: 1000}, "~7.15MiB"},
		{"no size info", VideoFormat{}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.format.getSize(video); got != tt.expect {
				t.Errorf("getSize() = %v, want %v", got, tt.expect)
			}
		})
	}
}

func TestVideoData_GetAudioList(t *testing.T) {
	video := &VideoData{
		Formats: []VideoFormat{
			{ACodec: "aac", VCodec: "none", FormatID: "140"},
			{ACodec: "opus", VCodec: "none", FormatID: "251"},
			{ACodec: "aac", VCodec: "h264", FormatID: "22"},
		},
	}

	formats, labels := video.GetAudioList()

	if len(formats) != 2 {
		t.Errorf("expected 2 audio formats, got %d", len(formats))
	}

	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestVideoData_GetVideoList(t *testing.T) {
	video := &VideoData{
		Formats: []VideoFormat{
			{ACodec: "aac", VCodec: "h264", FormatID: "22"},
			{ACodec: "none", VCodec: "vp9", FormatID: "244"},
			{ACodec: "aac", VCodec: "none", FormatID: "140"},
		},
	}

	formats, labels := video.GetVideoList()

	if len(formats) != 2 {
		t.Errorf("expected 2 video formats, got %d", len(formats))
	}

	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestVideoFormat_getDescription(t *testing.T) {
	video := &VideoData{Duration: 60}

	tests := []struct {
		name   string
		format VideoFormat
		check  string
	}{
		{
			"video + audio",
			VideoFormat{Ext: "mp4", Height: 720, VCodec: "h264", ACodec: "aac", FileSize: 1048576},
			"Video + Audio",
		},
		{
			"audio only",
			VideoFormat{Ext: "m4a", ACodec: "aac", VCodec: "none", AudioSampleRate: new(44100), FileSize: 524288},
			"Audio Only",
		},
		{
			"video only",
			VideoFormat{Ext: "webm", Height: 1080, VCodec: "vp9", ACodec: "none"},
			"Video Only",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			desc := tt.format.getDescription(video)
			if !strings.Contains(desc, tt.check) {
				t.Errorf("getDescription() = %v, should contain %v", desc, tt.check)
			}
		})
	}
}
