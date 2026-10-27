package models

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/builder/internal/ui/logo"
	"github.com/zeiss/builder/internal/versions"
)

const (
	headerDiag           = "╱"
	minHeaderDiags       = 3
	leftPadding          = 1
	rightPadding         = 1
	diagToDetailsSpacing = 1 // space between diagonal pattern and details section
)

// Header is the model for the header.
type Header struct {
	logo        string
	compactLogo string

	ui      *defaults.UI
	width   int
	compact bool
}

// creates a new header model.
func NewHeader(ui *defaults.UI) *Header {
	h := &Header{
		ui: ui,
	}

	return h
}

// Draw draws the header for the given session.
func (h *Header) Draw(
	scr uv.Screen,
	area uv.Rectangle,
	compact bool,
	detailsOpen bool,
	width int,
) {
	t := h.ui.Styles

	h.logo = logo.Render(
		t.Logo.GradCanvas,
		versions.VersionString,
		compact,
		logo.Opts{
			FieldColor:   t.Logo.FieldColor,
			TitleColorA:  t.Logo.TitleColorA,
			TitleColorB:  t.Logo.TitleColorB,
			CharmColor:   t.Logo.CharmColor,
			VersionColor: t.Logo.VersionColor,
			Width:        width,
		},
	)

	h.width = width
	h.compact = compact

	var b strings.Builder
	b.WriteString(h.logo)

	view := uv.NewStyledString(
		t.Header.Wrapper.Padding(0, rightPadding, 0, leftPadding).Render(b.String()),
	)

	view.Draw(scr, area)
}
