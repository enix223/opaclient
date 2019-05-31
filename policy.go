package opaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-policy-agent/opa/server/types"
)

// ListPolicies list policies
func (c *Client) ListPolicies() (*types.PolicyListResponseV1, error) {
	req, err := http.NewRequest(http.MethodGet, c.BuildURL(APIListPolicies), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrServerError
	}

	var result types.PolicyListResponseV1
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPolicy get specific policy
func (c *Client) GetPolicy(policyID string) (*types.PolicyGetResponseV1, error) {
	req, err := http.NewRequest(http.MethodGet, c.BuildURL(APIPolicy, policyID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.PolicyGetResponseV1
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		err = decoder.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return nil, ErrServerError
}

// CreateUpdatePolicy create policy if not exist, update it if exist
func (c *Client) CreateUpdatePolicy(policyID, policy string, pretty, metrics bool) (*types.PolicyPutResponseV1, error) {
	body := bytes.NewBufferString(policy)
	req, err := http.NewRequest(http.MethodPut, c.BuildURL(APIPolicy, policyID), body)
	if err != nil {
		return nil, err
	}
	if pretty {
		req.URL.Query().Add("pretty", "true")
	}
	if metrics {
		req.URL.Query().Add("metrics", "true")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.PolicyPutResponseV1
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		err = decoder.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return nil, ErrServerError
}

// DeletePolicy delete policy
func (c *Client) DeletePolicy(policyID string, pretty, metrics bool) (*types.PolicyDeleteResponseV1, error) {
	req, err := http.NewRequest(http.MethodDelete, c.BuildURL(APIPolicy, policyID), nil)
	if err != nil {
		return nil, err
	}
	if pretty {
		req.URL.Query().Add("pretty", "true")
	}
	if metrics {
		req.URL.Query().Add("metrics", "true")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.PolicyDeleteResponseV1
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		err = decoder.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	}

	return nil, ErrServerError
}
