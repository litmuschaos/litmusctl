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
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

type SaveExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data SavedExperimentDetails `json:"data"`
}

type SavedExperimentDetails struct {
	Message string `json:"saveChaosExperiment"`
}

type SaveChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                  string                           `json:"projectID"`
		SaveChaosExperimentRequest model.SaveChaosExperimentRequest `json:"request"`
	} `json:"variables"`
}

type RunExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data RunExperimentDetails `json:"data"`
}

type RunExperimentDetails struct {
	RunExperimentResponse model.RunChaosExperimentResponse `json:"runChaosExperiment"`
}

// CreateExperiment sends GraphQL API request for creating a Experiment
func CreateExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (RunExperimentData, error) {

	// Query to Save the Experiment
	var gqlReq SaveChaosExperimentGraphQLRequest

	gqlReq.Query = `mutation saveChaosExperiment($projectID: ID!, $request: SaveChaosExperimentRequest!) {
                      saveChaosExperiment(projectID: $projectID, request: $request)
                    }`
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.SaveChaosExperimentRequest = requestData

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return RunExperimentData{}, err
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
		return RunExperimentData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return RunExperimentData{}, errors.New("Error in creating Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var savedExperiment SaveExperimentData

		err = json.Unmarshal(bodyBytes, &savedExperiment)

		if err != nil {
			return RunExperimentData{}, errors.New("Error in creating Chaos Experiment: " + err.Error())
		}

		// Errors present
		if len(savedExperiment.Errors) > 0 {
			return RunExperimentData{}, errors.New(savedExperiment.Errors[0].Message)
		}

		if strings.Contains(savedExperiment.Data.Message, "experiment saved successfully") {
			fmt.Print("\n🚀 Chaos Experiment successfully Saved 🎉")
		}
	} else {
		return RunExperimentData{}, errors.New("graphql schema error")
	}

	// Query to Run the Chaos Experiment
	runQuery := `{"query":"mutation{ \n runChaosExperiment(experimentID:  \"` + requestData.ID + `\", projectID:  \"` + pid + `\"){\n notifyID \n}}"}`
	resp, err = SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(runQuery), string(types.Post))

	if err != nil {
		return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var runExperiment RunExperimentData
		err = json.Unmarshal(bodyBytes, &runExperiment)
		if err != nil {
			return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
		}

		if len(runExperiment.Errors) > 0 {
			return RunExperimentData{}, errors.New(runExperiment.Errors[0].Message)
		}
		return runExperiment, nil
	} else {
		return RunExperimentData{}, err
	}
}

func SaveExperiment(pid string, requestData model.SaveChaosExperimentRequest, cred types.Credentials) (SaveExperimentData, error) {

	// Query to Save the Experiment
	var gqlReq SaveChaosExperimentGraphQLRequest

	gqlReq.Query = `mutation saveChaosExperiment($projectID: ID!, $request: SaveChaosExperimentRequest!) {
                      saveChaosExperiment(projectID: $projectID, request: $request)
                    }`
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.SaveChaosExperimentRequest = requestData

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return SaveExperimentData{}, err
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
		return SaveExperimentData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()
	if err != nil {
		return SaveExperimentData{}, errors.New("Error in creating Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var savedExperiment SaveExperimentData

		err = json.Unmarshal(bodyBytes, &savedExperiment)

		if err != nil {
			return SaveExperimentData{}, errors.New("Error in creating Chaos Experiment: " + err.Error())
		}

		// Errors present
		if len(savedExperiment.Errors) > 0 {
			return SaveExperimentData{}, errors.New(savedExperiment.Errors[0].Message)
		}

		if strings.Contains(savedExperiment.Data.Message, "experiment saved successfully") {
			fmt.Print("\n🚀 Chaos Experiment successfully Saved 🎉")
		}
		return savedExperiment, nil

	} else {
		return SaveExperimentData{}, errors.New("graphql schema error")
	}

}

func RunExperiment(pid string, eid string, cred types.Credentials) (RunExperimentData, error) {
	var err error
	runQuery := `{"query":"mutation{ \n runChaosExperiment(experimentID:  \"` + eid + `\", projectID:  \"` + pid + `\"){\n notifyID \n}}"}`

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(runQuery), string(types.Post))

	if err != nil {
		return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var runExperiment RunExperimentData
		err = json.Unmarshal(bodyBytes, &runExperiment)
		if err != nil {
			return RunExperimentData{}, errors.New("Error in Running Chaos Experiment: " + err.Error())
		}

		if len(runExperiment.Errors) > 0 {
			return RunExperimentData{}, errors.New(runExperiment.Errors[0].Message)
		}
		return runExperiment, nil
	} else {
		return RunExperimentData{}, err
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
	ListExperimentDetails model.ListExperimentResponse `json:"listExperiment"`
}

type GetChaosExperimentsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosExperimentRequest model.ListExperimentRequest `json:"request"`
		ProjectID                 string                      `json:"projectID"`
	} `json:"variables"`
}

