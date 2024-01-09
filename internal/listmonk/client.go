package listmonk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	Host     string
	Username string
	Password string
	Headers  map[string]string
}

type Template struct {
	ID        int    `json:"id,omitempty"`
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	Name      string `json:"name"`
	Body      string `json:"body"`
	Type      string `json:"type"`
	IsDefault bool   `json:"is_default,omitempty"`
	Subject   string `json:"subject"`
}

type TemplatesResponse struct {
	Data []Template `json:"data"`
}

type TemplateResponse struct {
	Data Template `json:"data"`
}

type DeleteResponse struct {
	Data bool `json:"data"`
}

func NewClient(host, username, password string, headers map[string]string) *Client {
	return &Client{
		Host:     host,
		Username: username,
		Password: password,
		Headers:  headers,
	}
}

func (c *Client) sendRequest(method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.Headers {
		// remove qoutes from header values
		req.Header.Add(k, v)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: %s\n%s", resp.Status, string(responseBody))
	}

	return responseBody, nil
}

func (c *Client) GetTemplates() (*[]Template, error) {
	url := c.Host + "/api/templates"
	responseBody, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var templates TemplatesResponse
	err = json.Unmarshal(responseBody, &templates)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return &templates.Data, nil
}

// GetTemplate returns a template by ID.
func (c *Client) GetTemplate(id int) (*Template, error) {
	url := fmt.Sprintf("%s/api/templates/%d", c.Host, id)
	responseBody, err := c.sendRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var template TemplateResponse
	err = json.Unmarshal(responseBody, &template)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w\n%s", err, responseBody)
	}

	return &template.Data, nil
}

func (c *Client) CreateTemplate(template *Template) (*Template, error) {
	url := c.Host + "/api/templates"
	templateJSON, err := json.Marshal(&template)
	if err != nil {
		return nil, fmt.Errorf("error marshalling template: %w", err)
	}

	responseBody, err := c.sendRequest("POST", url, bytes.NewBuffer(templateJSON))
	if err != nil {
		return nil, err
	}

	var r TemplateResponse
	err = json.Unmarshal(responseBody, &r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w\n%s", err, responseBody)
	}

	return &r.Data, nil
}

func (c *Client) UpdateTemplate(template *Template) (*Template, error) {
	url := fmt.Sprintf("%s/api/templates/%d", c.Host, template.ID)
	templateJSON, err := json.Marshal(&template)
	if err != nil {
		return nil, fmt.Errorf("error marshalling template: %w", err)
	}

	responseBody, err := c.sendRequest("PUT", url, bytes.NewBuffer(templateJSON))
	if err != nil {
		return nil, err
	}

	var r TemplateResponse
	err = json.Unmarshal(responseBody, &r)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return &r.Data, nil
}

func (c *Client) DeleteTemplate(id int) error {
	url := fmt.Sprintf("%s/api/templates/%d", c.Host, id)
	_, err := c.sendRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	return nil
}
