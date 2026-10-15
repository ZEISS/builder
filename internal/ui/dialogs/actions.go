package dialogs

import tea "charm.land/bubbletea/v2"

// Action represents an action taken in a dialog after handling a message.
type Action any

// ActionClose is a message to close the current dialog.
type ActionClose struct{}

// ActionQuit is a message to quit the application.
type ActionQuit = tea.QuitMsg

// ActionOpenDialog is a message to open a dialog.
type ActionOpenDialog struct {
	DialogID string
}

// ActionCmd is a message to execute a command.
type ActionCmd struct {
	Command tea.Cmd
}

// ActionInitiateBuilder is sent when the OAuth builder is initiated.
type ActionInitiateBuilder struct{ URL string }

type (
	// ActionInitiateOAuth is sent when the device auth is initiated
	// successfully.
	ActionInitiateOAuth struct {
		DeviceCode      string
		UserCode        string
		ExpiresIn       int
		VerificationURL string
		Interval        int
	}

	// ActionCompleteOAuth is sent when the device flow completes successfully.
	ActionCompleteOAuth struct{}

	// ActionOAuthErrored is sent when the device flow encounters an error.
	ActionOAuthErrored struct {
		Error error
	}

	// ActionCloseOAuth closes the OAuth dialog and runs the given cleanup
	// command, cancelling any in-flight authorization. It exists so a
	// dismissed dialog does not leave a poller or loopback listener
	// running in the background.
	ActionCloseOAuth struct {
		Cmd tea.Cmd
	}
)
