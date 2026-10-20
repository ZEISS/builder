// Package exitbanner is
package exitbanner

import (
	"math/rand/v2"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ui/logo"
	"github.com/zeiss/builder/internal/ui/styles"
	"github.com/zeiss/builder/internal/versions"
)

// FallbackWidth is used when stdout is not a terminal, so the banner still
// wraps to something sane when output is piped or captured.
const FallbackWidth = 80

// Render returns the exit banner for the given style.
func Render(session *models.Session, version string, width int) string {
	if width <= 0 {
		width = FallbackWidth
	}

	style := lipgloss.NewStyle().Padding(1, 3)
	contentWidth := width - style.GetHorizontalFrameSize()

	sections := []string{
		logoSection(contentWidth),
		sessionSection(session, contentWidth),
	}
	return style.Render(strings.Join(sections, "\n\n"))
}

// sessionSection returns the session section of the exit banner.
func sessionSection(session *models.Session, contentWidth int) string {
	return lipgloss.NewStyle().Width(contentWidth).Render("Session: " + session.ID.String())
}

// logoSection returns the ASCII art logo followed by the parting message.
func logoSection(contentWidth int) string {
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

	return sparklesLogo + "\n" +
		lipgloss.NewStyle().Width(contentWidth).Render("It's sparkling! "+randomExitMessage())
}

// randomExitMessage returns a random exit message.
func randomExitMessage() string {
	messages := []string{
		"",
		"Catch you on the flip side.",
		"Stay shiny.",
		"Until next time.",
		"Go stretch your legs.",
		"Don’t be a stranger.",
		"You crushed it.",
		"Treat yourself to a cookie.",
		"Coffee break?",
		"That was delightful.",
		"See you tomorrow?",
		"Go touch some grass.",
		"Who’s up for tacos?",
		"Be kind to yourself.",
		"Drink some water.",
		"Off you go, superstar.",
		"You’re simply dazzling.",
		"Well played.",
		"You’re a legend.",
		"Where did I leave my sunglasses?",
		"It’s cocoa o’clock.",
	}
	return messages[rand.IntN(len(messages))]
}
