package main

import (
	"os"
	"runtime"

	"github.com/QuickOrBeDead/yt-dlp-console/cmd"
	"golang.org/x/sys/windows"
)

func main() {
	if runtime.GOOS == "windows" {
		enableVirtualTerminal()
	}

	cmd.Execute()
}

func enableVirtualTerminal() {
	stdout := windows.Handle(os.Stdout.Fd())
	var mode uint32
	windows.GetConsoleMode(stdout, &mode)
	windows.SetConsoleMode(stdout, mode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
