package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/huanghao/dida365-cli/internal/config"
	"github.com/huanghao/dida365-cli/internal/dida"
	"github.com/huanghao/dida365-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewAuthCommand(app *App) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "OAuth setup and token management",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newAuthInitCommand(app))
	cmd.AddCommand(newAuthLoginCommand(app))
	cmd.AddCommand(newAuthTokenCommand(app))
	cmd.AddCommand(newAuthRefreshCommand(app))
	cmd.AddCommand(newAuthStatusCommand(app))
	cmd.AddCommand(newAuthLogoutCommand(app))

	return cmd
}

func newAuthInitCommand(app *App) *cobra.Command {
	var clientID string
	var clientSecret string
	var redirectURI string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Save OAuth client settings",
		RunE: func(cmd *cobra.Command, args []string) error {
			if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" || strings.TrimSpace(redirectURI) == "" {
				return fmt.Errorf("--client-id, --client-secret, and --redirect-uri are required")
			}
			if app.DryRun {
				if asJSON {
					return output.PrintJSON(app.Out, map[string]any{
						"action":        "auth_init",
						"dry_run":       true,
						"client_id_set": true,
						"redirect_uri":  redirectURI,
					})
				}
				fmt.Fprintln(app.Out, "Would save oauth client settings")
				return nil
			}
			if _, err := app.ConfigStore.SetOAuth(clientID, clientSecret, redirectURI); err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"ok":          true,
					"action":      "auth_init",
					"config_path": app.ConfigStore.Path(),
				})
			}
			fmt.Fprintf(app.Out, "Saved OAuth settings to %s\n", app.ConfigStore.Path())
			return nil
		},
	}

	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client id")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret")
	cmd.Flags().StringVar(&redirectURI, "redirect-uri", "", "OAuth redirect URI")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}

func newAuthLoginCommand(app *App) *cobra.Command {
	var scope string
	var state string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Print authorization URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if strings.TrimSpace(cfg.OAuth.ClientID) == "" || strings.TrimSpace(cfg.OAuth.RedirectURI) == "" {
				return fmt.Errorf("missing oauth client settings; run 'dida365-cli auth init ...'")
			}
			if strings.TrimSpace(scope) == "" {
				scope = "tasks:read tasks:write"
			}
			if strings.TrimSpace(state) == "" {
				state = fmt.Sprintf("dida-%d", time.Now().Unix())
			}
			authURL, err := dida.BuildAuthorizeURL(cfg.OAuth.ClientID, cfg.OAuth.RedirectURI, scope, state)
			if err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"action":            "auth_login",
					"authorization_url": authURL,
					"next_command":      "dida365-cli auth token --code <authorization_code>",
				})
			}
			fmt.Fprintln(app.Out, authURL)
			fmt.Fprintln(app.Out)
			fmt.Fprintln(app.Out, "Open the URL, authorize, then run:")
			fmt.Fprintln(app.Out, "  dida365-cli auth token --code <authorization_code>")
			return nil
		},
	}

	cmd.Flags().StringVar(&scope, "scope", "tasks:read tasks:write", "OAuth scope")
	cmd.Flags().StringVar(&state, "state", "", "OAuth state parameter")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}

func newAuthTokenCommand(app *App) *cobra.Command {
	var code string
	var scope string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "token",
		Short: "Exchange authorization code for access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if err := ensureOAuthConfig(cfg); err != nil {
				return err
			}
			if strings.TrimSpace(code) == "" {
				return fmt.Errorf("--code is required")
			}
			if app.DryRun {
				if asJSON {
					return output.PrintJSON(app.Out, map[string]any{
						"action":   "auth_token_exchange",
						"dry_run":  true,
						"code_set": true,
					})
				}
				fmt.Fprintln(app.Out, "Would exchange authorization code for token")
				return nil
			}

			resp, err := dida.ExchangeAuthorizationCode(
				cfg.OAuth.ClientID,
				cfg.OAuth.ClientSecret,
				cfg.OAuth.RedirectURI,
				code,
				scope,
			)
			if err != nil {
				return err
			}

			_, err = app.ConfigStore.SetToken(config.Token{
				AccessToken:  resp.AccessToken,
				RefreshToken: resp.RefreshToken,
				TokenType:    resp.TokenType,
				Scope:        resp.Scope,
				ExpiresIn:    resp.ExpiresIn,
			})
			if err != nil {
				return err
			}

			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"ok":                true,
					"action":            "auth_token_exchange",
					"config_path":       app.ConfigStore.Path(),
					"access_token_set":  strings.TrimSpace(resp.AccessToken) != "",
					"refresh_token_set": strings.TrimSpace(resp.RefreshToken) != "",
					"expires_in":        resp.ExpiresIn,
					"scope":             resp.Scope,
				})
			}
			fmt.Fprintf(app.Out, "Token saved to %s\n", app.ConfigStore.Path())
			return nil
		},
	}

	cmd.Flags().StringVar(&code, "code", "", "Authorization code returned by redirect URI")
	cmd.Flags().StringVar(&scope, "scope", "tasks:read tasks:write", "OAuth scope")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}

