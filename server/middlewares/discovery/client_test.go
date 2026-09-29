package discovery_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zeiss/builder/server/middlewares/discovery"
)

func TestClient_Discover(t *testing.T) {
	tests := []struct {
		name       string
		response   string
		oidcIssuer string
		apiURL     string
		statusCode int
	}{
		{
			name:       "successful discovery",
			response:   `{"oidc_issuer": "https://example.com", "api_url": "https://example.com/api/v1"}`,
			oidcIssuer: "https://example.com",
			apiURL:     "https://example.com/api/v1",
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			client := discovery.NewClient(discovery.WithClient(server.Client()))

			config, err := client.Discover(context.Background(), server.URL)
			require.NoError(t, err)
			require.NotNil(t, config)
			require.Equal(t, tt.apiURL, config.ApiURL)
			require.Equal(t, tt.oidcIssuer, config.OidcIssuer)
		})
	}
}
