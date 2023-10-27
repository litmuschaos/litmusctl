package environment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClientCreate struct {
	mockError    error
	mockResponse CreateEnvironmentResponse
}

func (c *MockHTTPClientCreate) Do(req *http.Request) (*http.Response, error) {
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

type MockHTTPClientGet struct {
	mockResponse ListEnvironmentData
	mockError    error
}

func (c *MockHTTPClientGet) Do(req *http.Request) (*http.Response, error) {
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

func TestCreateEnvironmentSuccess(t *testing.T) {
	//define the mock request
	mockRequest := models.CreateEnvironmentRequest{
		EnvironmentID: "env123",
		Name:          "Test Environment",
		Type:          "Development",
		Description:   nil,
	}
	//custom input
	input := types.Credentials{
		Username: "testuser",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	//define the mock data response

	mockData := CreateEnvironmentResponse{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}(nil),
		Data: CreateEnvironmentData{
			EnvironmentDetails: models.Environment{
				ProjectID:     "project123",
				EnvironmentID: "env123",
				Name:          "Test Environment",
				Description:   nil,
				Tags:          []string{"tag1", "tag2"},
				Type:          "Development",
				CreatedAt:     "2023-10-13",
				CreatedBy: &models.UserDetails{
					UserID:   "user123",
					Username: "JohnDoe",
					Email:    "johndoe@example.com",
				},
				UpdatedBy: &models.UserDetails{
					UserID:   "user456",
					Username: "JaneSmith",
					Email:    "janesmith@example.com",
				},
				UpdatedAt: "2023-10-14",
				IsRemoved: nil,
				InfraIDs:  []string{"infra1", "infra2"},
			},
		},
	}

	mockClient := &MockHTTPClientCreate{
		mockResponse: mockData,
		mockError:    nil, // Set this to nil if you don't want to return an error.
	}

	result, err := CreateEnvironment("testpid", mockRequest, input, mockClient)

	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, mockData, result)
}

func TestCreateEnvironmentFailed(t *testing.T) {
	//define the mock request
	mockRequest := models.CreateEnvironmentRequest{
		EnvironmentID: "env123",
		Name:          "Test Environment",
		Type:          "Development",
		Description:   nil,
	}
	//custom input
	input := types.Credentials{
		Username: "testuser",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	//define the mock data response

	mockData := CreateEnvironmentResponse{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Error message 1",
				Path:    []string{"path1", "path2"},
			},
		},
		Data: CreateEnvironmentData{
			EnvironmentDetails: models.Environment{
				ProjectID:     "project123",
				EnvironmentID: "env123",
				Name:          "Test Environment",
				Description:   nil,
				Tags:          []string{"tag1", "tag2"},
				Type:          "Development",
				CreatedAt:     "2023-10-13",
				CreatedBy: &models.UserDetails{
					UserID:   "user123",
					Username: "JohnDoe",
					Email:    "johndoe@example.com",
				},
				UpdatedBy: &models.UserDetails{
					UserID:   "user456",
					Username: "JaneSmith",
					Email:    "janesmith@example.com",
				},
				UpdatedAt: "2023-10-14",
				IsRemoved: nil,
				InfraIDs:  []string{"infra1", "infra2"},
			},
		},
	}

	mockClient := &MockHTTPClientCreate{
		mockResponse: mockData,
		mockError:    nil,
	}

	result, _ := CreateEnvironment("testpid", mockRequest, input, mockClient)

	assert.NotEqual(t, mockData, result)
}

func TestGetEnvironmentListSuccess(t *testing.T) {

	mockData := ListEnvironmentData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			//keeping the error empty for asserting
		},
		Data: EnvironmentsList{
			ListEnvironmentDetails: models.ListEnvironmentResponse{
				TotalNoOfEnvironments: 2,
				Environments: []*models.Environment{
					{
						ProjectID:     "project1",
						EnvironmentID: "env1",
						Name:          "Environment 1",
						Description:   nil,
						Tags:          []string{"tag1", "tag2"},
						Type:          "Development",
						CreatedAt:     "2023-10-13",
						CreatedBy: &models.UserDetails{
							UserID:   "user1",
							Username: "JohnDoe",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user2",
							Username: "JaneSmith",
						},
						UpdatedAt: "2023-10-14",
						IsRemoved: nil,
						InfraIDs:  []string{"infra1", "infra2"},
					},
					{
						ProjectID:     "project2",
						EnvironmentID: "env2",
						Name:          "Environment 2",
						Description:   nil,
						Tags:          []string{"tag3", "tag4"},
						Type:          "Production",
						CreatedAt:     "2023-10-15",
						CreatedBy: &models.UserDetails{
							UserID:   "user3",
							Username: "AliceJohnson",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user4",
							Username: "BobWilliams",
						},
						UpdatedAt: "2023-10-16",
						IsRemoved: nil,
						InfraIDs:  []string{"infra3", "infra4"},
					},
				},
			},
		},
	}

	mockClient := &MockHTTPClientGet{
		mockResponse: mockData,
		mockError:    nil,
	}

	input := types.Credentials{
		Username: "testuser",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	result, err := GetEnvironmentList("testpid", input, mockClient)

	if err != nil {
		fmt.Println(err)
	}
	assert.Equal(t, mockData, result)
}

func TestGetEnvironmentListFailed(t *testing.T) {

	mockData := ListEnvironmentData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Error message 1",
				Path:    []string{"path1", "path2"},
			},
		},
		Data: EnvironmentsList{
			ListEnvironmentDetails: models.ListEnvironmentResponse{
				TotalNoOfEnvironments: 2,
				Environments: []*models.Environment{
					{
						ProjectID:     "project1",
						EnvironmentID: "env1",
						Name:          "Environment 1",
						Description:   nil,
						Tags:          []string{"tag1", "tag2"},
						Type:          "Development",
						CreatedAt:     "2023-10-13",
						CreatedBy: &models.UserDetails{
							UserID:   "user1",
							Username: "JohnDoe",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user2",
							Username: "JaneSmith",
						},
						UpdatedAt: "2023-10-14",
						IsRemoved: nil,
						InfraIDs:  []string{"infra1", "infra2"},
					},
					{
						ProjectID:     "project2",
						EnvironmentID: "env2",
						Name:          "Environment 2",
						Description:   nil,
						Tags:          []string{"tag3", "tag4"},
						Type:          "Production",
						CreatedAt:     "2023-10-15",
						CreatedBy: &models.UserDetails{
							UserID:   "user3",
							Username: "AliceJohnson",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user4",
							Username: "BobWilliams",
						},
						UpdatedAt: "2023-10-16",
						IsRemoved: nil,
						InfraIDs:  []string{"infra3", "infra4"},
					},
				},
			},
		},
	}

	mockClient := &MockHTTPClientGet{
		mockResponse: mockData,
		mockError:    nil,
	}

	input := types.Credentials{
		Username: "testuser",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	result, _ := GetEnvironmentList("testpid", input, mockClient)

	assert.NotEqual(t, mockData, result)
}