// GetExperimentList sends GraphQL API request for fetching a list of experiments.
func GetExperimentList(pid string, in model.ListExperimentRequest, cred types.Credentials) (ExperimentListData, error) {

	var gqlReq GetChaosExperimentsGraphQLRequest
	var err error

	gqlReq.Query = `query listExperiment($projectID: ID!, $request: ListExperimentRequest!) {
                      listExperiment(projectID: $projectID, request: $request) {
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
                          updatedBy{
                              username
                              email
                        }
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
		return ExperimentListData{}, errors.New("Error while fetching the Chaos Experiments")
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
	ListExperimentRunDetails model.ListExperimentRunResponse `json:"listExperimentRun"`
}

type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                         `json:"projectID"`
		GetChaosExperimentRunRequest model.ListExperimentRunRequest `json:"request"`
	} `json:"variables"`
}

// GetExperimentRunsList sends GraphQL API request for fetching a list of experiment runs.
func GetExperimentRunsList(pid string, in model.ListExperimentRunRequest, cred types.Credentials) (ExperimentRunListData, error) {

	var gqlReq GetChaosExperimentRunGraphQLRequest
	var err error

	gqlReq.Query = `query listExperimentRuns($projectID: ID!, $request: ListExperimentRunRequest!) {
                      listExperimentRun(projectID: $projectID, request: $request) {
                        totalNoOfExperimentRuns
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
                          executionData
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
		var experimentRunList ExperimentRunListData
		err = json.Unmarshal(bodyBytes, &experimentRunList)
		if err != nil {
			return ExperimentRunListData{}, err
		}

		if len(experimentRunList.Errors) > 0 {
			return ExperimentRunListData{}, errors.New(experimentRunList.Errors[0].Message)
		}

		return experimentRunList, nil
	} else {
		return ExperimentRunListData{}, errors.New("Error while fetching the Chaos Experiment runs")
	}
}

type DeleteChaosExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosExperimentDetails `json:"data"`
}

type DeleteChaosExperimentDetails struct {
	IsDeleted bool `json:"deleteChaosExperiment"`
}

type DeleteChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID       string  `json:"projectID"`
		ExperimentID    *string `json:"experimentID"`
		ExperimentRunID *string `json:"experimentRunID"`
	} `json:"variables"`
}

// DeleteChaosExperiment sends GraphQL API request for deleting a given Chaos Experiment.
func DeleteChaosExperiment(projectID string, experimentID *string, cred types.Credentials) (DeleteChaosExperimentData, error) {

	var gqlReq DeleteChaosExperimentGraphQLRequest
	var err error

	gqlReq.Query = `mutation deleteChaosExperiment($projectID: ID!, $experimentID: String!, $experimentRunID: String) {
                      deleteChaosExperiment(
                        projectID: $projectID
                        experimentID: $experimentID
                        experimentRunID: $experimentRunID
                      )
                    }`
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.ExperimentID = experimentID
	//var experiment_run_id string = ""
	//gqlReq.Variables.ExperimentRunID = &experiment_run_id

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DeleteChaosExperimentData{}, err
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
		return DeleteChaosExperimentData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DeleteChaosExperimentData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var deletedExperiment DeleteChaosExperimentData
		err = json.Unmarshal(bodyBytes, &deletedExperiment)
		if err != nil {
			return DeleteChaosExperimentData{}, err
		}

		if len(deletedExperiment.Errors) > 0 {
			return DeleteChaosExperimentData{}, errors.New(deletedExperiment.Errors[0].Message)
		}

		return deletedExperiment, nil
	} else {
		return DeleteChaosExperimentData{}, errors.New("Error while deleting the Chaos Experiment")
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
