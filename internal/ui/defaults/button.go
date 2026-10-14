package defaults

import (
	"fmt"

	"charm.land/lipgloss/v2"
	"github.com/zeiss/builder/internal/ui/styles"
)

// ButtonOpts defines the configuration for a single button
type ButtonOpts struct {
	// Text is the button label
	Text string
	// UnderlineIndex is the 0-based index of the character to underline (-1 for none)
	UnderlineIndex int
	// Selected indicates whether this button is currently selected
	Selected bool
	// Hovered indicates whether the mouse is hovering over the button
	Hovered bool
	// Inactive marks a button inside a prompt that is not in the active
	// pane. It renders lighter than the plain blurred style so the
	// choices stay legible while the chat has focus.
	Inactive bool
	// Padding inner horizontal padding defaults to 2 if this is 0
	Padding int
}

// Button creates a button with an underlined character and selection state
func Button(t *styles.Styles, opts ButtonOpts) string {
	style := t.Button.Blurred

	if opts.Selected && opts.Hovered {
		style = t.Button.Focused.Bold(true)
	} else if opts.Hovered {
		style = t.Button.Hovered.Bold(true)
	} else if opts.Selected {
		style = t.Button.Focused
	} else if opts.Inactive {
		style = t.Button.Inactive
	}

	text := opts.Text
	if opts.Padding == 0 {
		opts.Padding = 2
	}

	// the index is out of bound
	if opts.UnderlineIndex > -1 && opts.UnderlineIndex > len(text)-1 {
		opts.UnderlineIndex = -1
	}

	text = style.Padding(0, opts.Padding).Render(text)

	if opts.UnderlineIndex != -1 {
		text = lipgloss.StyleRanges(text, lipgloss.NewRange(opts.Padding+opts.UnderlineIndex, opts.Padding+opts.UnderlineIndex+1, style.Underline(true)))
	}

	return text
}

// HitButtonIndex checks a compositor for a button hit and returns
// the button index, or -1 if no button was hit.
func HitButtonIndex(c *lipgloss.Compositor, x, y int) int {
	if c == nil {
		return -1
	}

	hit := c.Hit(x, y)
	if hit.Empty() {
		return -1
	}

	var idx int
	if _, err := fmt.Sscanf(hit.ID(), "btn_%d", &idx); err != nil {
		return -1
	}
	return idx
}
