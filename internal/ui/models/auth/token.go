package auth

import (
	"context"

	"github.com/katallaxie/pkg/cast"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	tea "charm.land/bubbletea/v2"
)

type (
	tokenFetchAccountMsg struct {
		account models.Account
	}
	tokenErrorMsg struct {
		error error
	}
)

type tokenModel struct {
	quitting    bool
	err         error
	accountCtrl ports.AccountController
	ctx         context.Context
}

// NewToken returns a new token model.
func NewToken(ctx context.Context, accountCtrl ports.AccountController) tokenModel {
	return tokenModel{
		quitting:    false,
		ctx:         ctx,
		accountCtrl: accountCtrl,
	}
}

// Init initializes the token model.
func (m tokenModel) Init() tea.Cmd {
	return tea.Batch(m.fetchAccount())
}

// Update handles incoming messages for the token model.
func (m tokenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return m, nil

	case tokenErrorMsg: // handle deployment error
		m.err = msg.error
		m.quitting = true
		return m, tea.Sequence(
			tea.Printf("%s %s", errorMark, m.err),
			tea.Quit,
		)

	case tokenFetchAccountMsg:
		return m, tea.Sequence(
			tea.Println(cast.Value(msg.account.IDToken)),
			tea.Quit,
		)
	}

	return m, nil
}

func (m tokenModel) View() tea.View {
	return tea.NewView("")
}

func (m tokenModel) fetchAccount() tea.Cmd {
	return func() tea.Msg {
		account := &models.Account{}
		err := m.accountCtrl.GetCurrent(m.ctx, account)
		if err != nil {
			return tokenErrorMsg{error: err}
		}

		return tokenFetchAccountMsg{account: *account}
	}
}
