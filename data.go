package opaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/open-policy-agent/opa/server/types"
)

// GetDocument get document with given input and path
func (c *Client) GetDocument(path, input string, pretty, provenance, metrics, instrument, watch bool) (*types.DataResponseV1, error) {
	method := http.MethodGet
	var body io.Reader
	if input != "" {
		method = http.MethodPost
		inp := fmt.Sprintf(`{"input": %s}`, input)
		body = bytes.NewBufferString(inp)
	}

	req, err := http.NewRequest(method, c.BuildURL(APIData, path), body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	if pretty {
		req.URL.Query().Add("pretty", "true")
	}
	if metrics {
		req.URL.Query().Add("metrics", "true")
	}
	if provenance {
		req.URL.Query().Add("provenance", "true")
	}
	if instrument {
		req.URL.Query().Add("instrument", "true")
	}
	if watch {
		req.URL.Query().Add("watch", "")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.DataResponseV1
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

// CreateOverwriteDocument create document or overwrite the existing one
func (c *Client) CreateOverwriteDocument(path string, data string, ifNoneMatch bool) (bool, error) {
	body := bytes.NewBufferString(data)
	req, err := http.NewRequest(http.MethodPut, c.BuildURL(APIData, path), body)
	if err != nil {
		return false, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return false, err
	}

	if resp.StatusCode == http.StatusNoContent {
		return true, nil
	} else if resp.StatusCode == http.StatusNotModified {
		return false, nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return false, ErrBadRequest
	} else if resp.StatusCode == http.StatusNotFound {
		return false, ErrWriteConflict
	}

	return false, ErrServerError
}

// UpdateDocument update document
func (c *Client) UpdateDocument(path string, data string) error {
	body := bytes.NewBufferString(data)
	req, err := http.NewRequest(http.MethodPatch, c.BuildURL(APIData, path), body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return ErrBadRequest
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	return ErrServerError
}

// DeleteDocument delete document
func (c *Client) DeleteDocument(path string) error {
	req, err := http.NewRequest(http.MethodDelete, c.BuildURL(APIData, path), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrNotFound
	}

	return ErrServerError
}
