// Package experimental is a banner to display when the experimental flag is enabled.
package experimental

import (
	"strings"

	"github.com/zeiss/builder/internal/ui/logo"
	"github.com/zeiss/builder/internal/ui/styles"
	"github.com/zeiss/builder/internal/versions"

	"charm.land/lipgloss/v2"
)

// FallbackWidth is used when stdout is not a terminal, so the banner still
// wraps to something sane when output is piped or captured.
const FallbackWidth = 80

// Render returns the exit banner for the given style.
func Render(version string, width int) string {
	if width <= 0 {
		width = FallbackWidth
	}

	style := lipgloss.NewStyle().Padding(1, 3)
	contentWidth := width - style.GetHorizontalFrameSize()

	sections := []string{
		logoSection(),
		experimentalSection(contentWidth),
	}

	return style.Render(strings.Join(sections, "\n\n"))
}

// experimentalSection returns the experimental section of the exit banner.
func experimentalSection(contentWidth int) string {
	return lipgloss.NewStyle().Width(contentWidth).Render("You're trying use an experimental feature.\n\nUse --experimental to enable the feature.")
}

// logoSection returns the ASCII art logo followed by the parting message.
func logoSection() string {
	t := styles.ThemeForProvider("")

	sparklesLogo := logo.Render(
		t.Logo.GradCanvas,
		versions.VersionString,
		true,
		logo.Opts{
			FieldColor:   t.Logo.FieldColor,
			TitleColorA:  t.Logo.TitleColorA,
			TitleColorB:  t.Logo.TitleColorB,
			CharmColor:   t.Logo.CharmColor,
			VersionColor: t.Logo.VersionColor,
		},
	)

	return sparklesLogo
}
