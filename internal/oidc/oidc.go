package oidc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/zeiss/builder/server/middlewares/discovery"
	"github.com/zeiss/fiber-goth/v3/providers"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/zeiss/pkg/utilx"
	"golang.org/x/oauth2"
)

var (
	// ErrOidcNoVerifiedPrimaryEmail is returned when no verified primary email is found.
	ErrOidcNoVerifiedPrimaryEmail = errors.New("no verified primary email found")
	// ErrOidcFailedFetchUser is returned when the user fails to be fetched.
	ErrOidcFailedFetchUser = errors.New("no failed to fetch user")
	// ErrOidcNotAllowedOrg is returned when the user is not in the allowed organization.
	ErrOidcNotAllowedOrg = errors.New("user not in allowed org")
	// ErrOidcNoName is returned when the user has no display name set.
	ErrOidcNoName = errors.New("user has no display name set")
	// ErrOidcMissingIDToken is returned when no ID token is found.
	ErrOidcMissingIDToken = errors.New("no id token found")
)

// DefaultClient is the default HTTP client used.
// TODO: allows to configure the client via options.
var DefaultClient = &http.Client{
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 20,
	},
	Timeout: 10 * time.Second,
}

// NoopEmail is the email address used when no email is available.
const NoopEmail = ""

// DeviceAuth is a struct that holds the device authentication response data.
type DeviceAuth struct {
	// DeviceCode is the device code used to authenticate the device.
	DeviceCode string `json:"device_code"`
	// UserCode is the user code used to authenticate the device.
	UserCode string `json:"user_code"`
	// VerificationURI is the URI to verify the device code.
	VerificationURI string `json:"verification_uri"`
	// VerificationURIComplete is the URI to verify the device code with the client secret.
	VerificationURIComplete string `json:"verification_uri_complete"`
	// ExpiresIn is the number of seconds until the device code expires.
	ExpiresIn time.Time `json:"expires_in"`
	// Interval is the number of seconds between polling requests for device code status.
	Interval int64 `json:"interval"`
}

// Claims holds the claims from the ID token.
type Claims struct {
	// Name is the user's name.
	Name string `json:"name"`
	// Email is the user's email address.
	Email string `json:"email"`
	// Verified indicates whether the user's email address has been verified.
	Verified bool `json:"email_verified"`
	// Groups is the list of groups the user belongs to.
	Groups []string `json:"groups"`
}

// Token is an alias for oauth2.Token.
type Token struct {
	// ProviderType is the type of provider that issued the token.
	ProviderType string `json:"provider_type"`
	// AccessToken is the token that authorizes and authenticates
	// the requests.
	AccessToken string `json:"access_token"`
	// TokenType is the type of token.
	// The Type method returns either this or "Bearer", the default.
	TokenType string `json:"token_type,omitempty"`
	// RefreshToken is a token that's used by the application
	// (as opposed to the user) to refresh the access token
	// if it expires.
	RefreshToken string `json:"refresh_token,omitempty"`
	// Expiry is the optional expiration time of the access token.
	//
	// If zero, [TokenSource] implementations will reuse the same
	// token forever and RefreshToken or equivalent
	// mechanisms for that TokenSource will not be used.
	Expiry time.Time `json:"expiry,omitempty"`
	// ExpiresIn is the OAuth2 wire format "expires_in" field,
	// which specifies how many seconds later the token expires,
	// relative to an unknown time base approximately around "now".
	// It is the application's responsibility to populate
	// `Expiry` from `ExpiresIn` when required.
	ExpiresIn int64 `json:"expires_in,omitempty"`
	// IDToken is the raw ID token returned from the provider.
	IDToken string `json:"id_token,omitempty"`
	// Claims is the parsed claims from the ID token.
	Claims *Claims `json:"claims,omitempty"`
}

const (
	// ProviderTypeOAuth2 represents an OAuth2 account type.
	ProviderTypeOAuth2 = "oauth2"
	// ProviderTypeOIDC represents an OIDC account type.
	ProviderTypeOIDC = "oidc"
	// ProviderTypeSAML represents a SAML account type.
	ProviderTypeSAML = "saml"
	// ProviderTypeEmail represents an email account type.
	ProviderTypeEmail = "email"
	// ProviderTypeWebAuthn represents a WebAuthn account type.
	ProviderTypeWebAuthn = "webauthn"
	// ProviderTypeUnknown represents an unknown account type.
	ProviderTypeUnknown = "unknown"
)

