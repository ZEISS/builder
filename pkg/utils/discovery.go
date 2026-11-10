package utils

// WellknownConfigURL returns the URL for the well-known endpoint.
const WellknownConfigURL = "builder-configuration"

// WellknownConfig is the well-known configuration for the builder.
type WellknownConfig struct {
	// OidcIssuer is the URL to the OIDC issuer.
	//
	// Default: "http://builder.internal:5556/dex"
	OidcIssuer string `json:"oidc_issuer" xml:"oidc_issuer" default:"http://builder.internal:5556/dex"`
	// ApiURL is the URL to the API.
	//
	// Default: "http://builder.internal:8080/api/v1"
	ApiURL string `json:"api_url" xml:"api_url" default:"http://builder.internal:8080/api/v1"`
}

// DefaultWellknown returns the default well-known response.
func DefaultWellknown() WellknownConfig {
	return WellknownConfig{
		OidcIssuer: "http://builder.internal:5556/dex",
		ApiURL:     "http://builder.internal:8080/api/v1",
	}
}
