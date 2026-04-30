package console

import (
	"fmt"

	"charm.land/lipgloss/v2"
)

var (
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true)
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Bold(true)
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)
	mutedStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Faint(true)
)

func Error(format string, args ...any) {
	fmt.Println(errorStyle.Render(fmt.Sprintf(format, args...)))
}

func Success(format string, args ...any) {
	fmt.Println(successStyle.Render(fmt.Sprintf(format, args...)))
}

func Warning(format string, args ...any) {
	fmt.Println(warningStyle.Render(fmt.Sprintf(format, args...)))
}

func Info(format string, args ...any) {
	fmt.Println(infoStyle.Render(fmt.Sprintf(format, args...)))
}

func Muted(format string, args ...any) {
	fmt.Println(mutedStyle.Render(fmt.Sprintf(format, args...)))
}

func SuccessSameLine(format string, args ...any) {
	fmt.Print(successStyle.Render(fmt.Sprintf(format, args...)))
}
