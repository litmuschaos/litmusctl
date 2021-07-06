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
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"io/ioutil"
	"net/http"

	types "github.com/litmuschaos/litmusctl/pkg/types"
)

type AgentData struct {
	Data AgentList `json:"data"`
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
	resp, err := SendRequest(SendRequestParams{Endpoint: c.Endpoint + utils.GQLAPIPath, Token: c.Token}, []byte(query))
	if err != nil {
		fmt.Println("Error in getting agent list: ", err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("Error in getting agent list: ", err)
	}

	if resp.StatusCode == http.StatusOK {
		var agent AgentData
		err = json.Unmarshal(bodyBytes, &agent)
		if err != nil {
			fmt.Println("Error in getting agent list: ", err)
		}

		return agent, nil

	} else {
		return AgentData{}, err
	}

	return AgentData{}, err
}

type AgentConnectionData struct {
	Errors []Errors     `json:"errors"`
	Data   AgentConnect `json:"data"`
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
	query := `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + fmt.Sprintf("%s", agent.AgentName) + `\", \n    description: \"` + fmt.Sprintf("%s", agent.Description) + `\",\n  \tplatform_name: \"` + fmt.Sprintf("%s", agent.PlatformName) + `\",\n    project_id: \"` + fmt.Sprintf("%s", agent.ProjectId) + `\",\n    cluster_type: \"` + fmt.Sprintf("%s", agent.ClusterType) + `\",\n  agent_scope: \"` + fmt.Sprintf("%s", agent.Mode) + `\",\n    agent_namespace: \"` + fmt.Sprintf("%s", agent.Namespace) + `\",\n    serviceaccount: \"` + fmt.Sprintf("%s", agent.ServiceAccount) + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", agent.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", agent.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query))
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

		return connectAgent, nil
	} else {
		return AgentConnectionData{}, err
	}

	return AgentConnectionData{}, err
}
