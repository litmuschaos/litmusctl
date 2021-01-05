package cmd

import (
	"fmt"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/mayadata-io/kuberactl/pkg/types/propel"
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
	newAgent.AgentName = Scanner()
	for newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		fmt.Print("🤷 Agent Name: ")
		newAgent.AgentName = Scanner()
	}
	i := 0
	// Check if agent with the given name already exists
	for AgentExists(pid, newAgent.AgentName, t, cred) {
		// Print agent list if existing agent name is entered twice
		if i < 1 {
			fmt.Println("🚫 Agent with the given name already exists.\n❗ Please enter a different name.")
			fmt.Print("🤷 Agent Name: ")
			newAgent.AgentName = Scanner()
			i++
		} else {
			fmt.Println("🚫 Agent with the given name already exists.")
			GetAgentList(pid, t, cred)
			fmt.Println("❗ Please enter a different name.")
			fmt.Print("\n🤷 Agent Name: ")
			newAgent.AgentName = Scanner()
		}
	}
	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = Scanner()
	// Get platform name as input
	newAgent.PlatformName = getPlatformName()
	// Set agent type
	newAgent.ClusterType = agentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = ValidNs(chaosAgentLabel)

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

// GetPropelAgentDetails take details of agent as input
func GetPropelAgentDetails(pid string, t Token, cred Credentials) Agent {
	var newAgent Agent
	agentsList, err := ListPropelAgents(pid, t, cred)
	if err != nil {
		fmt.Printf("🚫 List propel agents failed: [%s]", err.Error())
		os.Exit(1)
	}
	// Get agent name as input
	fmt.Println("\n🔗 Enter the details of the agent ----")
	fmt.Print("🤷 Agent Name: ")
	newAgent.AgentName = Scanner()
	for newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		fmt.Print("🤷 Agent Name: ")
		newAgent.AgentName = Scanner()
	}
	i := 0
	// Check if agent with the given name already exists
	for IsPropelAgentExists(agentsList, newAgent.AgentName) {
		// Print agent list if existing agent name is entered twice
		if i < 1 {
			fmt.Println("🚫 Agent with the given name already exists.\n❗ Please enter a different name.")
			fmt.Print("🤷 Agent Name: ")
			newAgent.AgentName = Scanner()
			i++
		} else {
			fmt.Println("🚫 Agent with the given name already exists.")
			PrintPropelAgents(agentsList)
			fmt.Println("❗ Please enter a different name.")
			fmt.Print("\n🤷 Agent Name: ")
			newAgent.AgentName = Scanner()
		}
	}
	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = Scanner()
	// Get platform name as input
	newAgent.PlatformName = getPlatformName()
	// Set agent type
	newAgent.ClusterType = propelAgentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = ValidNs(propelAgentLabel)

	return newAgent
}

func ListPropelAgents(pid string, t Token, cred Credentials) (propel.PropelAgentList, error) {
	var agents propel.PropelAgentList
	client := resty.New()
	bodyData := `{"query":"query{\n  listClusters(input:{\n    project_id: \"` + fmt.Sprintf("%s", pid) + `\"\n  }){\n    clusters{\n      cluster_id\n      cluster_name\n      platform_name\n      cluster_type\n      description\n      token\n    }\n  }\n}"}`
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
				"%s/propel/api/graphql/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return propel.PropelAgentList{}, err
	}

	return agents, nil
}

// IsPropelAgentExists checks if an agent of the given name already exists
func IsPropelAgentExists(agents propel.PropelAgentList, agentName string) bool {
	for i, _ := range agents.Data.ListAgents.Agents {
		if agentName == agents.Data.ListAgents.Agents[i].AgentName {
			return true
		}
	}
	return false
}

func PrintPropelAgents(agents propel.PropelAgentList) {
	fmt.Println("\n📘 Registered agents list -----------")
	fmt.Println()
	for i, _ := range agents.Data.ListAgents.Agents {
		fmt.Println("-", agents.Data.ListAgents.Agents[i].AgentName)
	}
	fmt.Println("\n-------------------------------------")
}

// RegisterPropelAgent registers the agent with the given details
func RegisterPropelAgent(c Agent, t Token, cred Credentials) (propel.AgentDetails, error) {
	var cr propel.AgentDetails
	client := resty.New()
	// bodyData := `{"query":"mutation {\n  userClusterReg(clusterInput: \n    { \n    cluster_name: \"` + fmt.Sprintf("%s", c.AgentName) + `\", \n    description: \"` + fmt.Sprintf("%s", c.Description) + `\",\n  \tplatform_name: \"` + fmt.Sprintf("%s", c.PlatformName) + `\",\n    project_id: \"` + fmt.Sprintf("%s", c.ProjectId) + `\",\n    cluster_type: \"` + fmt.Sprintf("%s", c.ClusterType) + `\",\n  agent_scope: \"` + fmt.Sprintf("%s", c.Mode) + `\",\n    agent_namespace: \"` + fmt.Sprintf("%s", c.Namespace) + `\",\n    serviceaccount: \"` + fmt.Sprintf("%s", c.ServiceAccount) + `\",\n    agent_ns_exists: ` + fmt.Sprintf("%t", c.NsExists) + `,\n    agent_sa_exists: ` + fmt.Sprintf("%t", c.SAExists) + `,\n  }){\n    cluster_id\n    cluster_name\n    token\n  }\n}"}`
	bodyData := `{"query":"mutation{\n  addCluster(clusterInput: {\n    cluster_name: \"` + fmt.Sprintf("%s", c.AgentName) + `\"\n    description: \"` + fmt.Sprintf("%s", c.Description) + `\"\n    platform_name: \"` + fmt.Sprintf("%s", c.PlatformName) + `\"\n    project_id: \"` + fmt.Sprintf("%s", c.ProjectId) + `\"\n    cluster_type: ` + fmt.Sprintf("%s", c.ClusterType) + `\n    namespace: \"` + fmt.Sprintf("%s", c.Namespace) + `\"\n  }){\n    cluster_id\n    cluster_token\n    cluster_access_key\n    isClusterConfirmed\n    yaml_route\n    new_cluster{\n      cluster_id\n      project_id\n      cluster_name\n      description\n      platform_name\n      access_key\n      is_registered\n      is_cluster_confirmed\n      is_active\n      updated_at\n      created_at\n      cluster_type\n      token\n      cluster_yaml_route\n      cluster_url\n      is_self_cluster\n      namespace\n    }\n  }\n}"}`
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
				"%s/propel/api/graphql/query",
				cred.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		fmt.Println(err)
		fmt.Println(resp.IsSuccess())
		return propel.AgentDetails{}, err
	}
	return cr, nil
}
