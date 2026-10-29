package status

import (
	"time"

	tea "charm.land/bubbletea/v2"
)

// InfoMsgType is the type of an info message.
type InfoMsgType int

const (
	InfoMsgTypeInfo InfoMsgType = iota
	InfoMsgTypeSuccess
	InfoMsgTypeWarn
	InfoMsgTypeError
	InfoMsgTypeUpdate
	InfoMsgTypePlan
	InfoMsgTypeYolo
)

// InfoMsg is a message to be displayed in the status bar.
type InfoMsg struct {
	// Type is the type of the message.
	Type InfoMsgType
	// Msg is the message text.
	Msg string
	// TTL is the time-to-live for the message.
	TTL time.Duration
}

// ClearStatusMsg is a message to clear the status bar.
type ClearStatusMsg struct{}

// IsEmpty checks if the [InfoMsg] is empty.
func (m InfoMsg) IsEmpty() bool {
	var zero InfoMsg
	return m == zero
}

// DefaultStatusTTL is the default time-to-live for status messages.
const DefaultStatusTTL = 5 * time.Second

// ClearInfoMsgCmd returns a command that clears the info
// message after the given TTL.
func ClearInfoMsgCmd(ttl time.Duration) tea.Cmd {
	return tea.Tick(ttl, func(time.Time) tea.Msg {
		return ClearStatusMsg{}
	})
}

// Info returns a command that displays an info message.
func Info(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeInfo,
			Msg:  msg,
		}
	}
}

// Error returns a command that displays an error message.
func Error(err error) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeError,
			Msg:  err.Error(),
		}
	}
}

// Warn returns a command that displays a warning message.
func Warn(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeWarn,
			Msg:  msg,
		}
	}
}

// Success returns a command that displays a success message.
func Success(msg string) tea.Cmd {
	return func() tea.Msg {
		return InfoMsg{
			Type: InfoMsgTypeSuccess,
			Msg:  msg,
		}
	}
}

// Msg is a helper function that returns a [tea.Msg] from a [tea.Cmd].
func Msg(cmd tea.Cmd) tea.Msg {
	return cmd()
}

// Cmd returns a command that displays a message.
func Cmd(msg InfoMsg) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}
