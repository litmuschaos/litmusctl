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
	AgentName    string `json:"cluster_name"`
	IsActive     bool   `json:"is_active"`
	IsRegistered bool   `json:"is_registered"`
	ClusterID    string `json:"cluster_id"`
}

type AgentList struct {
	GetAgent []AgentDetails `json:"getCluster"`
}

// GetAgentList lists the agent connected to the specified project
func GetAgentList(c types.Credentials, pid string) (AgentData, error) {
	query := `{"query":"query{\n  getCluster(project_id: \"` + pid + `\"){\n  cluster_id cluster_name is_active \n  }\n}"}`
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
	UserAgentReg UserAgentReg `json:"userClusterReg"`
}

type UserAgentReg struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Token       string `json:"token"`
}

// ConnectAgent connects the agent with the given details
func ConnectAgent(agent types.Agent, cred types.Credentials) (AgentConnectionData, error) {
	query := `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`

	if agent.NodeSelector != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  node_selector: \"` + agent.NodeSelector + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	}

	if agent.Tolerations != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	}

	if agent.NodeSelector != "" && agent.Tolerations != "" {
		query = `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + agent.AgentName + `\", \n    description: \"` + agent.Description + `\",\n  node_selector: \"` + agent.NodeSelector + `\",\n  \tplatform_name: \"` + agent.PlatformName + `\",\n    project_id: \"` + agent.ProjectId + `\",\n    cluster_type: \"` + agent.ClusterType + `\",\n  agent_scope: \"` + agent.Mode + `\",\n    agent_namespace: \"` + agent.Namespace + `\",\n    skip_ssl: ` + fmt.Sprintf("%t", agent.SkipSSL) + `,\n    serviceaccount: \"` + agent.ServiceAccount + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n tolerations: ` + agent.Tolerations + ` }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
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
