package dialogs

import (
	"cmp"
	"fmt"
	"image/color"
	"slices"
	"strings"
	"time"

	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/builder/internal/ui/list"
	"github.com/zeiss/builder/internal/ui/styles"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	uv "github.com/charmbracelet/ultraviolet"
	"github.com/charmbracelet/x/ansi"
	"github.com/rivo/uniseg"
	"github.com/sahilm/fuzzy"
)

// Dialog sizing constants.
const (
	// defaultDialogMaxWidth is the maximum width for standard dialogs.
	defaultDialogMaxWidth = 70
	// defaultDialogHeight is the default height for standard dialogs.
	defaultDialogHeight = 20
	// titleContentHeight is the height of the title content line.
	titleContentHeight = 1
	// inputContentHeight is the height of the input content line.
	inputContentHeight = 1
)

// Dialog is a component that can be displayed on top of the UI.
type Dialog interface {
	// ID returns the unique identifier of the dialog.
	ID() string
	// HandleMsg processes a message and returns an action. An [Action] can be
	// anything and the caller is responsible for handling it appropriately.
	HandleMsg(msg tea.Msg) Action
	// Draw draws the dialog onto the provided screen within the specified area
	// and returns the desired cursor position on the screen.
	Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor
}

// ListItem represents a selectable and searchable item in a dialog list.
type ListItem interface {
	list.FilterableItem
	list.Focusable
	list.MatchSettable

	// ID returns the unique identifier of the item.
	ID() string
}

// HoverDialog is a dialog that reacts to mouse hover. While one is open,
// the terminal must report all mouse motion events, not just clicks.
type HoverDialog interface {
	HandlesHover() bool
}

// LoadingDialog is a dialog that can show a loading state.
type LoadingDialog interface {
	StartLoading() tea.Cmd
	StopLoading()
}

// CloseKey is the default key binding to close dialogs.
var CloseKey = key.NewBinding(
	key.WithKeys("esc", "alt+esc"),
	key.WithHelp("esc", "exit"),
)

// Grace period constants for dialogs that open asynchronously and may
// receive in-flight keystrokes from a previously focused component.
const (
	// graceQuietPeriod is how long input must be quiet before the dialog
	// arms. Each absorbed keystroke resets this timer.
	graceQuietPeriod = 425 * time.Millisecond
	// graceMaxDelay is the absolute ceiling: the dialog always arms after
	// this duration regardless of input activity. Prevents auto-repeat
	// from keeping the dialog disarmed indefinitely.
	graceMaxDelay = 1500 * time.Millisecond
)

// Overlay manages multiple dialogs as an overlay.
type Overlay struct {
	dialogs []Dialog

	// Grace period state for the front dialog. Only active when the
	// dialog was opened via OpenDialogWithGrace.
	graceOpenedAt    time.Time
	graceLastInputAt time.Time

	// Track recently closed dialog IDs so that reopening the same
	// dialog type can skip the grace period. This prevents rapid
	// successive dialogs (e.g. multiple permission prompts) from
	// each eating a keystroke.
	lastClosedID string
	lastClosedAt time.Time
}

// NewOverlay creates a new [Overlay] instance.
func NewOverlay(dialogs ...Dialog) *Overlay {
	return &Overlay{
		dialogs: dialogs,
	}
}

// HasDialogs checks if there are any active dialogs.
func (d *Overlay) HasDialogs() bool {
	return len(d.dialogs) > 0
}

// HandlesHover reports whether any open dialog reacts to mouse hover, in
// which case the terminal should report all mouse motion events.
func (d *Overlay) HandlesHover() bool {
	for _, dialog := range d.dialogs {
		if hd, ok := dialog.(HoverDialog); ok && hd.HandlesHover() {
			return true
		}
	}
	return false
}

// ContainsDialog checks if a dialog with the specified ID exists.
func (d *Overlay) ContainsDialog(dialogID string) bool {
	for _, dialog := range d.dialogs {
		if dialog.ID() == dialogID {
			return true
		}
	}
	return false
}

