package defaults

import (
	"image"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/zeiss/builder/internal/ui/styles"
)

// CenterRect returns a new [Rectangle] centered within the given area with the
// specified width and height.
func CenterRect(area uv.Rectangle, width, height int) uv.Rectangle {
	centerX := area.Min.X + area.Dx()/2
	centerY := area.Min.Y + area.Dy()/2
	minX := centerX - width/2
	minY := centerY - height/2
	maxX := minX + width
	maxY := minY + height

	return image.Rect(minX, minY, maxX, maxY)
}

// BottomLeftRect returns a new [Rectangle] positioned at the bottom-left within the given area with the
// specified width and height.
func BottomLeftRect(area uv.Rectangle, width, height int) uv.Rectangle {
	minX := area.Min.X
	maxX := minX + width
	maxY := area.Max.Y
	minY := maxY - height
	return image.Rect(minX, minY, maxX, maxY)
}

// DialogTitle renders a dialog title with a decorative line filling the
// remaining width. When the title alone exceeds the available width it is
// truncated with an ellipsis so it never wraps.
func DialogTitle(t *styles.Styles, title string, width int, fromColor, toColor color.Color) string {
	if width > 0 && lipgloss.Width(title) > width {
		return ansi.Truncate(title, width, "…")
	}

	char := "╱"
	length := lipgloss.Width(title) + 1
	remainingWidth := width - length

	if remainingWidth > 0 {
		lines := strings.Repeat(char, remainingWidth)
		lines = styles.ApplyForegroundGrad(t.Dialog.TitleLineBase, lines, fromColor, toColor)
		title = title + " " + lines
	}

	return title
}
