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

	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

type CreateWorkflowData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data WorkflowData `json:"data"`
}

type WorkflowData struct {
	CreateChaosWorkflow CreateChaosWorkflow `json:"createChaosWorkFlow"`
}

type CreateChaosWorkflow struct {
	WorkflowID          string `json:"workflow_id"`
	CronSyntax          string `json:"cronSyntax"`
	WorkflowName        string `json:"workflow_name"`
	WorkflowDescription string `json:"workflow_description"`
	IsCustomWorkflow    bool   `json:"isCustomWorkflow"`
}

// CreateWorkflow sends GraphQL API request for creating a workflow
func CreateWorkflow(in types.ChaosWorkFlowInput, cred types.Credentials) (CreateWorkflowData, error) {

	var gqlReq types.ChaosWorkFlowGraphQLRequest

	gqlReq.Query = `mutation createChaosWorkFlow($ChaosWorkFlowInput: ChaosWorkFlowInput!) {
                      createChaosWorkFlow(input: $ChaosWorkFlowInput) {
                        workflow_id
                        cronSyntax
                        workflow_name
                        workflow_description
                        isCustomWorkflow
                      }
                    }`
	gqlReq.Variables.ChaosWorkFlowInput = in

	query, _ := json.Marshal(gqlReq)

	resp, _ := SendRequest(
		SendRequestParams{
			Endpoint: cred.Endpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return CreateWorkflowData{}, errors.New("Error in creating workflow: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var createdWorkflow CreateWorkflowData

		err = json.Unmarshal(bodyBytes, &createdWorkflow)
		if err != nil {
			return CreateWorkflowData{}, errors.New("Error in creating workflow: " + err.Error())
		}

		// Errors present
		if len(createdWorkflow.Errors) > 0 {
			return CreateWorkflowData{}, errors.New(createdWorkflow.Errors[0].Message)
		}

		return createdWorkflow, nil
	} else {
		return CreateWorkflowData{}, err
	}
}
