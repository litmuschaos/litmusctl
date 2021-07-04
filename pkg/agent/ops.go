package agent

import (
	"fmt"
	"os"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

// GetProject display list of projects and returns the project id based on input
func GetProjectID(u apis.ProjectDetails) string {
	var pid int
	fmt.Println("\n✨ Projects List:")
	for index := range u.Data.GetUser.Projects {
		fmt.Printf("%d.  %s\n", index+1, u.Data.GetUser.Projects[index].Name)
	}

repeat:
	fmt.Print("\n🔎 Select Project: ")
	fmt.Scanln(&pid)

	for pid < 1 || pid > len(u.Data.GetUser.Projects) {
		fmt.Println("❗ Invalid Project. Please select a correct one.")
		goto repeat
	}

	return u.Data.GetUser.Projects[pid-1].ID
}

// GetMode gets mode of agent installation as input
func GetModeType() string {
	var mode int = 1
	fmt.Println("\n🔌 Installation Modes:\n1. Cluster\n2. Namespace")
	fmt.Print("\n👉 Select Mode [", utils.DefaultMode, "]: ")
	fmt.Scanln(&mode)

repeat:
	if mode == 1 {
		return "cluster"
	}

	if mode == 2 {
		return "namespace"
	}

	for mode < 1 || mode > 2 {
		fmt.Println("🚫 Invalid mode. Please enter the correct mode")
		goto repeat
	}

	return utils.DefaultMode
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

// GetAgentDetails take details of agent as input
func GetAgentDetails(mode string, pid string, c types.Credentials, kubeconfig *string) (types.Agent, error) {
	var newAgent types.Agent
	// Get agent name as input
	fmt.Println("\n🔗 Enter the details of the agent ----")
	// Label for goto statement in case of invalid agent name

AGENT_NAME:
	fmt.Print("🤷 Agent Name: ")
	newAgent.AgentName = utils.Scanner()
	if newAgent.AgentName == "" {
		fmt.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
		goto AGENT_NAME
	}

	// Check if agent with the given name already exists
	agent, err := apis.GetAgentList(c, pid)
	if err != nil {
		return types.Agent{}, err
	}

	var isAgentExist = false
	for i := range agent.Data.GetAgent {
		if newAgent.AgentName == agent.Data.GetAgent[i].AgentName {
			fmt.Println(agent.Data.GetAgent[i].AgentName)
			isAgentExist = true
		}
	}

	if isAgentExist {
		fmt.Println("🚫 Agent with the given name already exists.")
		// Print agent list if existing agent name is entered twice
		fmt.Print("\n📘 Connected agents list -----------\n\n")

		for i := range agent.Data.GetAgent {
			fmt.Println("-", agent.Data.GetAgent[i].AgentName)
		}

		fmt.Println("\n-------------------------------------")

		fmt.Println("❗ Please enter a different name.")
		goto AGENT_NAME
	}

	// Get agent description as input
	fmt.Print("📘 Agent Description: ")
	newAgent.Description = utils.Scanner()
	// Get platform name as input
	newAgent.PlatformName = GetPlatformName(kubeconfig)
	// Set agent type
	newAgent.ClusterType = utils.AgentType
	// Set project id
	newAgent.ProjectId = pid
	// Get namespace
	newAgent.Namespace, newAgent.NsExists = k8s.ValidNs(mode, utils.ChaosAgentLabel, kubeconfig)

	return newAgent, nil
}

func ValidateSAPermissions(mode string, kubeconfig *string) {
	var (
		pems [2]bool
		err  error
	)

	if mode == "cluster" {
		resources := [2]string{"clusterrole", "clusterrolebinding"}
		i := 0
		for _, resource := range resources {
			pems[i], err = k8s.CheckSAPermissions("create", resource, true, kubeconfig)
			if err != nil {
				fmt.Println(err)
			}
			i++
		}
	} else {
		resources := [2]string{"role", "rolebinding"}
		i := 0
		for _, resource := range resources {
			pems[i], err = k8s.CheckSAPermissions("create", resource, true, kubeconfig)
			if err != nil {
				fmt.Println(err)
			}
			i++
		}
	}

	for _, pem := range pems {
		if !pem {
			fmt.Println("\n🚫 You don't have sufficient permissions.\n🙄 Please use a service account with sufficient permissions.")
			os.Exit(1)
		}
	}

	fmt.Println("\n🌟 Sufficient permissions. Connecting Agent")
}

// Summary display the agent details based on input
func Summary(agent types.Agent, kubeconfig *string) {
	fmt.Println("\n📌 Summary --------------------------")
	fmt.Println("\nAgent Name:        ", agent.AgentName)
	fmt.Println("Agent Description: ", agent.Description)
	fmt.Println("Platform Name:     ", agent.PlatformName)
	if ok, _ := k8s.NsExists(agent.Namespace, kubeconfig); ok {
		fmt.Println("Namespace:         ", agent.Namespace)
	} else {
		fmt.Println("Namespace:         ", agent.Namespace, "(new)")
	}

	if k8s.SAExists(agent.Namespace, agent.ServiceAccount, kubeconfig) {
		fmt.Println("Service Account:   ", agent.ServiceAccount)
	} else {
		fmt.Println("Service Account:   ", agent.ServiceAccount, "(new)")
	}

	fmt.Println("Installation Mode: ", agent.Mode)

	fmt.Println("\n-------------------------------------")
}

func ConfirmInstallation() {
	var descision string
	fmt.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&descision)

	if strings.ToLower(descision) == "yes" || strings.ToLower(descision) == "y" {
		fmt.Println("👍 Continuing agent connection!!")
	} else {
		fmt.Println("✋ Exiting agent connection!!")
		os.Exit(1)
	}
}
