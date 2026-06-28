package dex

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/katallaxie/fiber-goth/v3/providers"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"

	"github.com/katallaxie/pkg/cast"
	"github.com/katallaxie/pkg/conv"
	"golang.org/x/oauth2"
)

var (
	ErrNoVerifiedPrimaryEmail = errors.New("goth: no verified primary email found")
	ErrFailedFetchUser        = errors.New("goth: no failed to fetch user")
	ErrNotAllowedOrg          = errors.New("goth: user not in allowed org")
	ErrNoName                 = errors.New("goth: user has no display name set")
	ErrMissingIDToken         = errors.New("goth: no id token found")
)

// DefaultClient is the default HTTP client used.
// TODO: allows to configure the client via options.
var DefaultClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	},
	Timeout: 10 * time.Second,
}

const NoopEmail = ""

var _ ports.DeviceAuthRepository = (*dexProvider)(nil)

// DefaultScopes holds the default scopes used for GitHub.
var DefaultScopes = []string{"openid", "profile", "email", "offline_access"}

type dexProvider struct {
	id           string
	name         string
	clientID     string
	clientSecret string
	callbackURL  string
	url          string
	allowedOrgs  []string
	providerType models.AuthProviderType
	client       *http.Client
	config       *oauth2.Config
	scopes       []string
}

// Opt is a function that configures the GitHub provider.
type Opt func(*dexProvider)

// WithScopes sets the scopes for the GitHub provider.
func WithScopes(scopes ...string) Opt {
	return func(p *dexProvider) {
		p.config.Scopes = scopes
	}
}

// New creates a new GitHub provider.
func New(url, clientID, clientSecret, callbackURL string, opts ...Opt) *dexProvider {
	p := &dexProvider{
		id:           "dex",
		name:         "Dex",
		clientID:     clientID,
		clientSecret: clientSecret,
		url:          url,
		callbackURL:  callbackURL,
		providerType: models.AuthProviderTypeOAuth2,
		client:       DefaultClient,
		allowedOrgs:  []string{},
		scopes:       DefaultScopes,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.config = newConfig(p, p.scopes...)

	return p
}

// Begin is a method that begins the device authentication process.
func (d *dexProvider) Begin(ctx context.Context) (*models.DeviceAuth, error) {
	resp, err := d.config.DeviceAuth(ctx, oauth2.SetAuthURLParam("client_secret", d.clientSecret))
	if err != nil {
		return nil, err
	}

	return &models.DeviceAuth{
		DeviceCode:              resp.DeviceCode,
		UserCode:                resp.UserCode,
		VerificationURI:         resp.VerificationURI,
		VerificationURIComplete: resp.VerificationURIComplete,
		ExpiresIn:               resp.Expiry,
		Interval:                resp.Interval,
	}, nil
}

// Finish is a method that finishes the device authentication process.
func (d *dexProvider) Finish(ctx context.Context, deviceAuth *models.DeviceAuth) (*models.Account, error) {
	code := &oauth2.DeviceAuthResponse{
		DeviceCode: deviceAuth.DeviceCode,
		Expiry:     deviceAuth.ExpiresIn,
	}

	token, err := d.config.DeviceAccessToken(ctx, code, oauth2.SetAuthURLParam("client_secret", d.clientSecret))
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, ErrMissingIDToken
	}

	provider, err := oidc.NewProvider(ctx, d.url)
	if err != nil {
		return nil, err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: d.clientID})
	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, providers.ErrFailedVerifyToken
	}

	var claims struct {
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Verified bool     `json:"email_verified"`
		Groups   []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	account := &models.Account{
		Type:         models.AccountTypeOAuth2,
		Email:        claims.Email,
		Name:         claims.Name,
		Provider:     d.id, // this is an internal reference
		AccessToken:  cast.Ptr(token.AccessToken),
		RefreshToken: cast.Ptr(token.RefreshToken),
		ExpiresAt:    cast.Ptr(token.Expiry),
		TokenType:    cast.Ptr(token.TokenType),
		IDToken:      cast.Ptr(conv.String(token.Extra("id_token"))), // save id token for the API
	}

	return account, nil
}

func newConfig(d *dexProvider, scopes ...string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     d.clientID,
		ClientSecret: d.clientSecret,
		RedirectURL:  d.callbackURL,
		Endpoint:     urlEndpointConfig(d.url),
		Scopes:       append(DefaultScopes, scopes...),
	}

	return c
}

func urlEndpointConfig(url string) oauth2.Endpoint {
	return oauth2.Endpoint{
		AuthURL:       fmt.Sprintf("%s/authorize", strings.TrimSuffix(url, "/")),
		TokenURL:      fmt.Sprintf("%s/token", strings.TrimSuffix(url, "/")),
		DeviceAuthURL: fmt.Sprintf("%s/device/code", strings.TrimSuffix(url, "/")),
	}
}
