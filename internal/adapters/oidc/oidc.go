package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/zeiss/builder/internal/models"
	"github.com/zeiss/builder/internal/ports"
	"github.com/zeiss/builder/server/middlewares/discovery"
	"github.com/zeiss/fiber-goth/v3/providers"

	"github.com/zeiss/pkg/cast"
	"github.com/zeiss/pkg/utilx"
	"golang.org/x/oauth2"
)

var (
	ErrNoVerifiedPrimaryEmail = errors.New("builder: no verified primary email found")
	ErrFailedFetchUser        = errors.New("builder: no failed to fetch user")
	ErrNotAllowedOrg          = errors.New("builder: user not in allowed org")
	ErrNoName                 = errors.New("builder: user has no display name set")
	ErrMissingIDToken         = errors.New("builder: no id token found")
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

var _ ports.DeviceAuthRepository = (*oidcProvider)(nil)

// DefaultScopes holds the default scopes used for GitHub.
var DefaultScopes = []string{"openid", "email", "offline_access"}

type oidcProvider struct {
	id              string
	name            string
	clientID        string
	callbackURL     string
	url             string
	allowedOrgs     []string
	providerType    models.AuthProviderType
	client          *http.Client
	scopes          []string
	wellKnownConfig *discovery.WellKnownConfig
}

// Opt is a function that configures the GitHub provider.
type Opt func(*oidcProvider)

// WithScopes sets the scopes for the GitHub provider.
func WithScopes(scopes ...string) Opt {
	return func(p *oidcProvider) {
		p.scopes = scopes
	}
}

// WithWellKnownConfig sets the well-known configuration for the GitHub provider.
func WithWellKnownConfig(wellKnownConfig *discovery.WellKnownConfig) Opt {
	return func(p *oidcProvider) {
		p.wellKnownConfig = wellKnownConfig
	}
}

// New creates a new GitHub provider.
func New(url, clientID string, opts ...Opt) *oidcProvider {
	p := &oidcProvider{
		allowedOrgs:  []string{},
		client:       DefaultClient,
		clientID:     clientID,
		id:           "oidc",
		name:         "ODIC",
		providerType: models.AuthProviderTypeOAuth2,
		scopes:       DefaultScopes,
		url:          url,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Begin is a method that begins the device authentication process.
func (o *oidcProvider) Begin(ctx context.Context) (*models.DeviceAuth, error) {
	cfg := newConfig(o, o.wellKnownConfig.OidcIssuer, o.scopes...)

	resp, err := cfg.DeviceAuth(ctx) // PKCE flow
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
func (o *oidcProvider) Finish(ctx context.Context, deviceAuth *models.DeviceAuth) (*models.Account, error) {
	cfg := newConfig(o, o.wellKnownConfig.OidcIssuer, o.scopes...)

	code := &oauth2.DeviceAuthResponse{
		DeviceCode: deviceAuth.DeviceCode,
		Expiry:     deviceAuth.ExpiresIn,
	}

	token, err := cfg.DeviceAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, ErrMissingIDToken
	}

	provider, err := oidc.NewProvider(ctx, o.wellKnownConfig.OidcIssuer)
	if err != nil {
		return nil, err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: o.clientID})
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
		Provider:     o.id, // this is an internal reference
		AccessToken:  cast.Ptr(token.AccessToken),
		RefreshToken: cast.Ptr(token.RefreshToken),
		ExpiresAt:    cast.Ptr(token.Expiry),
		TokenType:    cast.Ptr(token.TokenType),
		IDToken:      cast.Ptr(rawIDToken), // save id token for the API
	}

	return account, nil
}

// Refresh refreshes the access token for the given account.
func (o *oidcProvider) Refresh(ctx context.Context, account *models.Account) error {
	provider, err := oidc.NewProvider(ctx, o.wellKnownConfig.OidcIssuer)
	if err != nil {
		return err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: o.clientID})
	_, err = idTokenVerifier.Verify(ctx, cast.Value(account.IDToken))
	if utilx.IsNil(err) {
		return nil
	}

	if _, ok := errors.AsType[*oidc.TokenExpiredError](err); !ok {
		return providers.ErrFailedVerifyToken
	}

	cfg := newConfig(o, o.wellKnownConfig.OidcIssuer, o.scopes...)
	cfg.Endpoint = provider.Endpoint()

	ts := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: cast.Value(account.RefreshToken)})
	token, err := ts.Token()
	if err != nil {
		return err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return ErrMissingIDToken
	}

	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return providers.ErrFailedVerifyToken
	}

	var claims struct {
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Verified bool     `json:"email_verified"`
		Groups   []string `json:"groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return err
	}

	account.AccessToken = cast.Ptr(token.AccessToken)
	account.RefreshToken = cast.Ptr(token.RefreshToken)
	account.ExpiresAt = cast.Ptr(token.Expiry)
	account.TokenType = cast.Ptr(token.TokenType)
	account.IDToken = cast.Ptr(rawIDToken)

	return nil
}

func newConfig(o *oidcProvider, wellKnownURL string, scopes ...string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:    o.clientID,
		RedirectURL: o.callbackURL,
		Endpoint:    urlEndpointConfig(wellKnownURL),
		Scopes:      scopes,
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
