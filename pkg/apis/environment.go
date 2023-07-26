package apis

import (
	"encoding/json"
	"errors"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"io/ioutil"
	"net/http"
)

type CreateEnvironmentGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectId string                          `json:"projectID"`
		Request   models.CreateEnvironmentRequest `json:"request"`
	} `json:"variables"`
}

type EnvironmentConnectionData struct {
	Data   EnvironmentData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type EnvironmentData struct {
	EnvironmentDetails models.Environment `json:"connectEnvironment"`
}

// ConnectEnvironment connects the  Infra with the given details
func ConnectEnvironment(pid string, request models.CreateEnvironmentRequest, cred types.Credentials) (EnvironmentConnectionData, error) {
	var gqlReq CreateEnvironmentGQLRequest
	gqlReq.Query = `mutation createEnvironment($projectID: ID!, $request: CreateEnvironmentRequest!) {
					  createEnvironment(
						projectID: $projectID
						request: $request
					  ) {
						environmentID
						name	
					  }
					}
					`
	gqlReq.Variables.ProjectId = pid
	gqlReq.Variables.Request = request

	query, err := json.Marshal(gqlReq)
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, query, string(types.Post))
	if err != nil {
		return EnvironmentConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return EnvironmentConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectEnvironment EnvironmentConnectionData
		err = json.Unmarshal(bodyBytes, &connectEnvironment)
		if err != nil {
			return EnvironmentConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
		}

		if len(connectEnvironment.Errors) > 0 {
			return EnvironmentConnectionData{}, errors.New(connectEnvironment.Errors[0].Message)
		}
		return connectEnvironment, nil
	} else {
		return EnvironmentConnectionData{}, err
	}
}
