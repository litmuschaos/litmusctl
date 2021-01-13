package util

import (
	"fmt"
	"os"

	"github.com/mayadata-io/cli-utils/pkg/common"
	"github.com/mayadata-io/cli-utils/pkg/common/k8s"

	resty "github.com/go-resty/resty/v2"
	"github.com/mayadata-io/cli-utils/pkg/constants"
)

type PropelAgentList struct {
	Data PropelAgentData `json:"data"`
}
type Agents struct {
	AgentID      string `json:"cluster_id"`
	AgentName    string `json:"cluster_name"`
	PlatformName string `json:"platform_name"`
	AgentType    string `json:"cluster_type"`
	Description  string `json:"description"`
	Token        string `json:"token"`
}
type ListAgents struct {
	Agents []Agents `json:"clusters"`
}
type PropelAgentData struct {
	ListAgents ListAgents `json:"listClusters"`
}

type AgentDetails struct {
	Errors    []common.Errors `json:"errors"`
	AgentData AgentData     `json:"data"`
}
type NewCluster struct {
	ClusterID          string `json:"cluster_id"`
	ProjectID          string `json:"project_id"`
	ClusterName        string `json:"cluster_name"`
	Description        string `json:"description"`
	PlatformName       string `json:"platform_name"`
	AccessKey          string `json:"access_key"`
	IsRegistered       bool   `json:"is_registered"`
	IsClusterConfirmed bool   `json:"is_cluster_confirmed"`
	IsActive           bool   `json:"is_active"`
	UpdatedAt          string `json:"updated_at"`
	CreatedAt          string `json:"created_at"`
	ClusterType        string `json:"cluster_type"`
	Token              string `json:"token"`
	ClusterYamlRoute   string `json:"cluster_yaml_route"`
	ClusterURL         string `json:"cluster_url"`
	IsSelfCluster      bool   `json:"is_self_cluster"`
	Namespace          string `json:"namespace"`
}
type AddCluster struct {
	ClusterID          string     `json:"cluster_id"`
	ClusterToken       string     `json:"cluster_token"`
	ClusterAccessKey   string     `json:"cluster_access_key"`
	IsClusterConfirmed bool       `json:"isClusterConfirmed"`
	YamlRoute          string     `json:"yaml_route"`
	NewCluster         NewCluster `json:"new_cluster"`
}
type AgentData struct {
	AddCluster AddCluster `json:"addCluster"`
}

// GetPropelAgentDetails take details of agent as input
func GetPropelAgentDetails(pid string, t common.Token, cred common.Credentials) common.Agent {
	var newAgent common.Agent
	agentsList, err := ListPropelAgents(pid, t, cred)
	if err != nil {
		fmt.Printf("🚫 List propel agents failed: [%s]", err.Error())
		os.Exit(1)
	}
	// Get agent name as input
	fmt.Println("\n🔗 Enter the details of the agent ----")
	fmt.Print("🤷 Agent Name: ")
	newAgent.AgentName = common.Scanner()
	for newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		fmt.Print("🤷 Agent Name: ")
		newAgent.AgentName = common.Scanner()
	}
	i := 0
	// Check if agent with the given name already exists
	for IsPropelAgentExists(agentsList, newAgent.AgentName) {
		// Print agent list if existing agent name is entered twice
		if i < 1 {
			fmt.Println("🚫 Agent with the given name already exists.\n❗ Please enter a different name.")
			fmt.Print("🤷 Agent Name: ")
			newAgent.AgentName = common.Scanner()
			i++
		} else {
			fmt.Println("🚫 Agent with the given name already exists.")
			PrintPropelAgents(agentsList)
			fmt.Println("❗ Please enter a different name.")
			fmt.Print("\n🤷 Agent Name: ")
			newAgent.AgentName = common.Scanner()
		}
	}
	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = common.Scanner()
	// Get platform name as input
	newAgent.PlatformName = common.GetPlatformName()
	// Set agent type
	newAgent.ClusterType = constants.PropelAgentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = k8s.ValidNs(constants.PropelAgentLabel)

	return newAgent
}

func ListPropelAgents(pid string, t common.Token, cred common.Credentials) (PropelAgentList, error) {
	var agents PropelAgentList
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
		return PropelAgentList{}, err
	}

	return agents, nil
}

// IsPropelAgentExists checks if an agent of the given name already exists
func IsPropelAgentExists(agents PropelAgentList, agentName string) bool {
	for i, _ := range agents.Data.ListAgents.Agents {
		if agentName == agents.Data.ListAgents.Agents[i].AgentName {
			return true
		}
	}
	return false
}

func PrintPropelAgents(agents PropelAgentList) {
	fmt.Println("\n📘 Registered agents list -----------")
	fmt.Println()
	for i, _ := range agents.Data.ListAgents.Agents {
		fmt.Println("-", agents.Data.ListAgents.Agents[i].AgentName)
	}
	fmt.Println("\n-------------------------------------")
}

// RegisterPropelAgent registers the agent with the given details
func RegisterPropelAgent(c common.Agent, t common.Token, cred common.Credentials) (AgentDetails, error) {
	var cr AgentDetails
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
		return AgentDetails{}, err
	}
	return cr, nil
}
