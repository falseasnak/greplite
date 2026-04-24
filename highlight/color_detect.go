package highlight

import (
	"os"
	"runtime"
	"strings"
)

// ShouldUseColor returns true if color output is appropriate for the
// given file descriptor. It checks for a NO_COLOR env var, the TERM
// variable, and whether the fd is a terminal on supported platforms.
func ShouldUseColor(f *os.File) bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	term := os.Getenv("TERM")
	if term == "dumb" || term == "" {
		return false
	}
	// On Windows, color support depends on the terminal emulator.
	if runtime.GOOS == "windows" {
		return os.Getenv("WT_SESSION") != "" ||
			strings.Contains(os.Getenv("TERM_PROGRAM"), "vscode")
	}
	return isTerminal(f)
}

// isTerminal returns true if f is connected to a TTY.
func isTerminal(f *os.File) bool {
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	// ModeCharDevice is set for TTYs on Unix.
	return (fi.Mode() & os.ModeCharDevice) != 0
}