// OidcProvider is an interface that defines the methods for an OIDC provider.
type OidcProvider interface {
	// Begin is a method that begins the device authentication process.
	Begin(ctx context.Context) (*DeviceAuth, error)
	// Finish is a method that finishes the device authentication process.
	// TODO: could be returned in a token structure
	Finish(ctx context.Context, deviceAuth *DeviceAuth) (*Token, error)
	// Refresh is a method that refreshes the device authentication token.
	Refresh(ctx context.Context, token *Token) error
}

// DefaultScopes holds the default scopes used for GitHub.
var DefaultScopes = []string{"openid", "email", "offline_access"}

var _ OidcProvider = (*oidcProvider)(nil)

type oidcProvider struct {
	id           string
	name         string
	clientID     string
	callbackURL  string
	oidcIssuer   string
	allowedOrgs  []string
	providerType string
	client       *http.Client
	scopes       []string
	discv        discovery.Client
}

// Opt is a function that configures the GitHub provider.
type Opt func(*oidcProvider)

// WithScopes sets the scopes for the GitHub provider.
func WithScopes(scopes ...string) Opt {
	return func(p *oidcProvider) {
		p.scopes = scopes
	}
}

// New creates a new GitHub provider.
func New(oidcIssuer, clientID string, opts ...Opt) *oidcProvider {
	p := &oidcProvider{
		allowedOrgs:  []string{},
		client:       DefaultClient,
		clientID:     clientID,
		id:           "oidc",
		name:         "ODIC",
		providerType: ProviderTypeOAuth2,
		scopes:       DefaultScopes,
		oidcIssuer:   oidcIssuer,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

// Begin is a method that begins the device authentication process.
func (o *oidcProvider) Begin(ctx context.Context) (*DeviceAuth, error) {
	cfg := newConfig(o, o.oidcIssuer, o.scopes...)

	resp, err := cfg.DeviceAuth(ctx) // PKCE flow
	if err != nil {
		return nil, err
	}

	return &DeviceAuth{
		DeviceCode:              resp.DeviceCode,
		UserCode:                resp.UserCode,
		VerificationURI:         resp.VerificationURI,
		VerificationURIComplete: resp.VerificationURIComplete,
		ExpiresIn:               resp.Expiry,
		Interval:                resp.Interval,
	}, nil
}

// Finish is a method that finishes the device authentication process.
func (o *oidcProvider) Finish(ctx context.Context, deviceAuth *DeviceAuth) (*Token, error) {
	cfg := newConfig(o, o.oidcIssuer, o.scopes...)

	code := &oauth2.DeviceAuthResponse{
		DeviceCode: deviceAuth.DeviceCode,
		Expiry:     deviceAuth.ExpiresIn,
	}

	deviceToken, err := cfg.DeviceAccessToken(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := deviceToken.Extra("id_token").(string)
	if !ok {
		return nil, ErrOidcMissingIDToken
	}

	provider, err := oidc.NewProvider(ctx, o.oidcIssuer)
	if err != nil {
		return nil, err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: o.clientID})
	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, providers.ErrFailedVerifyToken
	}

	claims := &Claims{}
	if err := idToken.Claims(claims); err != nil {
		return nil, err
	}

	token := &Token{
		ProviderType: ProviderTypeOAuth2,
		AccessToken:  deviceToken.AccessToken,
		RefreshToken: deviceToken.RefreshToken,
		Expiry:       deviceToken.Expiry,
		TokenType:    deviceToken.TokenType,
		IDToken:      rawIDToken,
		Claims:       claims,
	}

	return token, nil
}

// Refresh refreshes the access token for the given account.
func (o *oidcProvider) Refresh(ctx context.Context, token *Token) error {
	provider, err := oidc.NewProvider(ctx, o.oidcIssuer)
	if err != nil {
		return err
	}

	idTokenVerifier := provider.Verifier(&oidc.Config{ClientID: o.clientID})
	_, err = idTokenVerifier.Verify(ctx, token.IDToken)
	if utilx.IsNil(err) {
		return nil
	}

	cfg := newConfig(o, o.oidcIssuer, o.scopes...)
	cfg.Endpoint = provider.Endpoint()

	ts := cfg.TokenSource(ctx, &oauth2.Token{RefreshToken: token.RefreshToken})
	sourceToken, err := ts.Token()
	if err != nil {
		return err
	}

	rawIDToken, ok := sourceToken.Extra("id_token").(string)
	if !ok {
		return ErrOidcMissingIDToken
	}

	idToken, err := idTokenVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return providers.ErrFailedVerifyToken
	}

	claims := &Claims{}
	if err := idToken.Claims(claims); err != nil {
		return err
	}

	token.AccessToken = sourceToken.AccessToken
	token.RefreshToken = sourceToken.RefreshToken
	token.Expiry = sourceToken.Expiry
	token.TokenType = sourceToken.TokenType
	token.IDToken = rawIDToken
	token.Claims = claims

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
