package models

import "time"

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

// AuthProviderType represents the type of authentication provider.
type AuthProviderType string

const (
	// AuthProviderTypeOAuth2 represents an OAuth2 account type.
	AuthProviderTypeOAuth2 AuthProviderType = "oauth2"
	// AuthProviderTypeOIDC represents an OIDC account type.
	AuthProviderTypeOIDC AuthProviderType = "oidc"
	// AuthProviderTypeSAML represents a SAML account type.
	AuthProviderTypeSAML AuthProviderType = "saml"
	// AuthProviderTypeEmail represents an email account type.
	AuthProviderTypeEmail AuthProviderType = "email"
	// AuthProviderTypeWebAuthn represents a WebAuthn account type.
	AuthProviderTypeWebAuthn AuthProviderType = "webauthn"
	// AuthProviderTypeUnknown represents an unknown account type.
	AuthProviderTypeUnknown AuthProviderType = "unknow"
)
