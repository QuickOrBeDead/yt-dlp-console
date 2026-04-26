package ytdlp

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type VideoData struct {
	Title    string        `json:"title"`
	Duration uint64        `json:"duration"`
	Formats  []VideoFormat `json:"formats"`
}

func (v *VideoData) GetAudioList() ([]VideoFormat, []string) {
	formats := make([]VideoFormat, 0)
	labels := make([]string, 0)

	for _, f := range v.Formats {
		if f.HasAudio() && !f.HasVideo() {
			formats = append(formats, f)
			labels = append(labels, f.getDescription(v))
		}
	}

	return formats, labels
}

func (v *VideoData) GetVideoList() ([]VideoFormat, []string) {
	formats := make([]VideoFormat, 0)
	labels := make([]string, 0)

	for _, f := range v.Formats {
		if f.HasVideo() {
			formats = append(formats, f)
			labels = append(labels, f.getDescription(v))
		}
	}

	return formats, labels
}

type VideoFormat struct {
	FormatID        string  `json:"format_id"`
	FormatNote      string  `json:"format_note,omitempty"`
	Ext             string  `json:"ext"`
	Resolution      string  `json:"resolution,omitempty"`
	FileSize        uint64  `json:"filesize,omitempty"`
	FileSizeApprox  uint64  `json:"filesize_approx,omitempty"`
	VCodec          string  `json:"vcodec,omitempty"`
	ACodec          string  `json:"acodec,omitempty"`
	AudioExt        string  `json:"audio_ext,omitempty"`
	VideoExt        string  `json:"video_ext,omitempty"`
	Width           int     `json:"width"`
	Height          int     `json:"height"`
	Language        string  `json:"language"`
	TotalBitRate    float64 `json:"tbr,omitempty"`
	AudioSampleRate int     `json:"asr,omitempty"`
	AudioBitrate    int     `json:"abr,omitempty"`
	VideoBitrate    int     `json:"vbr,omitempty"`
}

func (f *VideoFormat) getSize(v *VideoData) string {
	if f.FileSize != 0 {
		return formatBytes(f.FileSize)
	}

	if f.FileSizeApprox != 0 {
		return fmt.Sprintf("~%s", formatBytes(f.FileSizeApprox))
	}

	size := fileSizeFromTotalBitRate(f.TotalBitRate, v.Duration)
	if size != 0 {
		return fmt.Sprintf("~%s", formatBytes(size))
	}

	return ""
}

func (f *VideoFormat) getAudioCodec() string {
	if len(f.ACodec) == 0 {
		return "unknown"
	}

	if f.ACodec == "none" {
		return ""
	}

	codecMap := map[string]string{
		"mp4a.40.2": "aac",
		"mp4a.40.5": "aac",
		"mp4a.40":   "aac",
		"opus":      "opus",
		"vorbis":    "vorbis",
		"mp3":       "mp3",
		"aac":       "aac",
		"flac":      "flac",
		"wav":       "wav",
	}

	if mapped, ok := codecMap[f.ACodec]; ok {
		return mapped
	}

	return f.ACodec
}

func (f *VideoFormat) HasAudio() bool {
	return f.ACodec != "none"
}

func (f *VideoFormat) HasVideo() bool {
	return f.VCodec != "none"
}

func (f *VideoFormat) getVideoCodec() string {
	if len(f.VCodec) == 0 {
		return "unknown"
	}

	if f.VCodec == "none" {
		return ""
	}

	codecMap := map[string]string{
		"avc1.640015":               "h264",
		"avc1.64001f":               "h264",
		"avc1.64001e":               "h264",
		"avc1.64000d":               "h264",
		"avc1.42001e":               "h264",
		"avc1.42001f":               "h264",
		"avc1.420020":               "h264",
		"avc1.42001a":               "h264",
		"avc1.640032":               "h265",
		"hev1.1.6.L120.90":          "h265",
		"hev1.1.6.L93":              "h265",
		"vp9":                       "vp9",
		"vp09.00.10.08.01.01.01.00": "vp9",
		"av01":                      "av1",
		"av01.0.05M.08":             "av1",
		"av01.0.08M.10":             "av1",
		"h263":                      "h263",
		"mpeg4":                     "mpeg4",
		"theora":                    "theora",
	}

	if mapped, ok := codecMap[f.VCodec]; ok {
		return mapped
	}

	return f.VCodec
}

func (f *VideoFormat) getDescription(v *VideoData) string {
	size := f.getSize(v)

	var sb strings.Builder

	var t string
	if f.HasAudio() && f.HasVideo() {
		t = "Video + Audio"
	} else if f.HasAudio() {
		t = "Audio Only"
	} else if f.HasVideo() {
		t = "Video Only"
	}

	if len(t) > 0 {
		fmt.Fprintf(&sb, "%s - ", t)
	}

	if f.HasVideo() {
		fmt.Fprintf(&sb, "%s %dp", f.Ext, f.Height)
	} else if f.HasAudio() && f.AudioSampleRate > 0 {
		fmt.Fprintf(&sb, "%s %.1fk", f.Ext, float64(f.AudioSampleRate)/1000)
	} else {
		fmt.Fprintf(&sb, "%s", f.Ext)
	}

	if len(size) > 0 {
		fmt.Fprintf(&sb, " - %s", f.getSize(v))
	}

	videoCodec := f.getVideoCodec()
	audioCodec := f.getAudioCodec()

	properties := make([]string, 0)
	if f.HasAudio() && len(audioCodec) > 0 {
		properties = append(properties, fmt.Sprintf("Audio Codec: %s", audioCodec))
	}
	if f.HasVideo() && len(videoCodec) > 0 {
		properties = append(properties, fmt.Sprintf("Video Codec: %s", videoCodec))
	}
	if len(f.Language) > 0 {
		properties = append(properties, fmt.Sprintf("Language: %s", f.Language))
	}
	if len(f.FormatNote) > 0 {
		properties = append(properties, fmt.Sprintf("Note: %s", f.FormatNote))
	}

	if len(properties) > 0 {
		fmt.Fprintf(&sb, " (%s)", strings.Join(properties, ", "))
	}

	return sb.String()
}

func fileSizeFromTotalBitRate(totalBitRate float64, duration uint64) uint64 {
	if totalBitRate == 0 || duration == 0 {
		return 0
	}

	return uint64(float64(duration) * totalBitRate * (1000 / 8))
}

func formatBytes(bytes uint64) string {
	suffixes := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB", "ZiB", "YiB"}
	format := "%s%s"
	base := 1024.0

	fBytes := float64(bytes)
	if fBytes < base {
		return fmt.Sprintf(format, strconv.Itoa(int(bytes)), suffixes[0])
	}

	index := math.Floor(math.Log(fBytes) / math.Log(base))
	volume := base * fBytes / math.Pow(base, index+1)

	return fmt.Sprintf(format, strconv.FormatFloat(volume, 'f', 2, 64), suffixes[int(index)])
}
