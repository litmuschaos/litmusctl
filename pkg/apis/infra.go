/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a1 copy of the License at

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
	"github.com/litmuschaos/litmusctl/pkg/types"
	"io/ioutil"
	"net/http"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

type InfraData struct {
	Data   InfraList `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type InfraList struct {
	ListInfraDetails models.ListInfraResponse `json:"listInfras"`
}

type ListInfraGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID        string                  `json:"projectID"`
		ListInfraRequest models.ListInfraRequest `json:"request"`
	} `json:"variables"`
}

// GetInfraList lists the Chaos Infrastructure connected to the specified project
func GetInfraList(c types.Credentials, pid string, request models.ListInfraRequest) (InfraData, error) {
	var gplReq ListInfraGraphQLRequest
	gplReq.Query = `query listInfras($projectID: ID!, $request: ListInfraRequest!){
					listInfras(projectID: $projectID, request: $request){
						totalNoOfInfras
						infras {
							infraID
							name
							isActive
						}
					}
					}`
	gplReq.Variables.ProjectID = pid
	gplReq.Variables.ListInfraRequest = request

	query, err := json.Marshal(gplReq)
	if err != nil {
		return InfraData{}, err
	}
	resp, err := SendRequest(SendRequestParams{Endpoint: c.Endpoint + utils.GQLAPIPath, Token: c.Token}, query, string(types.Post))
	if err != nil {
		return InfraData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return InfraData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var Infra InfraData
		err = json.Unmarshal(bodyBytes, &Infra)
		if err != nil {
			return InfraData{}, err
		}

		if len(Infra.Errors) > 0 {
			return InfraData{}, errors.New(Infra.Errors[0].Message)
		}

		return Infra, nil
	} else {
		return InfraData{}, err
	}
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type RegisterInfraGqlRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectId            string                      `json:"projectID"`
		RegisterInfraRequest models.RegisterInfraRequest `json:"request"`
	} `json:"variables"`
}

type InfraConnectionData struct {
	Data   RegisterInfra `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type RegisterInfra struct {
	RegisterInfraDetails models.RegisterInfraResponse `json:"registerInfra"`
}

// ConnectInfra connects the  Infra with the given details
func ConnectInfra(infra types.Infra, cred types.Credentials) (InfraConnectionData, error) {
	var gqlReq RegisterInfraGqlRequest
	gqlReq.Query = `mutation registerInfra($projectID: ID!, $request: RegisterInfraRequest!) {
					  registerInfra(
						projectID: $projectID
						request: $request
					  ) {
						infraID
						name
						token
					  }
					}
					`
	gqlReq.Variables.ProjectId = infra.ProjectId
	gqlReq.Variables.RegisterInfraRequest = CreateRegisterInfraRequest(infra)

	if infra.NodeSelector != "" {
		gqlReq.Variables.RegisterInfraRequest.NodeSelector = &infra.NodeSelector
	}

	if infra.Tolerations != "" {
		var toleration []*models.Toleration
		err := json.Unmarshal([]byte(infra.Tolerations), &toleration)
		utils.PrintError(err)
		gqlReq.Variables.RegisterInfraRequest.Tolerations = toleration
	}

	if infra.NodeSelector != "" && infra.Tolerations != "" {
		gqlReq.Variables.RegisterInfraRequest.NodeSelector = &infra.NodeSelector

		var toleration []*models.Toleration
		err := json.Unmarshal([]byte(infra.Tolerations), &toleration)
		utils.PrintError(err)
		gqlReq.Variables.RegisterInfraRequest.Tolerations = toleration
	}

	query, err := json.Marshal(gqlReq)
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, query, string(types.Post))
	if err != nil {
		return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectInfra InfraConnectionData
		err = json.Unmarshal(bodyBytes, &connectInfra)
		if err != nil {
			return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
		}

		if len(connectInfra.Errors) > 0 {
			return InfraConnectionData{}, errors.New(connectInfra.Errors[0].Message)
		}
		return connectInfra, nil
	} else {
		return InfraConnectionData{}, err
	}
}

func CreateRegisterInfraRequest(infra types.Infra) (request models.RegisterInfraRequest) {
	var mode models.InfrastructureType
	if infra.InfraType == "external" {
		mode = models.InfrastructureTypeExternal
	} else {
		mode = models.InfrastructureTypeInternal
	}
	return models.RegisterInfraRequest{
		Name:               infra.InfraName,
		InfraScope:         infra.Mode,
		Description:        &infra.Description,
		PlatformName:       infra.PlatformName,
		EnvironmentID:      infra.EnvironmentID,
		InfrastructureType: mode,
		InfraNamespace:     &infra.Namespace,
		ServiceAccount:     &infra.ServiceAccount,
		InfraNsExists:      &infra.NsExists,
		InfraSaExists:      &infra.SAExists,
		SkipSsl:            &infra.SkipSSL,
	}
}

type DisconnectInfraData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DisconnectInfraDetails `json:"data"`
}

type DisconnectInfraDetails struct {
	Message string `json:"deleteInfra"`
}

type DisconnectInfraGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string `json:"projectID"`
		InfraID   string `json:"infraID"`
	} `json:"variables"`
}

// DisconnectInfra sends GraphQL API request for disconnecting Chaos Infra(s).
func DisconnectInfra(projectID string, infraID string, cred types.Credentials) (DisconnectInfraData, error) {

	var gqlReq DisconnectInfraGraphQLRequest
	var err error

	gqlReq.Query = `mutation deleteInfra($projectID: ID!, $infraID: String!) {
                      deleteInfra(
                        projectID: $projectID
                        infraID: $infraID
                      )
                    }`
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.InfraID = infraID

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DisconnectInfraData{}, err
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
		return DisconnectInfraData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DisconnectInfraData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var disconnectInfraData DisconnectInfraData
		err = json.Unmarshal(bodyBytes, &disconnectInfraData)
		if err != nil {
			return DisconnectInfraData{}, err
		}

		if len(disconnectInfraData.Errors) > 0 {
			return DisconnectInfraData{}, errors.New(disconnectInfraData.Errors[0].Message)
		}

		return disconnectInfraData, nil
	} else {
		return DisconnectInfraData{}, err
	}
}

type ListEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data EnvironmentsList `json:"data"`
}

type EnvironmentsList struct {
	ListEnvironmentDetails models.ListEnvironmentResponse `json:"listEnvironments"`
}
type EnvironmentListGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string                        `json:"projectID"`
		Request   models.ListEnvironmentRequest `json:"request"`
	}
}

func GetEnvironmentList(pid string, cred types.Credentials) (ListEnvironmentData, error) {
	var err error
	var gqlReq EnvironmentListGQLRequest
	gqlReq.Query = `query listEnvironments($projectID: ID!, $request: ListEnvironmentRequest) {
	                 listEnvironments(projectID: $projectID,request: $request){
						environments {
							environmentID
						}
					}
	               }`

	gqlReq.Variables.Request = models.ListEnvironmentRequest{}
	gqlReq.Variables.ProjectID = pid
	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ListEnvironmentData{}, err
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
