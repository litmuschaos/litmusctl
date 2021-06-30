package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	types "github.com/litmuschaos/litmusctl/pkg/types"
	"io/ioutil"
	"net/http"
)

type ProjectDetails struct {
	Data Data `json:"data"`
}

type Data struct {
	GetUser GetUser `json:"getUser"`
}

type GetUser struct {
	Projects []Project `json:"projects"`
}

type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}


// GetProjectDetails fetches details of the input user
func GetProjectDetails(c types.Credentials) (ProjectDetails, error) {
	query := `{"query":"query {\n  getUser(username: \"` + c.Username + `\"){\n projects{\n id\n name\n}\n}\n}"}`
	resp, err := SendRequest(c.Endpoint + "/api/query", c.Token, []byte(query))
	if err != nil {
		return ProjectDetails{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ProjectDetails{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var project ProjectDetails
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return ProjectDetails{}, err
		}

		return project, nil
	} else {
		return ProjectDetails{}, errors.New("Unmatached status code:" + string(bodyBytes))
	}

	return ProjectDetails{}, nil
}

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
	query := `{"query":"query{\n  getCluster(project_id: \"` + pid + `\"){\n    cluster_name\n  }\n}"}`
	resp, err := SendRequest(c.Endpoint +  "/api/query", c.Token, []byte(query))
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
	resp, err := SendRequest(cred.Endpoint + "/api/query", cred.Token, []byte(query))
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
