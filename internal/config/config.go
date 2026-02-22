package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultConfigDirName  = "dida365-cli"
	defaultConfigFileName = "config.json"
)

type OAuthConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	Scope        string `json:"scope,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
}

type Config struct {
	APIBaseURL string      `json:"api_base_url"`
	OAuth      OAuthConfig `json:"oauth"`
	Token      Token       `json:"token"`
}

type Store struct {
	path string
}

func NewStore(pathOverride string) (*Store, error) {
	path, err := resolvePath(pathOverride)
	if err != nil {
		return nil, err
	}
	return &Store{path: path}, nil
}

func (s *Store) Path() string {
	if s == nil {
		return ""
	}
	return s.path
}

func (s *Store) Load() (*Config, error) {
	if s == nil || s.path == "" {
		return nil, errors.New("config store not initialized")
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := defaultConfig()
			return &cfg, nil
		}
		return nil, err
	}

	cfg := defaultConfig()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	if strings.TrimSpace(cfg.APIBaseURL) == "" {
		cfg.APIBaseURL = defaultConfig().APIBaseURL
	}
	return &cfg, nil
}

func (s *Store) Save(cfg *Config) error {
	if s == nil || s.path == "" {
		return errors.New("config store not initialized")
	}
	if cfg == nil {
		return errors.New("config is nil")
	}
	if strings.TrimSpace(cfg.APIBaseURL) == "" {
		cfg.APIBaseURL = defaultConfig().APIBaseURL
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	payload, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	payload = append(payload, '\n')
	return os.WriteFile(s.path, payload, 0o600)
}

func (s *Store) SetOAuth(clientID, clientSecret, redirectURI string) (*Config, error) {
	cfg, err := s.Load()
	if err != nil {
		return nil, err
	}
	cfg.OAuth.ClientID = strings.TrimSpace(clientID)
	cfg.OAuth.ClientSecret = strings.TrimSpace(clientSecret)
	cfg.OAuth.RedirectURI = strings.TrimSpace(redirectURI)
	if err := s.Save(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Store) SetToken(token Token) (*Config, error) {
	cfg, err := s.Load()
	if err != nil {
		return nil, err
	}

	token.AccessToken = strings.TrimSpace(token.AccessToken)
	token.RefreshToken = strings.TrimSpace(token.RefreshToken)
	token.TokenType = strings.TrimSpace(token.TokenType)
	token.Scope = strings.TrimSpace(token.Scope)
	if token.ExpiresIn > 0 && token.ExpiresAt == "" {
		token.ExpiresAt = time.Now().UTC().Add(time.Duration(token.ExpiresIn) * time.Second).Format(time.RFC3339)
	}
	cfg.Token = token

	if err := s.Save(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (s *Store) ClearToken() (*Config, error) {
	cfg, err := s.Load()
	if err != nil {
		return nil, err
	}
	cfg.Token = Token{}
	if err := s.Save(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func resolvePath(pathOverride string) (string, error) {
	if v := strings.TrimSpace(pathOverride); v != "" {
		return expandHome(v), nil
	}
	if v := strings.TrimSpace(os.Getenv("DIDA_CONFIG")); v != "" {
		return expandHome(v), nil
	}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, defaultConfigDirName, defaultConfigFileName), nil
}

func expandHome(path string) string {
	if path == "" || path[0] != '~' {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	if strings.HasPrefix(path, "~/") {
		return filepath.Join(home, path[2:])
	}
	return path
}

func defaultConfig() Config {
	return Config{APIBaseURL: "https://api.dida365.com/open/v1"}
}
