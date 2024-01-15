package environment

import (
	"encoding/json"
	"errors"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"io/ioutil"
	"net/http"
)

// CreateEnvironment connects the  Infra with the given details
func CreateEnvironment(pid string, request models.CreateEnvironmentRequest, cred types.Credentials) (CreateEnvironmentResponse, error) {
	var gqlReq CreateEnvironmentGQLRequest
	gqlReq.Query = CreateEnvironmentQuery
	gqlReq.Variables.ProjectId = pid
	gqlReq.Variables.Request = request

	query, err := json.Marshal(gqlReq)
	resp, err := apis.SendRequest(apis.SendRequestParams{Endpoint: cred.ServerEndpoint + utils.GQLAPIPath, Token: cred.Token}, query, string(types.Post))
	if err != nil {
		return CreateEnvironmentResponse{}, errors.New("Error in Creating Chaos Infrastructure: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return CreateEnvironmentResponse{}, errors.New("Error in Creating Chaos Environment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectEnvironment CreateEnvironmentResponse
		err = json.Unmarshal(bodyBytes, &connectEnvironment)
		if err != nil {
			return CreateEnvironmentResponse{}, errors.New("Error in Creating Chaos Environment: " + err.Error())
		}

		if len(connectEnvironment.Errors) > 0 {
			return CreateEnvironmentResponse{}, errors.New(connectEnvironment.Errors[0].Message)
		}
		return connectEnvironment, nil
	} else {
		return CreateEnvironmentResponse{}, err
	}
}

func GetEnvironmentList(pid string, cred types.Credentials) (ListEnvironmentData, error) {
	var err error
	var gqlReq CreateEnvironmentListGQLRequest
	gqlReq.Query = ListEnvironmentQuery

	gqlReq.Variables.Request = models.ListEnvironmentRequest{}
	gqlReq.Variables.ProjectID = pid
	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ListEnvironmentData{}, err
	}
	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return ListEnvironmentData{}, errors.New("Error in Getting Chaos Environment List: " + err.Error())
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ListEnvironmentData{}, errors.New("Error in Getting Chaos Environment List: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var listEnvironment ListEnvironmentData
		err = json.Unmarshal(bodyBytes, &listEnvironment)
		if err != nil {
			return ListEnvironmentData{}, errors.New("Error in Getting Chaos Environment List: " + err.Error())
		}
		if len(listEnvironment.Errors) > 0 {
			return ListEnvironmentData{}, errors.New(listEnvironment.Errors[0].Message)
		}
		return listEnvironment, nil
	} else {
		return ListEnvironmentData{}, err
	}
}
