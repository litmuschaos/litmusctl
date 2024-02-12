package experiment

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/types"

	"github.com/stretchr/testify/assert"
)

type MockHTTPClientCreateExperiment struct {
	mockResponse RunExperimentResponse
	mockError    error
}

func (c *MockHTTPClientCreateExperiment) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func TestCreateExperimentSuccess(t *testing.T) {

	pid := "testprojectid"
	requestData := model.SaveChaosExperimentRequest{
		ID:          "test-experiment-id",
		Name:        "Test Experiment",
		Description: "A test experiment for unit testing",
		Manifest:    "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n  - name: nginx\n    image: nginx",
		InfraID:     "test-infras-id",
		Tags:        []string{"tag1", "tag2"},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	// Create a mock HTTP client with a successful response
	mockResponse := RunExperimentResponse{}
	mockClient := &MockHTTPClientCreateExperiment{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, _ := CreateExperiment(pid, requestData, cred, mockClient)

	// Assertions
	assert.Equal(t, result, mockResponse) // Ensure the result matches the expected response
}

func TestCreateExperimentFailedRequest(t *testing.T) {

	pid := "testprojectid"
	requestData := model.SaveChaosExperimentRequest{
		ID:          "test-experiment-id",
		Name:        "Test Experiment",
		Description: "A test experiment for unit testing",
		Manifest:    "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n  - name: nginx\n    image: nginx",
		InfraID:     "test-infras-id",
		Tags:        []string{"tag1", "tag2"},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientCreateExperiment{
		mockError:    errors.New("Some error occurred"),
		mockResponse: RunExperimentResponse{},
	}

	result, err := CreateExperiment(pid, requestData, cred, mockClient)

	// Assertions
	assert.Error(t, err)                             // Ensure an error occurred
	assert.Equal(t, result, RunExperimentResponse{}) // Ensure the result is an empty response
}

// TestGetExperimentSuccess tests the GetExperiment function with a successful response.
type MockHTTPClientSaveExperiment struct {
	mockResponse SaveExperimentData
	mockError    error
}

func (c *MockHTTPClientSaveExperiment) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}
func TestSaveExperimentSuccess(t *testing.T) {

	pid := "testprojectid"
	requestData := model.SaveChaosExperimentRequest{
		ID:          "test-experiment-id",
		Name:        "Test Experiment",
		Description: "A test experiment for unit testing",
		Manifest:    "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n  - name: nginx\n    image: nginx",
		InfraID:     "test-infras-id",
		Tags:        []string{"tag1", "tag2"},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	// Create a mock HTTP client with a successful response
	mockResponse := SaveExperimentData{}
	mockClient := &MockHTTPClientSaveExperiment{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, _ := SaveExperiment(pid, requestData, cred, mockClient)

	// Assertion
	assert.Equal(t, result, mockResponse)
}

func TestSaveExperimentFailedRequest(t *testing.T) {

	pid := "testprojectid"
	requestData := model.SaveChaosExperimentRequest{
		ID:          "test-experiment-id",
		Name:        "Test Experiment",
		Description: "A test experiment for unit testing",
		Manifest:    "apiVersion: v1\nkind: Pod\nmetadata:\n  name: nginx\nspec:\n  containers:\n  - name: nginx\n    image: nginx",
		InfraID:     "test-infras-id",
		Tags:        []string{"tag1", "tag2"},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	// Create a mock HTTP client with an error response
	mockClient := &MockHTTPClientSaveExperiment{
		mockError:    errors.New("Some error occurred"),
		mockResponse: SaveExperimentData{},
	}

	_, err := SaveExperiment(pid, requestData, cred, mockClient)

	assert.Error(t, err) // Ensure an error occurred

}

type MockHTTPClientRunExperiment struct {
	mockResponse RunExperimentResponse
	mockError    error
}

func (c *MockHTTPClientRunExperiment) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func TestRunExperimentSuccess(t *testing.T) {

	pid := "testprojectid"
	eid := "testexperimentid"
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	// Create a mock HTTP client with a successful response
	mockResponse := RunExperimentResponse{}
	mockClient := &MockHTTPClientRunExperiment{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, err := RunExperiment(pid, eid, cred, mockClient)

	// Assertions
	assert.NoError(t, err)                // Ensure no errors occurred
	assert.Equal(t, result, mockResponse) // Ensure the result matches the expected response
}

func TestRunExperimentFailedRequest(t *testing.T) {
	pid := "testprojectid"
	eid := "testexperimentid"
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientRunExperiment{
		mockError:    errors.New("Some error occurred"),
		mockResponse: RunExperimentResponse{},
	}

	result, _ := RunExperiment(pid, eid, cred, mockClient)

	// Assertions
	// assert.Error(t, err)                                       // Ensure an error occurred
	assert.Empty(t, result.Data.RunExperimentDetails.NotifyID) // Ensure NotifyID is empty
}

type MockHTTPClientGetExperimentList struct {
	mockResponse ExperimentListData
	mockError    error
}

func (c *MockHTTPClientGetExperimentList) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func TestGetExperimentListSuccess(t *testing.T) {
	pid := "testprojectid"
	requestData := model.ListExperimentRequest{
		ExperimentIDs: []*string{nil},
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
		Sort: &model.ExperimentSortInput{
			Field: "name",
		},
		Filter: &model.ExperimentFilterInput{},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}
	mockData := ExperimentListData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{},
		Data: ExperimentList{
			ListExperimentDetails: model.ListExperimentResponse{
				TotalNoOfExperiments: 1,
				Experiments: []*model.Experiment{
					{
						Name: "testexperiment",
					},
				},
			},
		},
	}

	mockResponse := mockData
	mockClient := &MockHTTPClientGetExperimentList{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, err := GetExperimentList(pid, requestData, cred, mockClient)

	// Assertions
	assert.NoError(t, err)                // Ensure no errors occurred
	assert.Equal(t, result, mockResponse) // Ensure the result matches the expected response
}

func TestGetExperimentListFailedRequest(t *testing.T) {

	pid := "testprojectid"
	requestData := model.ListExperimentRequest{
		ExperimentIDs: []*string{nil},
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
		Sort: &model.ExperimentSortInput{
			Field: "name",
		},
		Filter: &model.ExperimentFilterInput{},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientGetExperimentList{
		mockError:    errors.New("Some error occurred"),
		mockResponse: ExperimentListData{},
	}

	result, err := GetExperimentList(pid, requestData, cred, mockClient)

	// Assertions
	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, ExperimentListData{}, "Expected an empty response")
}

type MockHTTPClientGetExperimentRunsList struct {
	mockResponse ExperimentRunListData
	mockError    error
}

func (c *MockHTTPClientGetExperimentRunsList) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func TestGetExperimentRunsListSuccess(t *testing.T) {

	pid := "testprojectid"
	requestData := model.ListExperimentRunRequest{
		ExperimentRunIDs: []*string{nil},
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
		Sort: &model.ExperimentRunSortInput{
			Field: "name",
		},
		Filter: &model.ExperimentRunFilterInput{},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}
	mockData := ExperimentRunListData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{},
		Data: ExperimentRunsList{},
	}

	mockResponse := mockData
	mockClient := &MockHTTPClientGetExperimentRunsList{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, err := GetExperimentRunsList(pid, requestData, cred, mockClient)

	// Assertions
	assert.NoError(t, err)                // Ensure no errors occurred
	assert.Equal(t, result, mockResponse) // Ensure the result matches the expected response
}

type MockHTTPClientGetServerVersion struct {
	mockResponse ServerVersionResponse
	mockError    error
}

func (c *MockHTTPClientGetServerVersion) Do(req *http.Request) (*http.Response, error) {

	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}

func TestGetExperimentRunsListFailure(t *testing.T) {

	pid := "testprojectid"
	requestData := model.ListExperimentRunRequest{
		ExperimentRunIDs: []*string{nil},
		Pagination: &model.Pagination{
			Page:  1,
			Limit: 10,
		},
		Sort: &model.ExperimentRunSortInput{
			Field: "name",
		},
		Filter: &model.ExperimentRunFilterInput{},
	}
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientGetExperimentRunsList{
		mockError:    errors.New("Some error occurred"),
		mockResponse: ExperimentRunListData{},
	}

	result, err := GetExperimentRunsList(pid, requestData, cred, mockClient)

	// Assertions
	assert.Error(t, err, "Expected an error")
	assert.Equal(t, result, ExperimentRunListData{}, "Expected an empty response")
}

func TestGetServerVersionSuccess(t *testing.T) {
	endpoint := "https://example.com"
	mockResponse := ServerVersionResponse{
		Data: ServerVersionData{
			GetServerVersion: GetServerVersionData{
				Key:   "version",
				Value: "1.0",
			},
		},
		Errors: nil,
	}

	mockClient := &MockHTTPClientGetServerVersion{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, err := GetServerVersion(endpoint, mockClient)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, result, mockResponse)
}

// TestGetServerVersionFailedRequest tests the GetServerVersion function with a failed request.
func TestGetServerVersionFailedRequest(t *testing.T) {

	endpoint := "https://example.com"

	mockClient := &MockHTTPClientGetServerVersion{
		mockError:    errors.New("Some error occurred"),
		mockResponse: ServerVersionResponse{},
	}

	result, err := GetServerVersion(endpoint, mockClient)

	// Assertions
	assert.Error(t, err)                             // Ensure an error occurred
	assert.Equal(t, result, ServerVersionResponse{}) // Ensure the result is an empty response
}

type MockHTTPClientDeleteExperiment struct {
	mockResponse DeleteChaosExperimentData
	mockError    error
}

func (c *MockHTTPClientDeleteExperiment) Do(req *http.Request) (*http.Response, error) {
	if c.mockError != nil {
		return nil, c.mockError
	}
	data, err := json.Marshal(c.mockResponse)
	if err != nil {
		return nil, err
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewReader(data)),
	}, nil
}
func TestDeleteChaosExperimentSuccess(t *testing.T) {

	projectID := "testprojectid"
	experimentID := "test-experiment-id"
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockResponse := DeleteChaosExperimentData{
		Errors: nil, // No errors
		Data: DeleteChaosExperimentDetails{
			IsDeleted: true,
		},
	}

	mockClient := &MockHTTPClientDeleteExperiment{
		mockError:    nil,
		mockResponse: mockResponse,
	}

	result, err := DeleteChaosExperiment(projectID, &experimentID, cred, mockClient)

	// Assertions
	assert.NoError(t, err) // Ensure no errors occurred
	assert.Equal(t, result, mockResponse)
}

func TestDeleteChaosExperimentFailedRequest(t *testing.T) {

	projectID := "testprojectid"
	experimentID := "test-experiment-id"
	cred := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockResponse := DeleteChaosExperimentData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Some error occurred",
				Path:    []string{"path", "to", "error"},
			},
		},
		Data: DeleteChaosExperimentDetails{}, // Empty data
	}

	mockClient := &MockHTTPClientDeleteExperiment{
		mockError:    errors.New("Some error occurred"),
		mockResponse: mockResponse,
	}

	_, err := DeleteChaosExperiment(projectID, &experimentID, cred, mockClient)

	// Assertions
	assert.Error(t, err) // Ensure an error occurred
}
