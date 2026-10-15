// Package logo renders a Builder in a stylized way.
package logo

import (
	"fmt"
	"image/color"
	"math/rand/v2"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/ui/styles"
)

// letterform represents a letterform. It can be stretched horizontally by
// a given amount via the boolean argument.
type letterform func(bool) string

const diag = `/`

// Opts are the options for rendering the Sparkles title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	CharmColor   color.Color // Charm™ text color
	VersionColor color.Color // version text color
	Width        int         // width of the rendered logo, used for truncation

	// When true, stretch a random letterform on each render. Has no effect in
	// compact mode. Mainly for testing. In production you will want to cache
	// the stretched letterform to keep the logo from jittering on resize.
	Unstable bool
}

// Render renders the Sparkles logo. Set the argument to true to render the narrow
// version, intended for use in a sidebar.
//
// The compact argument determines whether it renders compact for the sidebar
// or wider for the main pane.
func Render(base lipgloss.Style, version string, compact bool, o Opts) string {
	builder := "Builder"

	fg := func(c color.Color, s string) string {
		return lipgloss.NewStyle().Foreground(c).Render(s)
	}

	// Title.
	const spacing = 1

	sparklesLetterforms := []letterform{
		LetterSAlt,
		LetterP,
		LetterA,
		LetterR,
		LetterK,
		LetterL,
		LetterE,
		LetterSAlt,
	}

	stretchIndex := -1 // -1 means no stretching.

	if !compact && !o.Unstable {
		// Always stretch the same letterform, which is picked once at random.
		stretchIndex = cachedRandN(len(sparklesLetterforms))
	} else if !compact && o.Unstable {
		// Stretch a random letterform on every render.
		stretchIndex = rand.IntN(len(sparklesLetterforms)) //nolint:gosec
	}

	sparkles := renderWord(spacing, stretchIndex, sparklesLetterforms...)

	sparklesWidth := lipgloss.Width(sparkles)
	b := new(strings.Builder)
	for r := range strings.SplitSeq(sparkles, "\n") {
		fmt.Fprintln(b, styles.ApplyForegroundGrad(base, r, o.TitleColorA, o.TitleColorB))
	}

	sparkles = b.String()

	metaRowGap := 1
	maxVersionWidth := sparklesWidth - lipgloss.Width(builder) - metaRowGap
	version = ansi.Truncate(version, maxVersionWidth, "…") // truncate version if too long.
	gap := max(0, sparklesWidth-lipgloss.Width(builder)-lipgloss.Width(version))
	metaRow := fg(o.CharmColor, builder) + strings.Repeat(" ", gap) + fg(o.VersionColor, version)

	sparkles = strings.TrimSpace(metaRow + "\n" + sparkles)

	// Narrow version.
	if compact {
		field := fg(o.FieldColor, strings.Repeat(diag, sparklesWidth))
		return strings.Join([]string{field, field, sparkles, field, ""}, "\n")
	}

	fieldHeight := lipgloss.Height(sparkles)

	// Left field.
	const leftWidth = 6
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field.
	rightWidth := max(15, o.Width-sparklesWidth-leftWidth-2) // 2 for the gap.
	const stepDownAt = 0
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		width := rightWidth
		if i >= stepDownAt {
			width = rightWidth - (i - stepDownAt)
		}
		fmt.Fprint(rightField, fg(o.FieldColor, strings.Repeat(diag, width)), "\n")
	}

	// Return the wide version.
	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, sparkles, hGap, rightField.String())
	if o.Width > 0 {
		// Truncate the logo to the specified width.
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}

	return logo
}
