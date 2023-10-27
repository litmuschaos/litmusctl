package infrastructure

import (
	"bytes"
	"encoding/json"
	"testing"

	"io/ioutil"
	"net/http"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/stretchr/testify/assert"
)

type MockHTTPClientInfraGet struct {
	mockResponse InfraData
	mockError    error
}

func (c *MockHTTPClientInfraGet) Do(req *http.Request) (*http.Response, error) {
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

func TestGetInfraListSuccess(t *testing.T) {
	mockData := InfraData{
		Data: InfraList{
			ListInfraDetails: models.ListInfraResponse{
				TotalNoOfInfras: 2,
				Infras: []*models.Infra{
					{
						ProjectID:               "project1",
						InfraID:                 "infra1",
						Name:                    "Infrastructure 1",
						Description:             nil,
						Tags:                    []string{"tag1", "tag2"},
						EnvironmentID:           "env1",
						PlatformName:            "GKE",
						IsActive:                true,
						IsInfraConfirmed:        true,
						IsRemoved:               false,
						UpdatedAt:               "2023-10-13",
						CreatedAt:               "2023-10-12",
						NoOfExperiments:         nil,
						NoOfExperimentRuns:      nil,
						Token:                   "token1",
						InfraNamespace:          nil,
						ServiceAccount:          nil,
						InfraScope:              "ns",
						InfraNsExists:           nil,
						InfraSaExists:           nil,
						LastExperimentTimestamp: nil,
						StartTime:               "2023-10-11",
						Version:                 "v1",
						CreatedBy: &models.UserDetails{
							UserID:   "user1",
							Username: "JohnDoe",
							Email:    "johndoe@example.com",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user2",
							Username: "JaneSmith",
							Email:    "janesmith@example.com",
						},
						InfraType:    nil,
						UpdateStatus: "In Progress",
					},
					{
						ProjectID:               "project2",
						InfraID:                 "infra2",
						Name:                    "Infrastructure 2",
						Description:             nil,
						Tags:                    []string{"tag3", "tag4"},
						EnvironmentID:           "env2",
						PlatformName:            "AWS",
						IsActive:                false,
						IsInfraConfirmed:        true,
						IsRemoved:               true,
						UpdatedAt:               "2023-10-15",
						CreatedAt:               "2023-10-14",
						NoOfExperiments:         nil,
						NoOfExperimentRuns:      nil,
						Token:                   "token2",
						InfraNamespace:          nil,
						ServiceAccount:          nil,
						InfraScope:              "cluster",
						InfraNsExists:           nil,
						InfraSaExists:           nil,
						LastExperimentTimestamp: nil,
						StartTime:               "2023-10-16",
						Version:                 "v2",
						CreatedBy: &models.UserDetails{
							UserID:   "user3",
							Username: "AliceJohnson",
							Email:    "alicejohnson@example.com",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user4",
							Username: "BobWilliams",
							Email:    "bobwilliams@example.com",
						},
						InfraType:    nil,
						UpdateStatus: "Completed",
					},
				},
			},
		},
		Errors: nil, // We can add errors if needed
	}

	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientInfraGet{
		mockError:    nil,
		mockResponse: mockData,
	}

	tag1 := "tag1"
	tag2 := "tag2"

	mockRequest := models.ListInfraRequest{
		InfraIDs:       []string{"infra1", "infra2"},
		EnvironmentIDs: []string{"env1", "env2"},
		Pagination: &models.Pagination{
			Page:  1,
			Limit: 10,
		},
		Filter: &models.InfraFilterInput{
			Name:         nil,
			InfraID:      nil,
			Description:  nil,
			PlatformName: nil,
			InfraScope:   nil,
			IsActive:     nil,
			Tags:         []*string{&tag1, &tag2},
		},
	}

	result, _ := GetInfraList(input, "testpid", mockRequest, mockClient)
	assert.Equal(t, result, mockData)

}

func TestGetInfraListFailed(t *testing.T) {
	mockData := InfraData{
		Data: InfraList{
			ListInfraDetails: models.ListInfraResponse{
				TotalNoOfInfras: 2,
				Infras: []*models.Infra{
					{
						ProjectID:               "project1",
						InfraID:                 "infra1",
						Name:                    "Infrastructure 1",
						Description:             nil,
						Tags:                    []string{"tag1", "tag2"},
						EnvironmentID:           "env1",
						PlatformName:            "GKE",
						IsActive:                true,
						IsInfraConfirmed:        true,
						IsRemoved:               false,
						UpdatedAt:               "2023-10-13",
						CreatedAt:               "2023-10-12",
						NoOfExperiments:         nil,
						NoOfExperimentRuns:      nil,
						Token:                   "token1",
						InfraNamespace:          nil,
						ServiceAccount:          nil,
						InfraScope:              "ns",
						InfraNsExists:           nil,
						InfraSaExists:           nil,
						LastExperimentTimestamp: nil,
						StartTime:               "2023-10-11",
						Version:                 "v1",
						CreatedBy: &models.UserDetails{
							UserID:   "user1",
							Username: "JohnDoe",
							Email:    "johndoe@example.com",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user2",
							Username: "JaneSmith",
							Email:    "janesmith@example.com",
						},
						InfraType:    nil,
						UpdateStatus: "In Progress",
					},
					{
						ProjectID:               "project2",
						InfraID:                 "infra2",
						Name:                    "Infrastructure 2",
						Description:             nil,
						Tags:                    []string{"tag3", "tag4"},
						EnvironmentID:           "env2",
						PlatformName:            "AWS",
						IsActive:                false,
						IsInfraConfirmed:        true,
						IsRemoved:               true,
						UpdatedAt:               "2023-10-15",
						CreatedAt:               "2023-10-14",
						NoOfExperiments:         nil,
						NoOfExperimentRuns:      nil,
						Token:                   "token2",
						InfraNamespace:          nil,
						ServiceAccount:          nil,
						InfraScope:              "cluster",
						InfraNsExists:           nil,
						InfraSaExists:           nil,
						LastExperimentTimestamp: nil,
						StartTime:               "2023-10-16",
						Version:                 "v2",
						CreatedBy: &models.UserDetails{
							UserID:   "user3",
							Username: "AliceJohnson",
							Email:    "alicejohnson@example.com",
						},
						UpdatedBy: &models.UserDetails{
							UserID:   "user4",
							Username: "BobWilliams",
							Email:    "bobwilliams@example.com",
						},
						InfraType:    nil,
						UpdateStatus: "Completed",
					},
				},
			},
		},
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Error message 1",
				Path:    []string{"path1", "path2"},
			},
		},
	}

	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockClient := &MockHTTPClientInfraGet{
		mockError:    nil,
		mockResponse: mockData,
	}

	tag1 := "tag1"
	tag2 := "tag2"

	mockRequest := models.ListInfraRequest{
		InfraIDs:       []string{"infra1", "infra2"},
		EnvironmentIDs: []string{"env1", "env2"},
		Pagination: &models.Pagination{
			Page:  1,
			Limit: 10,
		},
		Filter: &models.InfraFilterInput{
			Name:         nil,
			InfraID:      nil,
			Description:  nil,
			PlatformName: nil,
			InfraScope:   nil,
			IsActive:     nil,
			Tags:         []*string{&tag1, &tag2},
		},
	}

	result, _ := GetInfraList(input, "testpid", mockRequest, mockClient)
	assert.NotEqual(t, result, mockData)
}

