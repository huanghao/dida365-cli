package dida

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const oauthTokenURL = "https://dida365.com/oauth/token"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

type Project struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind,omitempty"`
	Color     string `json:"color,omitempty"`
	Closed    bool   `json:"closed,omitempty"`
	SortOrder int64  `json:"sortOrder,omitempty"`
}

type Task struct {
	ID            string     `json:"id,omitempty"`
	ProjectID     string     `json:"projectId,omitempty"`
	Title         string     `json:"title,omitempty"`
	Content       string     `json:"content,omitempty"`
	Desc          string     `json:"desc,omitempty"`
	IsAllDay      bool       `json:"isAllDay,omitempty"`
	StartDate     string     `json:"startDate,omitempty"`
	DueDate       string     `json:"dueDate,omitempty"`
	TimeZone      string     `json:"timeZone,omitempty"`
	Reminders     []string   `json:"reminders,omitempty"`
	RepeatFlag    string     `json:"repeatFlag,omitempty"`
	Priority      int        `json:"priority,omitempty"`
	Status        int        `json:"status,omitempty"`
	CompletedTime string     `json:"completedTime,omitempty"`
	SortOrder     int64      `json:"sortOrder,omitempty"`
	Items         []TaskItem `json:"items,omitempty"`
}

type TaskItem struct {
	ID            string `json:"id,omitempty"`
	Status        int    `json:"status,omitempty"`
	Title         string `json:"title,omitempty"`
	SortOrder     int64  `json:"sortOrder,omitempty"`
	StartDate     string `json:"startDate,omitempty"`
	IsAllDay      bool   `json:"isAllDay,omitempty"`
	TimeZone      string `json:"timeZone,omitempty"`
	CompletedTime string `json:"completedTime,omitempty"`
}

type ProjectData struct {
	Project Project `json:"project"`
	Tasks   []Task  `json:"tasks"`
	Columns []any   `json:"columns,omitempty"`
}

type ErrorResponse struct {
	Code    any    `json:"code"`
	Message string `json:"message"`
}

func NewClient(baseURL, token string) *Client {
	base := strings.TrimSpace(baseURL)
	if base == "" {
		base = "https://api.dida365.com/open/v1"
	}
	return &Client{
		BaseURL: strings.TrimRight(base, "/"),
		Token:   strings.TrimSpace(token),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func BuildAuthorizeURL(clientID, redirectURI, scope, state string) (string, error) {
	if strings.TrimSpace(clientID) == "" {
		return "", fmt.Errorf("client_id is required")
	}
	if strings.TrimSpace(redirectURI) == "" {
		return "", fmt.Errorf("redirect_uri is required")
	}
	if strings.TrimSpace(scope) == "" {
		scope = "tasks:read tasks:write"
	}
	u, err := url.Parse("https://dida365.com/oauth/authorize")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("client_id", clientID)
	q.Set("scope", scope)
	q.Set("state", state)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func ExchangeAuthorizationCode(clientID, clientSecret, redirectURI, code, scope string) (*TokenResponse, error) {
	if strings.TrimSpace(clientID) == "" || strings.TrimSpace(clientSecret) == "" {
		return nil, fmt.Errorf("client credentials are required")
	}
	if strings.TrimSpace(redirectURI) == "" {
		return nil, fmt.Errorf("redirect_uri is required")
	}
	if strings.TrimSpace(code) == "" {
		return nil, fmt.Errorf("code is required")
	}
	if strings.TrimSpace(scope) == "" {
		scope = "tasks:read tasks:write"
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("scope", scope)
	data.Set("redirect_uri", redirectURI)

	req, err := http.NewRequest(http.MethodPost, oauthTokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+basicAuth(clientID, clientSecret))

	httpClient := &http.Client{Timeout: 30 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("token exchange failed: %s", summarizeHTTPError(resp.StatusCode, body))
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("decode token response: %w", err)
	}
	return &tokenResp, nil
}

func (c *Client) GetProjects() ([]Project, error) {
	var projects []Project
	if err := c.doJSON(http.MethodGet, "/project", nil, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}

func (c *Client) GetProjectData(projectID string) (*ProjectData, error) {
	var data ProjectData
	path := fmt.Sprintf("/project/%s/data", url.PathEscape(strings.TrimSpace(projectID)))
	if err := c.doJSON(http.MethodGet, path, nil, &data); err != nil {
		return nil, err
	}
	return &data, nil
}

func (c *Client) GetTask(projectID, taskID string) (*Task, error) {
	var task Task
	path := fmt.Sprintf("/project/%s/task/%s", url.PathEscape(strings.TrimSpace(projectID)), url.PathEscape(strings.TrimSpace(taskID)))
	if err := c.doJSON(http.MethodGet, path, nil, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) CreateTask(input Task) (*Task, error) {
	var task Task
	if err := c.doJSON(http.MethodPost, "/task", input, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) UpdateTask(taskID string, input Task) (*Task, error) {
	var task Task
	path := fmt.Sprintf("/task/%s", url.PathEscape(strings.TrimSpace(taskID)))
	if err := c.doJSON(http.MethodPost, path, input, &task); err != nil {
		return nil, err
	}
	return &task, nil
}

func (c *Client) CompleteTask(projectID, taskID string) error {
	path := fmt.Sprintf("/project/%s/task/%s/complete", url.PathEscape(strings.TrimSpace(projectID)), url.PathEscape(strings.TrimSpace(taskID)))
	return c.doJSON(http.MethodPost, path, map[string]any{}, nil)
}

func (c *Client) DeleteTask(projectID, taskID string) error {
	path := fmt.Sprintf("/project/%s/task/%s", url.PathEscape(strings.TrimSpace(projectID)), url.PathEscape(strings.TrimSpace(taskID)))
	return c.doJSON(http.MethodDelete, path, nil, nil)
}

func (c *Client) doJSON(method, path string, requestBody any, out any) error {
	if c == nil {
		return fmt.Errorf("client is nil")
	}
	if strings.TrimSpace(c.Token) == "" {
		return fmt.Errorf("missing access token; run 'dida auth status' and 'dida auth token --code ...'")
	}

	var body io.Reader
	if requestBody != nil {
		payload, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, c.BaseURL+path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("api request failed: %s", summarizeHTTPError(resp.StatusCode, respBody))
	}

	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func summarizeHTTPError(status int, body []byte) string {
	if len(body) == 0 {
		return fmt.Sprintf("status=%d", status)
	}
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err == nil {
		if strings.TrimSpace(errResp.Message) != "" {
			return fmt.Sprintf("status=%d code=%v message=%s", status, errResp.Code, errResp.Message)
		}
	}
	msg := strings.TrimSpace(string(body))
	if len(msg) > 280 {
		msg = msg[:280] + "..."
	}
	return fmt.Sprintf("status=%d body=%s", status, msg)
}

func basicAuth(clientID, clientSecret string) string {
	raw := clientID + ":" + clientSecret
	return base64.StdEncoding.EncodeToString([]byte(raw))
}
