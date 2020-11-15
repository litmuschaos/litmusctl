package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type Agent struct {
	AgentName      string `json:"cluster_name"`
	Mode           string
	Description    string `json:"description,omitempty"`
	PlatformName   string `json:"platform_name"`
	ProjectId      string `json:"project_id"`
	ClusterType    string `json:"cluster_type"`
	Namespace      string
	ServiceAccount string
	NsExists       bool
	SAExists       bool
}

type AgentRegistrationData struct {
	Errors []Errors      `json:"errors"`
	Data   AgentRegister `json:"data"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type AgentRegister struct {
	UserAgentReg UserAgentReg `json:"userClusterReg"`
}

type UserAgentReg struct {
	ClusterID   string `json:"cluster_id"`
	ClusterName string `json:"cluster_name"`
	Token       string `json:"token"`
}

// GetAgentDetails take details of agent as input
func GetAgentDetails(pid string, t Token, cred Credentials) Agent {
	var newAgent Agent
	// Get agent name as input
	fmt.Println("\n🔗 Enter the details of the agent ----")
	fmt.Print("🤷 Agent Name: ")
	fmt.Scanln(&newAgent.AgentName)
	for newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		fmt.Print("🤷 Agent Name: ")
		fmt.Scanln(&newAgent.AgentName)
	}
	i := 0
	// Check if agent with the given name already exists
	for AgentExists(pid, newAgent.AgentName, t, cred) {
		// Print agent list if existing agent name is entered twice
		if i < 1 {
			fmt.Println("🚫 Agent with the given name already exists.\n❗ Please enter a different name.")
			fmt.Print("🤷 Agent Name: ")
			fmt.Scanln(&newAgent.AgentName)
			i++
		} else {
			fmt.Println("🚫 Agent with the given name already exists.")
			GetAgentList(pid, t, cred)
			fmt.Println("❗ Please enter a different name.")
			fmt.Print("\n🤷 Agent Name: ")
			fmt.Scanln(&newAgent.AgentName)
		}
	}
	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = Scanner()
	// Get platform name as input
	fmt.Print("📦 Platform Name [", defaultPlatform, "]: ")
	fmt.Scanln(&newAgent.PlatformName)
	if newAgent.PlatformName == "" {
		newAgent.PlatformName = defaultPlatform
	}
	// Set agent type
	newAgent.ClusterType = agentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = ValidNs()

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
func AgentExists(pid, agentName string, t Token, cred Credentials) bool {

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
				"%s/chaos/api/graphql/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return true
	}
	for i, _ := range agents.Data.GetAgent {
		if agentName == agents.Data.GetAgent[i].AgentName {
			return true
		}
	}
	return false
}

// GetAgentList lists the agent connected to the specified project
func GetAgentList(pid string, t Token, cred Credentials) {
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
				"%s/chaos/api/graphql/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println(err)
	}
	fmt.Println("\n📘 Registered agents list -----------")
	fmt.Println()
	for i, _ := range agents.Data.GetAgent {
		fmt.Println("-", agents.Data.GetAgent[i].AgentName)
	}
	fmt.Println("\n-------------------------------------")
}

// RegisterAgent registers the agent with the given details
func RegisterAgent(c Agent, t Token, cred Credentials) (AgentRegistrationData, error) {
	var cr AgentRegistrationData
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
				"%s/chaos/api/graphql/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println(err)
		fmt.Println(resp.IsSuccess())
		return AgentRegistrationData{}, err
	}
	return cr, nil
}
