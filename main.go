package main

import (
	"runtime"

	"github.com/QuickOrBeDead/yt-dlp-console/cmd"
	"github.com/muesli/termenv"
)

func main() {
	if runtime.GOOS == "windows" {
		// In the most recent versions, calling ColorProfile()
		// triggers the Windows internal ANSI initialization.
		_ = termenv.DefaultOutput().ColorProfile()
	}

	cmd.Execute()
}
