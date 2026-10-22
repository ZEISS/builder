package console

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

var (
	checkMark = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
	errorMark = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).SetString("✗")
)

// FallbackWidth is used when stdout is not a terminal, so the banner still
// wraps to something sane when output is piped or captured.
const FallbackWidth = 80

// View is a view of the exit banner.
type View func(width int) string

// Render returns the exit banner for the given style.
func Render(width int, views ...View) string {
	if width <= 0 {
		width = FallbackWidth
	}

	style := lipgloss.NewStyle()
	contentWidth := width - style.GetHorizontalFrameSize()

	output := make([]string, len(views))
	for _, view := range views {
		output = append(output, view(contentWidth))
	}

	return style.Render(strings.Join(output, "\n"))
}

// Success returns the success section of the exit banner.
func Success(msg string) View {
	return func(width int) string {
		return fmt.Sprintf("%s %s", checkMark, msg)
	}
}

// Error returns the error section of the exit banner.
func Error(err error) View {
	return func(width int) string {
		return fmt.Sprintf("%s %s", errorMark, err.Error())
	}
}
