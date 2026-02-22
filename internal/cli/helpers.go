package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/huanghao/dida365-cli/internal/config"
	"github.com/huanghao/dida365-cli/internal/dida"
)

func loadConfig(app *App) (*config.Config, error) {
	if app == nil || app.ConfigStore == nil {
		return nil, fmt.Errorf("config store not initialized")
	}
	cfg, err := app.ConfigStore.Load()
	if err != nil {
		return nil, err
	}
	if envBase := strings.TrimSpace(os.Getenv("DIDA_API_BASE_URL")); envBase != "" {
		cfg.APIBaseURL = envBase
	}
	return cfg, nil
}

func resolveAccessToken(cfg *config.Config) string {
	if v := strings.TrimSpace(os.Getenv("DIDA_ACCESS_TOKEN")); v != "" {
		return v
	}
	if cfg == nil {
		return ""
	}
	return strings.TrimSpace(cfg.Token.AccessToken)
}

func newAPIClient(cfg *config.Config) *dida.Client {
	if cfg == nil {
		return dida.NewClient("", "")
	}
	return dida.NewClient(cfg.APIBaseURL, resolveAccessToken(cfg))
}

func ensureOAuthConfig(cfg *config.Config) error {
	if cfg == nil {
		return fmt.Errorf("config is nil")
	}
	if strings.TrimSpace(cfg.OAuth.ClientID) == "" {
		return fmt.Errorf("missing oauth.client_id; run 'dida auth init --client-id ... --client-secret ... --redirect-uri ...'")
	}
	if strings.TrimSpace(cfg.OAuth.ClientSecret) == "" {
		return fmt.Errorf("missing oauth.client_secret; run 'dida auth init --client-id ... --client-secret ... --redirect-uri ...'")
	}
	if strings.TrimSpace(cfg.OAuth.RedirectURI) == "" {
		return fmt.Errorf("missing oauth.redirect_uri; run 'dida auth init --client-id ... --client-secret ... --redirect-uri ...'")
	}
	return nil
}