// OpenDialog opens a new dialog to the stack.
func (d *Overlay) OpenDialog(dialog Dialog) {
	d.dialogs = append(d.dialogs, dialog)
	d.graceOpenedAt = time.Time{}
	d.graceLastInputAt = time.Time{}
}

// OpenDialogWithGrace opens a dialog with an input grace period. All
// keystrokes are absorbed until either the input has been quiet for
// graceQuietPeriod or graceMaxDelay has elapsed since opening, whichever
// comes first. Use this for dialogs that open asynchronously (e.g.
// permission prompts) where in-flight keystrokes from a previously
// focused component could act on the dialog before the user sees it.
//
// If the same dialog ID was just closed (within reopenGraceWindow),
// the grace period is skipped — the user is already focused on this
// dialog type and rapid successive prompts should not eat keystrokes.
func (d *Overlay) OpenDialogWithGrace(dialog Dialog) {
	now := time.Now()
	d.dialogs = append(d.dialogs, dialog)

	// Skip grace when reopening the same dialog type immediately.
	if dialog.ID() == d.lastClosedID && now.Sub(d.lastClosedAt) < reopenGraceWindow {
		d.graceOpenedAt = time.Time{}
		d.graceLastInputAt = time.Time{}
		return
	}

	d.graceOpenedAt = now
	d.graceLastInputAt = now
}

// inGracePeriod reports whether the front dialog is still within its
// input grace period. Returns false if no grace period is active.
func (d *Overlay) inGracePeriod() bool {
	if d.graceOpenedAt.IsZero() {
		return false
	}
	if time.Since(d.graceOpenedAt) >= graceMaxDelay {
		return false
	}
	if time.Since(d.graceLastInputAt) >= graceQuietPeriod {
		return false
	}
	return true
}

// CloseDialog closes the dialog with the specified ID from the stack.
func (d *Overlay) CloseDialog(dialogID string) {
	for i, dialog := range d.dialogs {
		if dialog.ID() == dialogID {
			d.removeDialog(i)
			return
		}
	}
}

// CloseFrontDialog closes the front dialog in the stack.
func (d *Overlay) CloseFrontDialog() {
	if len(d.dialogs) == 0 {
		return
	}
	d.removeDialog(len(d.dialogs) - 1)
}

// reopenGraceWindow is how long after closing a dialog we consider
// a reopen of the same dialog ID to be "immediate" and skip grace.
const reopenGraceWindow = 500 * time.Millisecond

func (d *Overlay) removeDialog(idx int) {
	if idx == len(d.dialogs)-1 {
		d.lastClosedID = d.dialogs[idx].ID()
		d.lastClosedAt = time.Now()
	}
	d.dialogs = append(d.dialogs[:idx], d.dialogs[idx+1:]...)
	// Clear grace state when the front dialog changes.
	if idx == len(d.dialogs) {
		d.graceOpenedAt = time.Time{}
		d.graceLastInputAt = time.Time{}
	}
}

// Dialog returns the dialog with the specified ID, or nil if not found.
func (d *Overlay) Dialog(dialogID string) Dialog {
	for _, dialog := range d.dialogs {
		if dialog.ID() == dialogID {
			return dialog
		}
	}
	return nil
}

// DialogLast returns the front dialog, or nil if there are no dialogs.
func (d *Overlay) DialogLast() Dialog {
	if len(d.dialogs) == 0 {
		return nil
	}
	return d.dialogs[len(d.dialogs)-1]
}

// BringToFront brings the dialog with the specified ID to the front.
func (d *Overlay) BringToFront(dialogID string) {
	for i, dialog := range d.dialogs {
		if dialog.ID() == dialogID {
			d.dialogs = append(d.dialogs[:i], d.dialogs[i+1:]...)
			d.dialogs = append(d.dialogs, dialog)
			break
		}
	}
}

