package opaclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

// Client client for opa server
type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
}

// ClientConfig client config
type ClientConfig struct {
	TLSConfig *tls.Config
	Token     string
	BaseURL   string
}

const (
	// APIListPolicies list policies
	APIListPolicies = "/v1/policies"
	// APIPolicy get/create/update policy
	APIPolicy = "/v1/policies/%s"
	// APIData get data
	APIData = "/v1/data/%s"
	// APIDataWebhook get data (webhook)
	APIDataWebhook = "/v0/data/%s"
	// APISimpleQuery simple query
	APISimpleQuery = "/"
	// APIAdHocQuery ad-hoc query
	APIAdHocQuery = "/v1/query"
	// APICompile compile
	APICompile = "/v1/compile"
)

// NewClient create a client base on config
func NewClient(config *ClientConfig) *Client {
	c := &Client{
		token:   config.Token,
		baseURL: config.BaseURL,
	}

	if config.BaseURL == "" {
		return nil
	}

	if config.TLSConfig != nil {
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: config.TLSConfig,
			},
		}
	} else {
		c.httpClient = http.DefaultClient
	}

	return c
}

// BuildURL build url with host and path
func (c *Client) BuildURL(path string, args ...interface{}) string {
	p := c.baseURL + path
	return fmt.Sprintf(p, args...)
}