func newAuthStatusCommand(app *App) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show auth and token status",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}

			accessToken := resolveAccessToken(cfg)
			payload := map[string]any{
				"config_path":       app.ConfigStore.Path(),
				"api_base_url":      cfg.APIBaseURL,
				"client_id_set":     strings.TrimSpace(cfg.OAuth.ClientID) != "",
				"client_secret_set": strings.TrimSpace(cfg.OAuth.ClientSecret) != "",
				"redirect_uri_set":  strings.TrimSpace(cfg.OAuth.RedirectURI) != "",
				"access_token_set":  strings.TrimSpace(accessToken) != "",
				"refresh_token_set": strings.TrimSpace(cfg.Token.RefreshToken) != "",
				"token_scope":       cfg.Token.Scope,
				"token_expires_at":  cfg.Token.ExpiresAt,
			}
			if asJSON {
				return output.PrintJSON(app.Out, payload)
			}

			fmt.Fprintf(app.Out, "Config: %s\n", app.ConfigStore.Path())
			fmt.Fprintf(app.Out, "API Base URL: %s\n", cfg.APIBaseURL)
			fmt.Fprintf(app.Out, "Client ID set: %t\n", strings.TrimSpace(cfg.OAuth.ClientID) != "")
			fmt.Fprintf(app.Out, "Client Secret set: %t\n", strings.TrimSpace(cfg.OAuth.ClientSecret) != "")
			fmt.Fprintf(app.Out, "Redirect URI set: %t\n", strings.TrimSpace(cfg.OAuth.RedirectURI) != "")

			fmt.Fprintf(app.Out, "Access token set: %t\n", strings.TrimSpace(accessToken) != "")
			if strings.TrimSpace(cfg.Token.ExpiresAt) != "" {
				fmt.Fprintf(app.Out, "Token expires at: %s\n", cfg.Token.ExpiresAt)
			}
			if strings.TrimSpace(cfg.Token.Scope) != "" {
				fmt.Fprintf(app.Out, "Token scope: %s\n", cfg.Token.Scope)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}

func newAuthRefreshCommand(app *App) *cobra.Command {
	var refreshToken string
	var scope string
	var asJSON bool

	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Refresh access token by refresh token",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(app)
			if err != nil {
				return err
			}
			if err := ensureOAuthConfig(cfg); err != nil {
				return err
			}
			tokenToUse := strings.TrimSpace(refreshToken)
			if tokenToUse == "" {
				tokenToUse = strings.TrimSpace(cfg.Token.RefreshToken)
			}
			if tokenToUse == "" {
				return fmt.Errorf("missing refresh token; pass --refresh-token or run 'dida365-cli auth token --code ...' first")
			}
			if app.DryRun {
				if asJSON {
					return output.PrintJSON(app.Out, map[string]any{
						"action":            "auth_refresh",
						"dry_run":           true,
						"refresh_token_set": true,
					})
				}
				fmt.Fprintln(app.Out, "Would refresh access token")
				return nil
			}

			resp, err := dida.ExchangeRefreshToken(
				cfg.OAuth.ClientID,
				cfg.OAuth.ClientSecret,
				tokenToUse,
				scope,
			)
			if err != nil {
				if strings.Contains(err.Error(), "Unauthorized grant type: refresh_token") {
					return fmt.Errorf("refresh_token is not supported by Dida OAuth for this app; re-authenticate with 'dida365-cli auth login' then 'dida365-cli auth token --code ...'")
				}
				return err
			}

			nextRefreshToken := resp.RefreshToken
			if strings.TrimSpace(nextRefreshToken) == "" {
				nextRefreshToken = tokenToUse
			}
			_, err = app.ConfigStore.SetToken(config.Token{
				AccessToken:  resp.AccessToken,
				RefreshToken: nextRefreshToken,
				TokenType:    resp.TokenType,
				Scope:        resp.Scope,
				ExpiresIn:    resp.ExpiresIn,
			})
			if err != nil {
				return err
			}

			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"ok":                true,
					"action":            "auth_refresh",
					"config_path":       app.ConfigStore.Path(),
					"access_token_set":  strings.TrimSpace(resp.AccessToken) != "",
					"refresh_token_set": strings.TrimSpace(nextRefreshToken) != "",
					"expires_in":        resp.ExpiresIn,
					"scope":             resp.Scope,
				})
			}
			fmt.Fprintf(app.Out, "Token refreshed and saved to %s\n", app.ConfigStore.Path())
			return nil
		},
	}

	cmd.Flags().StringVar(&refreshToken, "refresh-token", "", "Refresh token (default: from config)")
	cmd.Flags().StringVar(&scope, "scope", "tasks:read tasks:write", "OAuth scope")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}

func newAuthLogoutCommand(app *App) *cobra.Command {
	var asJSON bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if app.DryRun {
				if asJSON {
					return output.PrintJSON(app.Out, map[string]any{
						"action":  "auth_logout",
						"dry_run": true,
					})
				}
				fmt.Fprintln(app.Out, "Would clear stored token")
				return nil
			}
			if _, err := app.ConfigStore.ClearToken(); err != nil {
				return err
			}
			if asJSON {
				return output.PrintJSON(app.Out, map[string]any{
					"ok":     true,
					"action": "auth_logout",
				})
			}
			fmt.Fprintln(app.Out, "Stored token cleared")
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output JSON")
	return cmd
}