// Update handles dialog updates.
func (d *Overlay) Update(msg tea.Msg) tea.Msg {
	if len(d.dialogs) == 0 {
		return nil
	}

	// Absorb keystrokes during the grace period for async dialogs.
	if _, ok := msg.(tea.KeyPressMsg); ok && d.inGracePeriod() {
		d.graceLastInputAt = time.Now()
		return nil
	}

	idx := len(d.dialogs) - 1 // active dialog is the last one
	dialog := d.dialogs[idx]
	if dialog == nil {
		return nil
	}

	return dialog.HandleMsg(msg)
}

// StartLoading starts the loading state for the front dialog if it
// implements [LoadingDialog].
func (d *Overlay) StartLoading() tea.Cmd {
	dialog := d.DialogLast()
	if ld, ok := dialog.(LoadingDialog); ok {
		return ld.StartLoading()
	}
	return nil
}

// StopLoading stops the loading state for the front dialog if it
// implements [LoadingDialog].
func (d *Overlay) StopLoading() {
	dialog := d.DialogLast()
	if ld, ok := dialog.(LoadingDialog); ok {
		ld.StopLoading()
	}
}

// Draw renders the overlay and its dialogs.
func (d *Overlay) Draw(scr uv.Screen, area uv.Rectangle) *tea.Cursor {
	var cur *tea.Cursor

	for _, dialog := range d.dialogs {
		cur = dialog.Draw(scr, area)
	}

	return cur
}

// ListItemStyles holds the styles for list items in a dialog.
type ListItemStyles struct {
	ItemBlurred     lipgloss.Style
	ItemFocused     lipgloss.Style
	InfoTextBlurred lipgloss.Style
	InfoTextFocused lipgloss.Style
}

func renderItem(t ListItemStyles, title, info string, focused bool, width int, cache map[int]string, m *fuzzy.Match) string {
	if cache == nil {
		cache = make(map[int]string)
	}

	cached, ok := cache[width]
	if ok {
		return cached
	}

	style := t.ItemBlurred
	if focused {
		style = t.ItemFocused
	}

	// Build content to the width left after the item style's own padding,
	// so the final rendered line is exactly `width` wide. Otherwise the
	// padding pushes the line past the list width and it wraps.
	lineWidth := max(0, width-style.GetHorizontalFrameSize())

	var infoText string
	var infoWidth int
	if len(info) > 0 {
		// Cap the info column so a long value (e.g. a provider name) can
		// truncate instead of overflowing the row or squeezing the title
		// to nothing; the title keeps at least half the width.
		if maxInfo := lineWidth / 2; lipgloss.Width(info)+2 > maxInfo {
			info = ansi.Truncate(info, max(0, maxInfo-2), "…")
		}
		infoText = fmt.Sprintf(" %s ", info)
		if focused {
			infoText = t.InfoTextFocused.Render(infoText)
		} else {
			infoText = t.InfoTextBlurred.Render(infoText)
		}

		infoWidth = lipgloss.Width(infoText)
	}

	title = ansi.Truncate(title, max(0, lineWidth-infoWidth), "…")
	titleWidth := lipgloss.Width(title)
	gap := strings.Repeat(" ", max(0, lineWidth-titleWidth-infoWidth))
	content := title
	if m != nil && len(m.MatchedIndexes) > 0 {
		var lastPos int
		parts := make([]string, 0)
		ranges := matchedRanges(m.MatchedIndexes)
		for _, rng := range ranges {
			start, stop := bytePosToVisibleCharPos(title, rng)
			if start > lastPos {
				parts = append(parts, ansi.Cut(title, lastPos, start))
			}
			// NOTE: We're using [ansi.Style] here instead of [lipglosStyle]
			// because we can control the underline start and stop more
			// precisely via [ansi.AttrUnderline] and [ansi.AttrNoUnderline]
			// which only affect the underline attribute without interfering
			// with other style attributes.
			parts = append(
				parts,
				ansi.NewStyle().Underline(true).String(),
				ansi.Cut(title, start, stop+1),
				ansi.NewStyle().Underline(false).String(),
			)
			lastPos = stop + 1
		}
		if lastPos < ansi.StringWidth(title) {
			parts = append(parts, ansi.Cut(title, lastPos, ansi.StringWidth(title)))
		}

		content = strings.Join(parts, "")
	}

	content = style.Render(content + gap + infoText)
	cache[width] = content
	return content
}

