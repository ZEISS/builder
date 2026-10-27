// Package clipboard provides cross-platform clipboard access behind build-tag
// guards so that unsupported platforms (e.g., Android, iOS) compile without
// requiring CGO or platform-specific dependencies.
package clipboard

import (
	"errors"

	tea "charm.land/bubbletea/v2"

	"github.com/zeiss/builder/internal/ui/status"
)

// Format identifies the kind of data to read from the clipboard.
type Format int

const (
	// FormatText is plain UTF-8 text.
	FormatText Format = iota
	// FormatImage is binary image data (PNG).
	FormatImage
)

var (
	// ErrUnsupported is returned when clipboard access is not available on
	// the current platform.
	ErrUnsupported = errors.New("clipboard operations are not supported on this platform")
	// ErrEmpty is returned when the clipboard holds no data for the
	// requested format.
	ErrEmpty = errors.New("clipboard is empty or holds an unsupported format")
	// ErrWriteFailed is returned when the platform has a clipboard but the
	// written text is not there afterwards.
	ErrWriteFailed = errors.New("clipboard write did not land")
)

// Init initializes the clipboard subsystem. On unsupported platforms it
// returns ErrUnsupported but is otherwise safe to call.
func Init() error {
	return initClipboard()
}

// WriteText writes plain text to the system clipboard and reads it back to
// confirm the write landed. It returns ErrUnsupported on platforms without a
// clipboard or when Init did not succeed, and ErrWriteFailed when the
// clipboard does not hold the text afterwards.
func WriteText(text string) error {
	return writeText(text)
}

// Read returns the clipboard contents for the given format. It returns
// ErrEmpty when no matching data is present and ErrUnsupported on platforms
// without clipboard support.
func Read(f Format) ([]byte, error) {
	return read(f)
}

// CopyToClipboard copies the given text to the clipboard using both OSC 52
// (terminal escape sequence) and native clipboard for maximum compatibility.
// Returns a command that reports the outcome to the user, using the given
// message on success.
func CopyToClipboard(text, successMessage string) tea.Cmd {
	return CopyToClipboardWithCallback(text, successMessage, nil)
}

// CopyToClipboardWithCallback copies text to clipboard and executes a callback
// before showing the success message.
// This is useful when you need to perform additional actions like clearing UI state.
//
// The callback and the success message only run when the copy is believed to
// have worked, so callers can safely use the callback to discard the copied
// state (a selection, say) without losing it on a failed copy.
func CopyToClipboardWithCallback(text, successMessage string, callback tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		// The native write goes first and is verified before OSC 52 goes out,
		// because the terminal handles OSC 52 by writing the very same
		// clipboard on its own schedule. Verifying afterwards means reading
		// back in the middle of somebody else's write, which reports a good
		// copy as lost.
		err := WriteText(text)
		// OSC 52 is fire and forget: the terminal never answers, so a platform
		// without a native clipboard (an SSH session, say) gets the benefit of
		// the doubt. Only a native clipboard that accepted the write and then
		// does not hold the text is a real failure.
		osc52 := tea.SetClipboard(text)

		if errors.Is(err, ErrWriteFailed) {
			return tea.Sequence(osc52, status.Error(errors.New("failed to copy to clipboard")))()
		}

		return tea.Sequence(osc52, callback, status.Success(successMessage))()
	}
}
