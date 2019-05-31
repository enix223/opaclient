package opaclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/open-policy-agent/opa/server/types"
)

// Compile compile
func (c *Client) Compile(query string, input *interface{}, unknowns *[]string) (*types.CompileResponseV1, error) {
	request := types.CompileRequestV1{
		Input:    input,
		Query:    query,
		Unknowns: unknowns,
	}
	buf := make([]byte, 0)
	body := bytes.NewBuffer(buf)
	encoder := json.NewEncoder(body)
	err := encoder.Encode(&request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, c.BuildURL(APICompile), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.token))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var result types.CompileResponseV1
		decoder := json.NewDecoder(resp.Body)
		defer resp.Body.Close()
		err = decoder.Decode(&result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	} else if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	}

	return nil, ErrServerError
}
