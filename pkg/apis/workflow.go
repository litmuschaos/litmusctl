/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package apis

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

type ChaosWorkflowCreationData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data CreatedChaosWorkflow `json:"data"`
}

type ServerVersionResponse struct {
	Data   ServerVersionData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type ServerVersionData struct {
	GetServerVersion GetServerVersionData `json:"getServerVersion"`
}

type GetServerVersionData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type CreatedChaosWorkflow struct {
	CreateChaosWorkflow model.ChaosWorkFlowResponse `json:"createChaosWorkFlow"`
}

type CreateChaosWorkFlowGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		CreateChaosWorkFlowRequest model.ChaosWorkFlowRequest `json:"request"`
	} `json:"variables"`
}

// CreateWorkflow sends GraphQL API request for creating a workflow
func CreateWorkflow(requestData model.ChaosWorkFlowRequest, cred types.Credentials) (ChaosWorkflowCreationData, error) {

	var gqlReq CreateChaosWorkFlowGraphQLRequest

	gqlReq.Query = `mutation createChaosWorkFlow($request: ChaosWorkFlowRequest!) {
                      createChaosWorkFlow(request: $request) {
                        workflowID
                        cronSyntax
                        workflowName
                        workflowDescription
                        isCustomWorkflow
                      }
                    }`
	gqlReq.Variables.CreateChaosWorkFlowRequest = requestData

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ChaosWorkflowCreationData{}, err
	}

	resp, err := SendRequest(
		SendRequestParams{
			Endpoint: cred.Endpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return ChaosWorkflowCreationData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ChaosWorkflowCreationData{}, errors.New("Error in creating workflow: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var createdWorkflow ChaosWorkflowCreationData

		err = json.Unmarshal(bodyBytes, &createdWorkflow)
		if err != nil {
			return ChaosWorkflowCreationData{}, errors.New("Error in creating workflow: " + err.Error())
		}

		// Errors present
		if len(createdWorkflow.Errors) > 0 {
			return ChaosWorkflowCreationData{}, errors.New(createdWorkflow.Errors[0].Message)
		}

		return createdWorkflow, nil
	} else {
		return ChaosWorkflowCreationData{}, err
	}
}

// GetServerVersion lists the agent connected to the specified project
func GetServerVersion() (ServerVersionResponse, error) {
	query := `{"query":"query{\n  getServerVersion{\n  key value\n  }\n}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: "http://localhost:8080/query"}, []byte(query), string(types.Post))
	if err != nil {
		return ServerVersionResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ServerVersionResponse{}, err
	}
	if resp.StatusCode == http.StatusOK {
		var version ServerVersionResponse
		err = json.Unmarshal(bodyBytes, &version)
		if err != nil {
			return ServerVersionResponse{}, err
		}

		if len(version.Errors) > 0 {
			return ServerVersionResponse{}, errors.New(version.Errors[0].Message)
		}

		return version, nil
	} else {
		return ServerVersionResponse{}, err
	}
}