func matchedRanges(in []int) [][2]int {
	if len(in) == 0 {
		return [][2]int{}
	}
	current := [2]int{in[0], in[0]}
	if len(in) == 1 {
		return [][2]int{current}
	}
	var out [][2]int
	for i := 1; i < len(in); i++ {
		if in[i] == current[1]+1 {
			current[1] = in[i]
		} else {
			out = append(out, current)
			current = [2]int{in[i], in[i]}
		}
	}
	out = append(out, current)
	return out
}

func bytePosToVisibleCharPos(str string, rng [2]int) (int, int) {
	bytePos, byteStart, byteStop := 0, rng[0], rng[1]
	pos, start, stop := 0, 0, 0
	gr := uniseg.NewGraphemes(str)
	for byteStart > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Str())
		pos += max(1, gr.Width())
	}
	start = pos
	for byteStop > bytePos {
		if !gr.Next() {
			break
		}
		bytePos += len(gr.Str())
		pos += max(1, gr.Width())
	}
	stop = pos
	return start, stop
}

// sameFuzzyMatch returns true if the two fuzzy matches are equal.
func sameFuzzyMatch(a, b fuzzy.Match) bool {
	return a.Str == b.Str &&
		a.Index == b.Index &&
		a.Score == b.Score &&
		slices.Equal(a.MatchedIndexes, b.MatchedIndexes)
}

// DrawCenter draws the given string view centered in the screen area.
func DrawCenter(scr uv.Screen, area uv.Rectangle, view string) {
	DrawCenterCursor(scr, area, view, nil)
}

