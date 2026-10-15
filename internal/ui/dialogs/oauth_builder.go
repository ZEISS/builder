package dialogs

import (
	"context"
	"fmt"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/zeiss/builder/internal/adapters/oidc"
	"github.com/zeiss/builder/internal/controllers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/internal/ui/defaults"
	"github.com/zeiss/builder/server/middlewares/discovery"
)

// NewOAuthBuilder returns a new OAuth dialog for the given provider.
func NewOAuthBuilder(
	ui *defaults.UI,
	isOnboarding bool,
	oAuthProvider OAuthProvider,
) (*OAuth, tea.Cmd) {
	return newOAuth(ui, isOnboarding, oAuthProvider)
}

type OAuthBuilder struct {
	URL        string
	deviceAuth *models.DeviceAuth
	authCtrl   ports.DeviceAuthController
	cancelFunc func()
}

var _ OAuthProvider = (*OAuthBuilder)(nil)

func (m *OAuthBuilder) name() string {
	return "Builder"
}

func (m *OAuthBuilder) initiateAuth() tea.Msg {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	discv := discovery.NewClient()
	wellknownConfig, err := discv.Discover(ctx, m.URL)
	if err != nil {
		return ActionOAuthErrored{Error: fmt.Errorf("failed to initiate device auth: %w", err)}
	}

	oidcProvider := oidc.New(m.URL, "builder-cli", oidc.WithWellKnownConfig(wellknownConfig))
	authCtrl := controllers.NewDeviceAuthController(oidcProvider, nil)
	m.authCtrl = authCtrl

	deviceAuth, err := authCtrl.Begin(ctx)
	if err != nil {
		return ActionOAuthErrored{Error: fmt.Errorf("failed to initiate device auth: %w", err)}
	}

	m.deviceAuth = deviceAuth

	return ActionInitiateOAuth{
		DeviceCode:      deviceAuth.DeviceCode,
		UserCode:        deviceAuth.UserCode,
		VerificationURL: deviceAuth.VerificationURI,
		ExpiresIn:       int(time.Until(deviceAuth.ExpiresIn)),
		Interval:        int(deviceAuth.Interval),
	}
}

func (m *OAuthBuilder) startPolling(deviceCode string, expiresIn int) tea.Cmd {
	return func() tea.Msg {
		return ActionCompleteOAuth{}
	}
}

func (m *OAuthBuilder) stopPolling() tea.Msg {
	if m.cancelFunc != nil {
		m.cancelFunc()
	}

	return nil
}

func (m *OAuthBuilder) supportsCodeEntry() bool {
	return true
}

func (m *OAuthBuilder) submitCode(string) tea.Cmd {
	return nil
}
