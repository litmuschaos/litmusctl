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

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
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
		return ChaosWorkflowCreationData{}, errors.New("Error in creating Chaos Scenario: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var createdWorkflow ChaosWorkflowCreationData

		err = json.Unmarshal(bodyBytes, &createdWorkflow)

		if err != nil {
			return ChaosWorkflowCreationData{}, errors.New("Error in creating Chaos Scenario: " + err.Error())
		}

		// Errors present
		if len(createdWorkflow.Errors) > 0 {
			return ChaosWorkflowCreationData{}, errors.New(createdWorkflow.Errors[0].Message)
		}

		return createdWorkflow, nil
	} else {
		return ChaosWorkflowCreationData{}, errors.New("graphql schema error")
	}
}

type ExperimentListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentList `json:"data"`
}

type ExperimentList struct {
	ListExperimentDetails models.ListExperimentResponse `json:"listExperiment"`
}

type GetChaosExperimentsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosExperimentRequest models.ListExperimentRequest `json:"request"`
		ProjectID                 string                       `json:"projectID"`
	} `json:"variables"`
}

// GetExperimentList sends GraphQL API request for fetching a list of experiments.
func GetExperimentList(pid string, in models.ListExperimentRequest, cred types.Credentials) (ExperimentListData, error) {

	var gqlReq GetChaosExperimentsGraphQLRequest
	var err error

	gqlReq.Query = `query listExperiment($projectID: String!, $request: ListExperimentRequest!) {
                      listExperiment(project: $projectID, request: $request) {
                        totalNoOfExperiments
                        experiments {
                          experimentID
                          experimentManifest
                          cronSyntax
                          name
                          description
                          weightages {
                            faultName
                            weightage
                          }
                          isCustomExperiment
                          updatedAt
                          createdAt
                          infra {
                            projectID
                            name
                            infraID
                            infraType
                          }
                          isRemoved
                          updatedBy
                        }
                      }
                    }`
	gqlReq.Variables.GetChaosExperimentRequest = in
	gqlReq.Variables.ProjectID = pid

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ExperimentListData{}, err
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
		return ExperimentListData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ExperimentListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var experimentList ExperimentListData
		err = json.Unmarshal(bodyBytes, &experimentList)
		if err != nil {
			return ExperimentListData{}, err
		}

		if len(experimentList.Errors) > 0 {
			return ExperimentListData{}, errors.New(experimentList.Errors[0].Message)
		}

		return experimentList, nil
	} else {
		return ExperimentListData{}, errors.New("Error while fetching the Chaos Scenarios")
	}
}

type ExperimentRunListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentRunsList `json:"data"`
}

type ExperimentRunsList struct {
	ListExperimentRunDetails model.ListWorkflowRunsResponse `json:"listExperimentRun"`
}

type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                          `json:"projectID"`
		GetChaosExperimentRunRequest models.ListExperimentRunRequest `json:"request"`
	} `json:"variables"`
}

// GetExperimentRunsList sends GraphQL API request for fetching a list of workflow runs.
func GetExperimentRunsList(pid string, in models.ListExperimentRunRequest, cred types.Credentials) (ExperimentRunListData, error) {

	var gqlReq GetChaosExperimentRunGraphQLRequest
	var err error

	gqlReq.Query = `query listExperimentRuns($projectID: String!, $request: ListExperimentRunRequest!) {
                      listWorkflowRuns(projectID: $projectID, request: $request) {
                        totalNoOfExperimentsRuns
                        experimentRuns {
                          experimentRunID
                          experimentID
                          experimentName
                          infra {
                          name
                          projectID
                          infraID
                          infraType
                          }
                          isRemoved
                          updatedAt
                          phase
                          resiliencyScore
                          faultsPassed
                          faultsFailed
                          faultsAwaited
                          faultsStopped
                          faultsNa
                          totalFaults
                          executedBy
                        }
                      }
                    }`
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.GetChaosExperimentRunRequest = in

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ExperimentRunListData{}, err
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
		return ExperimentRunListData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ExperimentRunListData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var workflowRunsList ExperimentRunListData
		err = json.Unmarshal(bodyBytes, &workflowRunsList)
		if err != nil {
			return ExperimentRunListData{}, err
		}

		if len(workflowRunsList.Errors) > 0 {
			return ExperimentRunListData{}, errors.New(workflowRunsList.Errors[0].Message)
		}

		return workflowRunsList, nil
	} else {
		return ExperimentRunListData{}, errors.New("Error while fetching the Chaos Scenario runs")
	}
}

type DeleteChaosWorkflowData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosWorkflowDetails `json:"data"`
}

type DeleteChaosWorkflowDetails struct {
	IsDeleted bool `json:"deleteChaosWorkflow"`
}

type DeleteChaosWorkflowGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID     string  `json:"projectID"`
		WorkflowID    *string `json:"workflowID"`
		WorkflowRunID *string `json:"workflowRunID"`
	} `json:"variables"`
}

// DeleteChaosWorkflow sends GraphQL API request for deleting a given Chaos Workflow.
func DeleteChaosWorkflow(projectID string, workflowID *string, cred types.Credentials) (DeleteChaosWorkflowData, error) {

	var gqlReq DeleteChaosWorkflowGraphQLRequest
	var err error

	gqlReq.Query = `mutation deleteChaosWorkflow($projectID: String!, $workflowID: String, $workflowRunID: String) {
                      deleteChaosWorkflow(
                        projectID: $projectID
                        workflowID: $workflowID
                        workflowRunID: $workflowRunID
                      )
                    }`
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.WorkflowID = workflowID
	var workflow_run_id string = ""
	gqlReq.Variables.WorkflowRunID = &workflow_run_id

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DeleteChaosWorkflowData{}, err
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
		return DeleteChaosWorkflowData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DeleteChaosWorkflowData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var deletedWorkflow DeleteChaosWorkflowData
		err = json.Unmarshal(bodyBytes, &deletedWorkflow)
		if err != nil {
			return DeleteChaosWorkflowData{}, err
		}

		if len(deletedWorkflow.Errors) > 0 {
			return DeleteChaosWorkflowData{}, errors.New(deletedWorkflow.Errors[0].Message)
		}

		return deletedWorkflow, nil
	} else {
		return DeleteChaosWorkflowData{}, errors.New("Error while deleting the Chaos Scenario")
	}
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

// GetServerVersion fetches the GQL server version
func GetServerVersion(endpoint string) (ServerVersionResponse, error) {
	query := `{"query":"query{\n getServerVersion{\n key value\n }\n}"}`
	resp, err := SendRequest(
		SendRequestParams{
			Endpoint: endpoint + utils.GQLAPIPath,
		},
		[]byte(query),
		string(types.Post),
	)
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
		return ServerVersionResponse{}, errors.New(resp.Status)
	}
}
