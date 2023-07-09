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
	"fmt"
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
		ProjectID        string                  `json:"pid"`
		ListInfraRequest models.ListInfraRequest `json:"request"`
	} `json:"variables"`
}

// GetInfraList lists the Chaos Infrastructure connected to the specified project
func GetInfraList(c types.Credentials, pid string, request models.ListInfraRequest) (InfraData, error) {
	var gplReq ListInfraGraphQLRequest
	gplReq.Query = `query listInfras($pid: projectID!, $request: ListInfraRequest!){
					listInfras(projectID: $pid, ListInfraRequest: $request){
						totalNoOfInfras
						Infras {
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

type InfraConnectionData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data InfraConnect `json:"data"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type InfraConnect struct {
	UserInfraReg InfraAgentReg `json:"registerInfra"`
}

type InfraAgentReg struct {
	InfraID   string `json:"infraID"`
	InfraName string `json:"name"`
	Token     string `json:"token"`
}

// ConnectInfra connects the  Infra with the given details
func ConnectInfra(infra types.Infra, cred types.Credentials) (InfraConnectionData, error) {
	query := `{"query":"mutation {\n  registerInfra(projectID: \"` + infra.ProjectId + `\", request: \n    { \n    name: \"` + infra.InfraName + `\", \n    description: \"` + infra.Description + `\",\n  \tplatformName: \"` + infra.PlatformName + `\",\n    infrastructureType: \"` + infra.InfraType + `\",\n  infraScope: \"` + infra.Mode + `\",\n    infraNamespace: \"` + infra.Namespace + `\",\n    serviceAccount: \"` + infra.ServiceAccount + `\",\n    skipSsl: ` + fmt.Sprintf("%t", infra.SkipSSL) + `,\n    infraNsExists: ` + fmt.Sprintf("%t", infra.NsExists) + `,\n    agentSaExists: ` + fmt.Sprintf("%t", infra.SAExists) + `,\n  }){\n    infraID\n    name\n    token\n  }\n}"}`

	if infra.NodeSelector != "" {
		query = `{"query":"mutation {\n  registerInfra(projectID: \"` + infra.ProjectId + `\", request: \n    { \n    name: \"` + infra.InfraName + `\", \n    description: \"` + infra.Description + `\",\n  nodeSelector: \"` + infra.NodeSelector + `\",\n  \tplatformName: \"` + infra.PlatformName + `\",\n   infrastructureType: \"` + infra.InfraType + `\",\n  infraScope: \"` + infra.Mode + `\",\n   infraNamespace: \"` + infra.Namespace + `\",\n    skipSsl: ` + fmt.Sprintf("%t", infra.SkipSSL) + `,\n    serviceAccount: \"` + infra.ServiceAccount + `\",\n    infraNsExists: ` + fmt.Sprintf("%t", infra.NsExists) + `,\n    infraSaExists: ` + fmt.Sprintf("%t", infra.SAExists) + `,\n  }){\n    infraID\n    name\n    token\n  }\n}"}`
	}

	if infra.Tolerations != "" {
		query = `{"query":"mutation {\n  registerInfra(projectID: \"` + infra.ProjectId + `\",request: \n    { \n    infraName: \"` + infra.InfraName + `\", \n    description: \"` + infra.Description + `\",\n  \tplatformName: \"` + infra.PlatformName + `\",\n    infraType: \"` + infra.InfraType + `\",\n  infraScope: \"` + infra.Mode + `\",\n    infraNamespace: \"` + infra.Namespace + `\",\n    serviceAccount: \"` + infra.ServiceAccount + `\",\n    skipSsl: ` + fmt.Sprintf("%t", infra.SkipSSL) + `,\n    infraExists: ` + fmt.Sprintf("%t", infra.NsExists) + `,\n    infraSaExists: ` + fmt.Sprintf("%t", infra.SAExists) + `,\n tolerations: ` + infra.Tolerations + ` }){\n    infraID\n    name\n    token\n  }\n}"}`
	}

	if infra.NodeSelector != "" && infra.Tolerations != "" {
		query = `{"query":"mutation {\n  registerInfra(projectID: \"` + infra.ProjectId + `\", request: \n    { \n    infraName: \"` + infra.InfraName + `\", \n    description: \"` + infra.Description + `\",\n  nodeSelector: \"` + infra.NodeSelector + `\",\n  \tplatformName: \"` + infra.PlatformName + `\",\n    infraType: \"` + infra.InfraType + `\",\n  infraScope: \"` + infra.Mode + `\",\n    infraNamespace: \"` + infra.Namespace + `\",\n    serviceAccount: \"` + infra.ServiceAccount + `\",\n    skipSsl: ` + fmt.Sprintf("%t", infra.SkipSSL) + `,\n    infraExists: ` + fmt.Sprintf("%t", infra.NsExists) + `,\n    infraSaExists: ` + fmt.Sprintf("%t", infra.SAExists) + `,\n tolerations: ` + infra.Tolerations + ` }){\n    infraID\n    name\n    token\n  }\n}"}`
	}

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query), string(types.Post))
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
		ProjectID string    `json:"projectID"`
		InfraID   []*string `json:"infraID"`
	} `json:"variables"`
}

// DisconnectInfra sends GraphQL API request for disconnecting Chaos Delegate(s).
func DisconnectInfra(projectID string, infraID []*string, cred types.Credentials) (DisconnectInfraData, error) {

	var gqlReq DisconnectInfraGraphQLRequest
	var err error

	gqlReq.Query = `mutation deleteInfra($projectID: String!, $infraID: [String]!) {
                      deleteClusters(
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
