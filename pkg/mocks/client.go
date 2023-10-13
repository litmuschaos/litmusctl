package mocks

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

type MockHTTPClient struct {
	MockResponse types.AuthResponse
	MockError    error
}

func (c *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	if c.MockError != nil {
		return nil, c.MockError
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"AccessToken": "mocked_token"}`))),
	}, nil
}
