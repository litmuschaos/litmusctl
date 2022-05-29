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

type WorkflowListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data WorkflowList `json:"data"`
}

type WorkflowList struct {
	ListWorkflowDetails model.ListWorkflowsResponse `json:"listWorkflows"`
}

type GetChaosWorkFlowsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosWorkFlowsRequest model.ListWorkflowsRequest `json:"request"`
	} `json:"variables"`
}

// GetWorkflowList sends GraphQL API request for fetching a list of workflows.
func GetWorkflowList(in model.ListWorkflowsRequest, cred types.Credentials) (WorkflowListData, error) {

	var gqlReq GetChaosWorkFlowsGraphQLRequest
	var err error

	gqlReq.Query = `query listWorkflows($request: ListWorkflowsRequest!) {
                      listWorkflows(request: $request) {
                        totalNoOfWorkflows
                        workflows {
                          workflowID
                          workflowManifest
                          cronSyntax
                          clusterName
                          workflowName
                          workflowDescription
                          weightages {
                            experimentName
                            weightage
                          }
                          isCustomWorkflow
                          updatedAt
                          createdAt
                          projectID
                          clusterID
                          clusterType
                          isRemoved
                          lastUpdatedBy
                        }
                      }
                    }`
	gqlReq.Variables.GetChaosWorkFlowsRequest = in

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return WorkflowListData{}, err
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
		return WorkflowListData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return WorkflowListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var workflowList WorkflowListData
		err = json.Unmarshal(bodyBytes, &workflowList)
		if err != nil {
			return WorkflowListData{}, err
		}

		if len(workflowList.Errors) > 0 {
			return WorkflowListData{}, errors.New(workflowList.Errors[0].Message)
		}

		return workflowList, nil
	} else {
		return WorkflowListData{}, err
	}
}

type WorkflowRunsListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data WorkflowRunsList `json:"data"`
}

type WorkflowRunsList struct {
	ListWorkflowRunsDetails model.ListWorkflowRunsResponse `json:"listWorkflowRuns"`
}

type GetChaosWorkFlowRunsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosWorkFlowRunsRequest model.ListWorkflowRunsRequest `json:"request"`
	} `json:"variables"`
}

// GetWorkflowRunsList sends GraphQL API request for fetching a list of workflow runs.
func GetWorkflowRunsList(in model.ListWorkflowRunsRequest, cred types.Credentials) (WorkflowRunsListData, error) {

	var gqlReq GetChaosWorkFlowRunsGraphQLRequest
	var err error

	gqlReq.Query = `query listWorkflowRuns($request: ListWorkflowRunsRequest!) {
                      listWorkflowRuns(request: $request) {
                        totalNoOfWorkflowRuns
                        workflowRuns {
                          workflowRunID
                          workflowID
                          clusterName
                          workflowName
                          projectID
                          clusterID
                          clusterType
                          isRemoved
                          lastUpdated
                          phase
                          resiliencyScore
                          experimentsPassed
                          experimentsFailed
                          experimentsAwaited
                          experimentsStopped
                          experimentsNa
                          totalExperiments
                          executedBy
                        }
                      }
                    }`
	gqlReq.Variables.GetChaosWorkFlowRunsRequest = in

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return WorkflowRunsListData{}, err
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
		return WorkflowRunsListData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return WorkflowRunsListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var workflowRunsList WorkflowRunsListData
		err = json.Unmarshal(bodyBytes, &workflowRunsList)
		if err != nil {
			return WorkflowRunsListData{}, err
		}

		if len(workflowRunsList.Errors) > 0 {
			return WorkflowRunsListData{}, errors.New(workflowRunsList.Errors[0].Message)
		}

		return workflowRunsList, nil
	} else {
		return WorkflowRunsListData{}, err
	}
}