// DrawCenterCursor draws the given string view centered in the screen area and
// adjusts the cursor position accordingly. Content larger than the area is
// clamped to fit.
func DrawCenterCursor(scr uv.Screen, area uv.Rectangle, view string, cur *tea.Cursor) {
	width, height := lipgloss.Size(view)
	// Clamp to available area so oversized dialogs don't draw outside bounds.

	width = min(width, area.Dx())
	height = min(height, area.Dy())
	center := defaults.CenterRect(area, width, height)

	if cur != nil {
		cur.X += center.Min.X
		cur.Y += center.Min.Y
	}

	uv.NewStyledString(view).Draw(scr, center)
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

// RenderContext is a dialog rendering context that can be used to render
// common dialog layouts.
type RenderContext struct {
	// Styles is the styles to use for rendering.
	Styles *styles.Styles
	// TitleStyle is the style of the dialog title by default it uses Styles.Dialog.Title
	TitleStyle lipgloss.Style
	// ViewStyle is the style of the dialog title by default it uses Styles.Dialog.View
	ViewStyle lipgloss.Style
	// TitleGradientFromColor is the color the title gradient starts by default
	// its Styles.Dialog.TitleGradFromColor
	TitleGradientFromColor color.Color
	// TitleGradientToColor is the color the title gradient ends by default its
	// Styles.Dialog.TitleGradToColor
	TitleGradientToColor color.Color
	// Width is the total width of the dialog including any margins, borders,
	// and paddings.
	Width int
	// Gap is the gap between content parts. Zero means no gap.
	Gap int
	// Title is the title of the dialog. This will be styled using the default
	// dialog title style and prepended to the content parts slice.
	Title string
	// TitleInfo is additional information to display next to the title. This
	// part is displayed as is, any styling must be applied before setting this
	// field.
	TitleInfo string
	// Parts are the rendered parts of the dialog.
	Parts []string
	// Help is the fully rendered help footer line. Produce it with
	// renderDialogHelp so it is sized and padded consistently; it is
	// appended as-is without further styling.
	Help string
	// IsOnboarding indicates whether to render the dialog as part of the
	// onboarding flow. This means that the content will be rendered at the
	// bottom left of the screen.
	IsOnboarding bool
}

// NewRenderContext creates a new RenderContext with the provided styles and width.
func NewRenderContext(t *styles.Styles, width int) *RenderContext {
	return &RenderContext{
		Styles:                 t,
		TitleStyle:             t.Dialog.Title,
		ViewStyle:              t.Dialog.View,
		TitleGradientFromColor: t.Dialog.TitleGradFromColor,
		TitleGradientToColor:   t.Dialog.TitleGradToColor,
		Width:                  width,
		Parts:                  []string{},
	}
}

// AddPart adds a rendered part to the dialog.
func (rc *RenderContext) AddPart(part string) {
	if len(part) > 0 {
		rc.Parts = append(rc.Parts, part)
	}
}

// Render renders the dialog using the provided context.
func (rc *RenderContext) Render() string {
	titleStyle := rc.TitleStyle
	dialogStyle := rc.ViewStyle.Width(rc.Width)

	var parts []string

	if len(rc.Title) > 0 {
		contentWidth := rc.Width - dialogStyle.GetHorizontalFrameSize() -
			titleStyle.GetHorizontalFrameSize()
		titleInfo := rc.TitleInfo
		titleInfoWidth := lipgloss.Width(titleInfo)
		// Drop the title info entirely when it can't sit beside the title
		// text with at least a one-cell gap. Title info is often styled
		// (e.g. radio toggles with backgrounds and padding); truncating it
		// mid-segment leaves broken colored fragments, so hide it instead.
		if titleInfoWidth > 0 && lipgloss.Width(rc.Title)+1+titleInfoWidth > contentWidth {
			titleInfo = ""
			titleInfoWidth = 0
		}
		title := DialogTitle(rc.Styles, rc.Title,
			max(0, contentWidth-titleInfoWidth), rc.TitleGradientFromColor, rc.TitleGradientToColor)
		if len(titleInfo) > 0 {
			title += titleInfo
		}
		parts = append(parts, titleStyle.Render(title))
		if rc.Gap > 0 {
			parts = append(parts, make([]string, rc.Gap)...)
		}
	}

	if rc.Gap <= 0 {
		parts = append(parts, rc.Parts...)
	} else {
		for i, p := range rc.Parts {
			if len(p) > 0 {
				parts = append(parts, p)
			}
			if i < len(rc.Parts)-1 {
				parts = append(parts, make([]string, rc.Gap)...)
			}
		}
	}

	if len(rc.Help) > 0 {
		if rc.Gap > 0 {
			parts = append(parts, make([]string, rc.Gap)...)
		}
		parts = append(parts, rc.Help)
	}

	content := strings.Join(parts, "\n")
	if rc.IsOnboarding {
		return content
	}
	return dialogStyle.Render(content)
}

// InputCursor adjusts the cursor position for an input field within a dialog.
func InputCursor(t *styles.Styles, cur *tea.Cursor) *tea.Cursor {
	if cur != nil {
		titleStyle := t.Dialog.Title
		dialogStyle := t.Dialog.View
		inputStyle := t.Dialog.InputPrompt
		// Adjust cursor position to account for dialog layout
		cur.X += inputStyle.GetBorderLeftSize() +
			inputStyle.GetMarginLeft() +
			inputStyle.GetPaddingLeft() +
			dialogStyle.GetBorderLeftSize() +
			dialogStyle.GetPaddingLeft() +
			dialogStyle.GetMarginLeft()
		cur.Y += titleStyle.GetVerticalFrameSize() +
			inputStyle.GetBorderTopSize() +
			inputStyle.GetMarginTop() +
			inputStyle.GetPaddingTop() +
			inputStyle.GetBorderBottomSize() +
			inputStyle.GetMarginBottom() +
			inputStyle.GetPaddingBottom() +
			dialogStyle.GetPaddingTop() +
			dialogStyle.GetMarginTop() +
			dialogStyle.GetBorderTopSize()
	}
	return cur
}

func dialogInputTextWidth(t *styles.Styles, input textinput.Model, contentWidth int) int {
	const cursorPadding = 1
	return max(0, contentWidth-
		t.Dialog.InputPrompt.GetHorizontalFrameSize()-
		lipgloss.Width(input.Prompt)-
		cursorPadding)
}

// renderDialogHelp renders keybind hints as a single padded footer line at
// contentWidth (the dialog's inner width: total minus the View border). The
// hints are packed greedily and truncated with an ellipsis so the line never
// wraps or overflows the border, and never ends on a dangling separator.
func renderDialogHelp(t *styles.Styles, h *help.Model, km help.KeyMap, contentWidth int) string {
	textWidth := max(0, contentWidth-t.Dialog.HelpView.GetHorizontalFrameSize())
	return t.Dialog.HelpView.Render(shortHelpLine(h, km.ShortHelp(), textWidth))
}

// shortHelpLine builds a single-line short help view truncated to width.
// It reimplements the bubbles help packing to avoid a component bug where
// items are kept even when they overflow (when the ellipsis itself does not
// fit), and to guarantee the line ends cleanly rather than on a separator.
func shortHelpLine(h *help.Model, bindings []key.Binding, width int) string {
	if width <= 0 {
		return ""
	}
	sep := h.Styles.ShortSeparator.Inline(true).Render(h.ShortSeparator)
	ellipsis := h.Styles.Ellipsis.Inline(true).Render(cmp.Or(h.Ellipsis, "…"))

	var b strings.Builder
	total := 0
	for _, kb := range bindings {
		if !kb.Enabled() {
			continue
		}
		seg := ""
		if total > 0 {
			seg = sep
		}
		seg += h.Styles.ShortKey.Inline(true).Render(kb.Help().Key) + " " +
			h.Styles.ShortDesc.Inline(true).Render(kb.Help().Desc)
		w := lipgloss.Width(seg)
		if total+w > width {
			// The next item doesn't fit; add an ellipsis if there's room.
			// The separator belongs to this dropped item, so what we've
			// written already ends on a real hint, not a dangling dot. A
			// leading space joins the ellipsis to prior hints, but only
			// when there are prior hints.
			tail := ellipsis
			if total > 0 {
				tail = " " + ellipsis
			}
			if total+lipgloss.Width(tail) <= width {
				b.WriteString(tail)
			}
			break
		}
		total += w
		b.WriteString(seg)
	}
	return b.String()
}

// adjustOnboardingInputCursor removes the dialog view frame offset from an
// input cursor. Onboarding dialogs render without Dialog.View frame, while
// InputCursor includes that frame offset for regular dialogs.
func adjustOnboardingInputCursor(t *styles.Styles, cur *tea.Cursor) *tea.Cursor {
	if cur == nil {
		return nil
	}

	dialogStyle := t.Dialog.View
	cur.X -= dialogStyle.GetBorderLeftSize() +
		dialogStyle.GetPaddingLeft() +
		dialogStyle.GetMarginLeft()
	cur.Y -= dialogStyle.GetBorderTopSize() +
		dialogStyle.GetPaddingTop() +
		dialogStyle.GetMarginTop()
	return cur
}

// DrawOnboarding draws the given string view centered in the screen area.
func DrawOnboarding(scr uv.Screen, area uv.Rectangle, view string) {
	DrawOnboardingCursor(scr, area, view, nil)
}

// DrawOnboardingCursor draws the given string view positioned at the bottom
// left area of the screen. Content larger than the area is clamped to fit.
func DrawOnboardingCursor(scr uv.Screen, area uv.Rectangle, view string, cur *tea.Cursor) {
	width, height := lipgloss.Size(view)
	// Clamp to available area so oversized dialogs don't draw outside bounds.
	width = min(width, area.Dx())
	height = min(height, area.Dy())
	bottomLeft := defaults.BottomLeftRect(area, width, height)
	if cur != nil {
		cur.X += bottomLeft.Min.X
		cur.Y += bottomLeft.Min.Y
	}
	uv.NewStyledString(view).Draw(scr, bottomLeft)
}
