package tests

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"
)

var originalClient = apis.Client

type MockHTTPClientAuth struct {
	mockResponse types.AuthResponse
	mockError    error
}

func (c *MockHTTPClientAuth) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(`{"AccessToken": "mocked_token"}`))),
	}, nil
}

func TestAuthSuccess(t *testing.T) {
	// Store the original HTTP client and restore it after the test.
	defer func() {
		apis.Client = originalClient
	}()

	// Create an instance of the MockHTTPClient.
	mockClient := &MockHTTPClientAuth{}

	// Replace the global HTTP client with the mock client for this test.
	apis.Client = mockClient

	input := types.AuthInput{
		Username: "testuser",
		Password: "testpassword",
		Endpoint: "https://example.com",
	}

	// Call the Auth function with the test input.
	authResponse, err := apis.Auth(input, mockClient)

	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if authResponse.AccessToken != "mocked_token" {
		t.Fatalf("Expected token 'mocked_token', but got %s", authResponse.AccessToken)
	}
}

func TestAuthFailed(t *testing.T) {

	defer func() {
		apis.Client = originalClient
	}()
	mockClient := &MockHTTPClientAuth{
		mockError: fmt.Errorf("mocked error"),
	}

	apis.Client = mockClient

	input := types.AuthInput{
		Username: "testuser",
		Password: "testpassword",
		Endpoint: "https://example.com",
	}

	_, err := apis.Auth(input, mockClient)
	// fmt.Println("Response:", res)
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
	expectedErrorMessage := "mocked error"
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}
