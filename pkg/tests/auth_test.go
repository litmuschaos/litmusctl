package tests

import (
	"fmt"
	"testing"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/mocks"
	"github.com/litmuschaos/litmusctl/pkg/types"
)

var originalClient = apis.Client

func TestAuthSuccess(t *testing.T) {
	// Storing the original HTTP client and restoring it after the test.
	defer func() {
		apis.Client = originalClient
	}()

	//  MockHTTPClient instance
	mockClient := &mocks.MockHTTPClient{}

	// Replacing the global HTTP client with the mock client for this test.
	apis.Client = mockClient

	input := types.AuthInput{
		Username: "testuser",
		Password: "testpassword",
		Endpoint: "https://example.com",
	}

	// Calling Auth with test input.
	authResponse, err := apis.Auth(input, mockClient)
	fmt.Println("Response:", authResponse)

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
	mockClient := &mocks.MockHTTPClient{
		MockError: fmt.Errorf("mocked error"),
	}

	apis.Client = mockClient

	input := types.AuthInput{
		Username: "testuser",
		Password: "testpassword",
		Endpoint: "https://example.com",
	}

	res, err := apis.Auth(input, mockClient)
	fmt.Println("Response:", res)
	if err == nil {
		t.Fatal("Expected an error, but got nil")
	}
	expectedErrorMessage := "mocked error"
	if err.Error() != expectedErrorMessage {
		t.Fatalf("Expected error message '%s', but got '%s'", expectedErrorMessage, err.Error())
	}
}
