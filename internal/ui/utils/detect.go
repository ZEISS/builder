package utils

import (
	"os"
	"strings"

	xstrings "github.com/charmbracelet/x/exp/strings"
	"github.com/charmbracelet/x/term"
)

// SupportsProgressBar tries to determine whether the current terminal supports
// progress bars by looking into environment variables.
func SupportsProgressBar() bool {
	if !term.IsTerminal(os.Stderr.Fd()) {
		return false
	}

	termProg := os.Getenv("TERM_PROGRAM")
	_, isWindowsTerminal := os.LookupEnv("WT_SESSION")

	return isWindowsTerminal || xstrings.ContainsAnyOf(strings.ToLower(termProg), "ghostty", "iterm2", "rio", "zed")
}
