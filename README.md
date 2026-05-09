# yt-dlp-console

![GitHub release](https://img.shields.io/github/v/release/QuickOrBeDead/yt-dlp-console)
![Tests](https://img.shields.io/github/actions/workflow/status/QuickOrBeDead/yt-dlp-console/test.yml?label=tests)
[![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/QuickOrBeDead/yt-dlp-console/badges/coverage.json)](https://github.com/QuickOrBeDead/yt-dlp-console/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/QuickOrBeDead/yt-dlp-console)](https://goreportcard.com/report/github.com/QuickOrBeDead/yt-dlp-console)

A user-friendly CLI tool that makes downloading videos with [yt-dlp](https://github.com/yt-dlp/yt-dlp) easy through an interactive terminal interface.

## Why this tool?

yt-dlp is powerful but requires remembering complex command-line options. This interactive tool guides you through the download process step-by-step in your terminal—no flags or format codes needed.

## Features
- Support for password-protected videos
- Lists and lets you select from available video/audio formats
- Downloads with configurable concurrent fragment downloads
- Real-time download progress display
- Persistent configuration

## Requirements

- [Go](https://go.dev/) 1.26.2+
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) installed and in PATH

## Installation

```bash
go install github.com/QuickOrBeDead/yt-dlp-console@latest
```

Or clone and build:

```bash
git clone https://github.com/QuickOrBeDead/yt-dlp-console.git
cd yt-dlp-console
go build -o yt-dlp-console .
```

## Usage

### Download a video

Run without arguments to start an interactive download flow:

```bash
yt-dlp-console
```

You'll be prompted for:
1. Video URL
2. Authentication (None, Password, or Username + Password)
3. Video format selection
4. Audio format selection (if video doesn't include audio)

### Configure settings

```bash
yt-dlp-console config
```

This lets you set:
- yt-dlp command path (default: `yt-dlp`)
- Concurrent fragments `-N` (1-32, default: 4)

## Configuration

Settings are stored in:
- Linux/macOS: `~/.config/yt-dlp-console/config.json`
- Windows: `%APPDATA%\yt-dlp-console\config.json`

## License

MIT License - see [LICENSE](LICENSE) file.