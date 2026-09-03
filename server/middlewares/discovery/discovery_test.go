package discovery_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/require"
	"github.com/zeiss/builder/server/middlewares/discovery"
)

func TestDefaultConfig(t *testing.T) {
	t.Parallel()

	defaultCfg := discovery.DefaultWellKnownFunc()
	require.Equal(t, defaultCfg.ApiURL, "http://builder.internal:8080/api/v1")
	require.Equal(t, defaultCfg.OidcIssuer, "http://builder.internal:5556/dex")
}

func TestNew(t *testing.T) {
	t.Parallel()

	app := fiber.New()
	app.Get(discovery.WellKnownConfigurationURL, discovery.New())

	req := httptest.NewRequestWithContext(t.Context(), fiber.MethodGet, "/.well-known/builder-configuration", http.NoBody)
	req.Header.Set("Accept", "application/json")

	resp, err := app.Test(req)
	defer resp.Body.Close() //nolint:staticcheck
	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 200, resp.StatusCode, "Status code")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NotNil(t, body)

	wlkn := discovery.WellKnownConfig{}
	err = json.Unmarshal(body, &wlkn)
	require.NoError(t, err)

	require.Equal(t, discovery.DefaultWellKnownFunc(), wlkn)
}
