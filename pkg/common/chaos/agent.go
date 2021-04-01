package chaos

import (
	"fmt"

	util "github.com/litmuschaos/litmusctl/pkg/common"
	"github.com/litmuschaos/litmusctl/pkg/common/k8s"

	resty "github.com/go-resty/resty/v2"
	"github.com/litmuschaos/litmusctl/pkg/constants"
)

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

// GetAgentDetails take details of agent as input
func GetAgentDetails(pid string, t util.Token, cred util.Credentials) util.Agent {
	var newAgent util.Agent
	// Get agent name as input
	fmt.Println("\n🔗 Enter the details of the agent ----")
	// Label for goto statement in case of invalid agent name
AGENT_NAME:
	fmt.Print("🤷 Agent Name: ")
	newAgent.AgentName = util.Scanner()
	if newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		goto AGENT_NAME
	}

	// Check if agent with the given name already exists
	if AgentExists(pid, newAgent.AgentName, t, cred) {
		fmt.Println("🚫 Agent with the given name already exists.")
		// Print agent list if existing agent name is entered twice
		GetAgentList(pid, t, cred)
		fmt.Println("❗ Please enter a different name.")
		goto AGENT_NAME
	}
	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = util.Scanner()
	// Get platform name as input
	newAgent.PlatformName = util.GetPlatformName()
	// Set agent type
	newAgent.ClusterType = constants.AgentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = k8s.ValidNs(constants.ChaosAgentLabel)

	return newAgent
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

// AgentExists checks if an agent of the given name already exists
func AgentExists(pid, agentName string, t util.Token, cred util.Credentials) bool {

	var agents AgentData
	client := resty.New()
	bodyData := `{"query":"query{\n  getCluster(project_id: \"` + fmt.Sprintf("%s", pid) + `\"){\n    cluster_name\n  }\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", t.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(bodyData).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&agents).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println("Error getting agent names: ", err)
		return true
	}
	for i := range agents.Data.GetAgent {
		if agentName == agents.Data.GetAgent[i].AgentName {
			fmt.Println(agents.Data.GetAgent[i].AgentName)
			return true
		}
	}
	return false
}

// GetAgentList lists the agent connected to the specified project
func GetAgentList(pid string, t util.Token, cred util.Credentials) {
	var agents AgentData
	client := resty.New()
	bodyData := `{"query":"query{\n  getCluster(project_id: \"` + fmt.Sprintf("%s", pid) + `\"){\n    cluster_name\n  }\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", t.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(bodyData).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&agents).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println("Error in geting agent list: ", err)
	}
	fmt.Print("\n📘 Connected agents list -----------\n\n")

	for i := range agents.Data.GetAgent {
		fmt.Println("-", agents.Data.GetAgent[i].AgentName)
	}
	fmt.Println("\n-------------------------------------")
}

// ConnectAgent connects the agent with the given details
func ConnectAgent(c util.Agent, t util.Token, cred util.Credentials) (AgentConnectionData, error) {
	var cr AgentConnectionData
	client := resty.New()
	bodyData := `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + fmt.Sprintf("%s", c.AgentName) + `\", \n    description: \"` + fmt.Sprintf("%s", c.Description) + `\",\n  \tplatform_name: \"` + fmt.Sprintf("%s", c.PlatformName) + `\",\n    project_id: \"` + fmt.Sprintf("%s", c.ProjectId) + `\",\n    cluster_type: \"` + fmt.Sprintf("%s", c.ClusterType) + `\",\n  agent_scope: \"` + fmt.Sprintf("%s", c.Mode) + `\",\n    agent_namespace: \"` + fmt.Sprintf("%s", c.Namespace) + `\",\n    serviceaccount: \"` + fmt.Sprintf("%s", c.ServiceAccount) + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", c.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", c.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", t.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(bodyData).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&cr).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println(err)
		fmt.Println(resp.IsSuccess())
		return AgentConnectionData{}, err
	}
	return cr, nil
}
