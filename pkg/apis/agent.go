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

	"github.com/litmuschaos/litmusctl/pkg/utils"

	types "github.com/litmuschaos/litmusctl/pkg/types"
)

type AgentData struct {
	Data   AgentList `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type AgentDetails struct {
	AgentName    string `json:"clusterName"`
	IsActive     bool   `json:"isActive"`
	IsRegistered bool   `json:"isRegistered"`
	ClusterID    string `json:"clusterID"`
}

type AgentList struct {
	GetAgent []AgentDetails `json:"listClusters"`
}

// GetAgentList lists the agent connected to the specified project
func GetAgentList(c types.Credentials, pid string) (AgentData, error) {
	query := `{"query":"query{\n  listClusters(projectID: \"` + pid + `\"){\n  clusterID clusterName isActive \n  }\n}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: c.Endpoint + utils.GQLAPIPath, Token: c.Token}, []byte(query), string(types.Post))
	if err != nil {
		return AgentData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return AgentData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var agent AgentData
		err = json.Unmarshal(bodyBytes, &agent)
		if err != nil {
			return AgentData{}, err
		}

		if len(agent.Errors) > 0 {
			return AgentData{}, errors.New(agent.Errors[0].Message)
		}

		return agent, nil
	} else {
		return AgentData{}, err
	}
}

type AgentConnectionData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data AgentConnect `json:"data"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type AgentConnect struct {
	UserAgentReg UserAgentReg `json:"registerCluster"`
}

type UserAgentReg struct {
	ClusterID   string `json:"clusterID"`
	ClusterName string `json:"clusterName"`
	Token       string `json:"token"`
}

// ConnectAgent connects the agent with the given details
func ConnectAgent(agent types.Agent, cred types.Credentials) (AgentConnectionData, error) {
	query := `{"query":"mutation {\n  registerCluster(request: \n    { \n    clusterName: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatformName: \"` + agent.PlatformName + `\",\n    projectID: \"` + agent.ProjectId + `\",\n    clusterType: \"` + agent.ClusterType + `\",\n  agentScope: \"` + agent.Mode + `\",\n    agentNamespace: \"` + agent.Namespace + `\",\n    serviceAccount: \"` + agent.ServiceAccount + `\",\n    skipSsl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agentNsExists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agentSaExists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    clusterID\n    clusterName\n    token\n  }\n}"}`

	if agent.NodeSelector != "" {
		query = `{"query":"mutation {\n  registerCluster(request: \n    { \n    clusterName: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  nodeSelector: \"` + agent.NodeSelector + `\",\n  \tplatformName: \"` + agent.PlatformName + `\",\n    projectID: \"` + agent.ProjectId + `\",\n    clusterType: \"` + agent.ClusterType + `\",\n  agentScope: \"` + agent.Mode + `\",\n    agentNamespace: \"` + agent.Namespace + `\",\n    skipSsl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceAccount: \"` + agent.ServiceAccount + `\",\n    agentNsExists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agentSaExists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    clusterID\n    clusterName\n    token\n  }\n}"}`
	}

	if agent.Tolerations != "" {
		query = `{"query":"mutation {\n  registerCluster(request: \n    { \n    clusterName: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatformName: \"` + agent.PlatformName + `\",\n    projectID: \"` + agent.ProjectId + `\",\n    clusterType: \"` + agent.ClusterType + `\",\n  agentScope: \"` + agent.Mode + `\",\n    agentNamespace: \"` + agent.Namespace + `\",\n    serviceAccount: \"` + agent.ServiceAccount + `\",\n    skipSsl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agentNsExists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agentSaExists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    clusterID\n    clusterName\n    token\n  }\n}"}`
	}

	if agent.NodeSelector != "" && agent.Tolerations != "" {
		query = `{"query":"mutation {\n  registerCluster(request: \n    { \n    clusterName: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  nodeSelector: \"` + agent.NodeSelector + `\",\n  \tplatformName: \"` + agent.PlatformName + `\",\n    projectID: \"` + agent.ProjectId + `\",\n    clusterType: \"` + agent.ClusterType + `\",\n  agentScope: \"` + agent.Mode + `\",\n    agentNamespace: \"` + agent.Namespace + `\",\n    skipSsl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceAccount: \"` + agent.ServiceAccount + `\",\n    agentNsExists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agentSaExists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    clusterID\n    clusterName\n    token\n  }\n}"}`
	}

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query), string(types.Post))
	if err != nil {
		return AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectAgent AgentConnectionData
		err = json.Unmarshal(bodyBytes, &connectAgent)
		if err != nil {
			return AgentConnectionData{}, errors.New("Error in registering agent: " + err.Error())
		}

		if len(connectAgent.Errors) > 0 {
			return AgentConnectionData{}, errors.New(connectAgent.Errors[0].Message)
		}
		return connectAgent, nil
	} else {
		return AgentConnectionData{}, err
	}
}

type DeleteAgentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteAgentDetails `json:"data"`
}

type DeleteAgentDetails struct {
	Message string `json:"deleteClusters"`
}

type DeleteAgentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID  string    `json:"projectID"`
		ClusterIDs []*string `json:"clusterIDs"`
	} `json:"variables"`
}

// DeleteAgent sends GraphQL API request for deleting ChaosAgent(s).
func DeleteAgent(projectID string, clusterIDs []*string, cred types.Credentials) (DeleteAgentData, error) {

	var gqlReq DeleteAgentGraphQLRequest
	var err error

	gqlReq.Query = `mutation deleteClusters($projectID: String!, $clusterIDs: [String]!) {
                      deleteClusters(
                        projectID: $projectID
                        clusterIDs: $clusterIDs
                      )
                    }`
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.ClusterIDs = clusterIDs

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DeleteAgentData{}, err
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
		return DeleteAgentData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DeleteAgentData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var deleteAgentData DeleteAgentData
		err = json.Unmarshal(bodyBytes, &deleteAgentData)
		if err != nil {
			return DeleteAgentData{}, err
		}

		if len(deleteAgentData.Errors) > 0 {
			return DeleteAgentData{}, errors.New(deleteAgentData.Errors[0].Message)
		}

		return deleteAgentData, nil
	} else {
		return DeleteAgentData{}, err
	}
}