type MockHTTPClientConnect struct {
	mockResponse InfraConnectionData
	mockError    error
}

func (c *MockHTTPClientConnect) Do(req *http.Request) (*http.Response, error) {
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

func TestConnectInfraSuccess(t *testing.T) {

	mockRequest := types.Infra{
		InfraName:      "Infraname",
		Mode:           "mode",
		Description:    "desc",
		PlatformName:   "platform",
		EnvironmentID:  "estid",
		ProjectId:      "projectid",
		InfraType:      "infratype",
		NodeSelector:   "nodeselector",
		Tolerations:    "",
		Namespace:      "namespace",
		ServiceAccount: "serviceaccount",
		NsExists:       false,
		SAExists:       false,
		SkipSSL:        false,
	}
	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockData := InfraConnectionData{
		Data: RegisterInfra{
			RegisterInfraDetails: models.RegisterInfraResponse{
				Token:    "dummyToken",
				InfraID:  "dummyInfraID",
				Name:     "DummyInfraName",
				Manifest: "DummyInfraManifest",
			},
		},
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{},
	}
	mockClient := &MockHTTPClientConnect{
		mockError:    nil,
		mockResponse: mockData,
	}
	result, _ := ConnectInfra(mockRequest, input, mockClient)
	assert.Equal(t, result, mockData)
}

func TestConnectInfraFail(t *testing.T) {

	mockRequest := types.Infra{
		InfraName:      "Infraname",
		Mode:           "mode",
		Description:    "desc",
		PlatformName:   "platform",
		EnvironmentID:  "estid",
		ProjectId:      "projectid",
		InfraType:      "infratype",
		NodeSelector:   "nodeselector",
		Tolerations:    "",
		Namespace:      "namespace",
		ServiceAccount: "serviceaccount",
		NsExists:       false,
		SAExists:       false,
		SkipSSL:        false,
	}
	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockData := InfraConnectionData{
		Data: RegisterInfra{
			RegisterInfraDetails: models.RegisterInfraResponse{
				Token:    "dummyToken",
				InfraID:  "dummyInfraID",
				Name:     "DummyInfraName",
				Manifest: "DummyInfraManifest",
			},
		},
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Error message 1",
				Path:    []string{"path1", "path2"},
			},
		},
	}
	mockClient := &MockHTTPClientConnect{
		mockError:    nil,
		mockResponse: mockData,
	}
	result, _ := ConnectInfra(mockRequest, input, mockClient)
	assert.NotEqual(t, result, mockData)
}
func TestCreateRegisterInfraRequestSuccess(t *testing.T) {

	mockRequest := types.Infra{
		InfraName:      "Infraname",
		Mode:           "mode",
		Description:    "desc",
		PlatformName:   "platform",
		EnvironmentID:  "estid",
		ProjectId:      "projectid",
		InfraType:      "infratype",
		NodeSelector:   "nodeselector",
		Tolerations:    "",
		Namespace:      "namespace",
		ServiceAccount: "serviceaccount",
		NsExists:       false,
		SAExists:       false,
		SkipSSL:        false,
	}

	mockData := models.RegisterInfraRequest{
		Name:               "Infraname",
		EnvironmentID:      "estid",
		InfrastructureType: "INTERNAL",
		Description:        stringPointer("desc"),
		PlatformName:       "platform",
		InfraNamespace:     stringPointer("namespace"),
		ServiceAccount:     stringPointer("serviceaccount"),
		InfraScope:         "mode",
		InfraNsExists:      boolPointer(false),
		InfraSaExists:      boolPointer(false),
		SkipSsl:            boolPointer(false),
		NodeSelector:       nil,
		Tolerations:        []*models.Toleration(nil),
		Tags:               []string(nil),
	}

	result := CreateRegisterInfraRequest(mockRequest)
	assert.Equal(t, result, mockData)
}

