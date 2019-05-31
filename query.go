package opaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/open-policy-agent/opa/server/types"
)

// ExecSimpleQuery execute simple query
func (c *Client) ExecSimpleQuery(query string, pretty bool) ([]byte, error) {
	body := bytes.NewBufferString(query)
	req, err := http.NewRequest(http.MethodGet, c.BuildURL(APISimpleQuery), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	return nil, ErrServerError
}

// ExecAdHocQuery execute ad-hoc query
func (c *Client) ExecAdHocQuery(query string, pretty, explain, metrics, watch bool) (*types.QueryResponseV1, error) {
	body := bytes.NewBufferString(fmt.Sprintf(`{"query": %s}`, query))
	req, err := http.NewRequest(http.MethodPost, c.BuildURL(APISimpleQuery), body)
	if err != nil {
		return nil, err
	}
	if pretty {
		req.URL.Query().Add("pretty", "true")
	}
	if metrics {
		req.URL.Query().Add("metrics", "true")
	}
	if watch {
		req.URL.Query().Add("watch", "true")
	}
	if explain {
		req.URL.Query().Add("explain", "true")
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("Content-Type", "text/plain")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.QueryResponseV1
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		err = decoder.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	} else if resp.StatusCode == http.StatusNotImplemented {
		return nil, ErrStreamingNotImplemented
	}

	return nil, ErrServerError
}
