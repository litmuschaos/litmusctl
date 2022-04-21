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

type WorkflowCreationData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data CreatedWorkflow `json:"data"`
}

type CreatedWorkflow struct {
	CreateChaosWorkflow model.ChaosWorkFlowResponse `json:"createChaosWorkFlow"`
}

type CreateChaosWorkFlowGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		CreateChaosWorkFlowInput model.ChaosWorkFlowInput `json:"ChaosWorkFlowInput"`
	} `json:"variables"`
}

// CreateWorkflow sends GraphQL API request for creating a workflow
func CreateWorkflow(in model.ChaosWorkFlowInput, cred types.Credentials) (WorkflowCreationData, error) {

	var gqlReq CreateChaosWorkFlowGraphQLRequest

	gqlReq.Query = `mutation createChaosWorkFlow($ChaosWorkFlowInput: ChaosWorkFlowInput!) {
                      createChaosWorkFlow(input: $ChaosWorkFlowInput) {
                        workflow_id
                        cronSyntax
                        workflow_name
                        workflow_description
                        isCustomWorkflow
                      }
                    }`
	gqlReq.Variables.CreateChaosWorkFlowInput = in

	query, _ := json.Marshal(gqlReq)

	resp, err := SendRequest(
		SendRequestParams{
			Endpoint: cred.Endpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)
	if err != nil {
		return WorkflowCreationData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return WorkflowCreationData{}, errors.New("Error in creating workflow: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var createdWorkflow WorkflowCreationData

		err = json.Unmarshal(bodyBytes, &createdWorkflow)
		if err != nil {
			return WorkflowCreationData{}, errors.New("Error in creating workflow: " + err.Error())
		}

		// Errors present
		if len(createdWorkflow.Errors) > 0 {
			return WorkflowCreationData{}, errors.New(createdWorkflow.Errors[0].Message)
		}

		return createdWorkflow, nil
	} else {
		return WorkflowCreationData{}, err
	}
}