func TestCreateRegisterInfraRequestFail(t *testing.T) {

	mockRequest := types.Infra{
		InfraName:      "Infraname",
		Mode:           "mode",
		Description:    "desc",
		PlatformName:   "platform",
		EnvironmentID:  "estid",
		ProjectId:      "projectid",
		InfraType:      "infratype",
		NodeSelector:   "nodeselector",
		Tolerations:    "",
		Namespace:      "namespace",
		ServiceAccount: "serviceaccount",
		NsExists:       false,
		SAExists:       false,
		SkipSSL:        false,
	}

	mockData := models.RegisterInfraRequest{
		Name:               "Infraname",
		EnvironmentID:      "estid",
		InfrastructureType: "INTERNAL",
		Description:        stringPointer("desc"),
		PlatformName:       "platform",
		InfraNamespace:     stringPointer("namespace"),
		ServiceAccount:     stringPointer("serviceaccount"),
		InfraScope:         "mode",
		InfraNsExists:      boolPointer(false),
		InfraSaExists:      boolPointer(false),
		SkipSsl:            boolPointer(false),
		NodeSelector:       nil,
		Tolerations:        []*models.Toleration(nil),
		Tags:               []string{"tag1", "tag2"},
	}

	result := CreateRegisterInfraRequest(mockRequest)
	assert.NotEqual(t, result, mockData)
}
func stringPointer(s string) *string {
	return &s
}

func boolPointer(b bool) *bool {
	return &b
}

type MockHTTPClientDisconnect struct {
	mockResponse DisconnectInfraData
	mockError    error
}

func (c *MockHTTPClientDisconnect) Do(req *http.Request) (*http.Response, error) {
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

func TestDisconnectInfraSuccess(t *testing.T) {

	projectID := "testprojectid"
	infraID := "infratestid"

	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockData := DisconnectInfraData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			// {
			// 	Message: "Error message 1",
			// 	Path:    []string{"path", "to", "error1"},
			// },
			// {
			// 	Message: "Error message 2",
			// 	Path:    []string{"path", "to", "error2"},
			// },
		},
		Data: DisconnectInfraDetails{
			Message: "Infra deleted successfully",
		},
	}

	mockClient := &MockHTTPClientDisconnect{
		mockError:    nil,
		mockResponse: mockData,
	}
	result, _ := DisconnectInfra(projectID, infraID, input, mockClient)
	assert.Equal(t, result, mockData)
}

func TestDisconnectInfraFail(t *testing.T) {

	projectID := "testprojectid"
	infraID := "infratestid"

	input := types.Credentials{
		Username: "testusername",
		Token:    "testtoken",
		Endpoint: "https://example.com",
	}

	mockData := DisconnectInfraData{
		Errors: []struct {
			Message string   `json:"message"`
			Path    []string `json:"path"`
		}{
			{
				Message: "Error message 1",
				Path:    []string{"path", "to", "error1"},
			},
			{
				Message: "Error message 2",
				Path:    []string{"path", "to", "error2"},
			},
		},
		Data: DisconnectInfraDetails{
			Message: "Infra deleted successfully",
		},
	}

	mockClient := &MockHTTPClientDisconnect{
		mockError:    nil,
		mockResponse: mockData,
	}
	result, _ := DisconnectInfra(projectID, infraID, input, mockClient)
	assert.NotEqual(t, result, mockData)
}
