package output

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/sqve/kamaji/internal/config"
)

// MessageType represents different output message categories.
type MessageType int

const (
	Success MessageType = iota
	Error
	Info
	Warning
	Debug
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	debugStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

var prefixes = map[MessageType]string{
	Success: "[ok] ",
	Error:   "Error: ",
	Info:    "-> ",
	Warning: "Warning: ",
	Debug:   "[DEBUG] ",
}

// Prefix returns the prefix for a message type.
func Prefix(t MessageType) string {
	p := prefixes[t]
	if config.IsPlain() {
		return p
	}
	return styleFor(t).Render(p)
}

// Style returns the styled string for the given type and message.
func Style(t MessageType, msg string) string {
	prefix := Prefix(t)
	return prefix + msg
}

func styleFor(t MessageType) lipgloss.Style {
	switch t {
	case Success:
		return successStyle
	case Error:
		return errorStyle
	case Info:
		return infoStyle
	case Warning:
		return warningStyle
	case Debug:
		return debugStyle
	default:
		return lipgloss.NewStyle()
	}
}

// PrintSuccess writes a success message to stdout.
func PrintSuccess(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, Style(Success, msg))
}

// PrintError writes an error message to stderr.
func PrintError(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, Style(Error, msg))
}

// PrintInfo writes an info message to stdout.
func PrintInfo(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, Style(Info, msg))
}

// PrintWarning writes a warning message to stdout.
func PrintWarning(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, Style(Warning, msg))
}

// PrintDebug writes a debug message to stdout.
func PrintDebug(msg string) {
	_, _ = fmt.Fprintln(os.Stdout, Style(Debug, msg))
}

// SuccessMsg returns a styled success message.
func SuccessMsg(msg string) string {
	return Style(Success, msg)
}

// ErrorMsg returns a styled error message.
func ErrorMsg(msg string) string {
	return Style(Error, msg)
}

// InfoMsg returns a styled info message.
func InfoMsg(msg string) string {
	return Style(Info, msg)
}

// WarningMsg returns a styled warning message.
func WarningMsg(msg string) string {
	return Style(Warning, msg)
}

// DebugMsg returns a styled debug message.
func DebugMsg(msg string) string {
	return Style(Debug, msg)
}
