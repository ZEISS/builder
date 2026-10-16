package dialogs

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/oidc"
)

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
	Cmd tea.Cmd
}

// ActionInitiateBuilder is sent when the OAuth builder is initiated.
type ActionInitiateBuilder struct{ URL string }

type (
	// ActionOAuthInitialized is sent when the OAuth flow is initialized.
	ActionOAuthInitialized struct {
		DeviceAuth *oidc.DeviceAuth
	}

	// ActionOAuthCompleted is sent when the OAuth flow is completed successfully.
	ActionOAuthCompleted struct {
		Token *oidc.Token
	}

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
