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
package agent

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

var (
	cyan = color.New(color.FgCyan, color.Bold)
	red = color.New(color.FgRed)
)

func PrintExistingAgents(agent apis.AgentData) {
	red.Println("\nAgent with the given name already exists.")
	// Print agent list if existing agent name is entered twice
	cyan.Println("\nConnected agents list:")

	for i := range agent.Data.GetAgent {
		cyan.Println("-", agent.Data.GetAgent[i].AgentName)
	}

	cyan.Println("\n❗ Please enter a different name.")
}

// GetProject display list of projects and returns the project id based on input
func GetProjectID(u apis.ProjectDetails) string {
	var pid int
	cyan.Println("Project list:")
	for index := range u.Data.GetUser.Projects {
		cyan.Printf("%d.  %s\n", index+1, u.Data.GetUser.Projects[index].Name)
	}

repeat:
	cyan.Printf("\nSelect a project [Range: 1-%s]: ", fmt.Sprint(len(u.Data.GetUser.Projects)))
	fmt.Scanln(&pid)

	for pid < 1 || pid > len(u.Data.GetUser.Projects) {
		red.Println("❗ Invalid Project. Please select a correct one.")
		goto repeat
	}

	return u.Data.GetUser.Projects[pid-1].ID
}

// GetMode gets mode of agent installation as input
func GetModeType() string {
repeat:
	var (
		cluster_no   = 1
		namespace_no = 2
		mode         = cluster_no
	)

	cyan.Println("\nInstallation Modes:\n1. Cluster\n2. Namespace")
	cyan.Print("\nSelect Mode [Default: ", utils.DefaultMode, "] [Range: 1-2]: ")
	fmt.Scanln(&mode)

	if mode == 1 {
		return "cluster"
	}

	if mode == 2 {
		return "namespace"
	}

	if (mode != cluster_no) || (mode != namespace_no) {
		red.Println("🚫 Invalid mode. Please enter the correct mode")
		goto repeat
	}

	return utils.DefaultMode
}

// GetAgentDetails take details of agent as input
func GetAgentDetails(mode string, pid string, c types.Credentials, kubeconfig *string) (types.Agent, error) {
	var newAgent types.Agent
	// Get agent name as input
	cyan.Println("\nEnter the details of the agent")
	// Label for goto statement in case of invalid agent name

AGENT_NAME:
	cyan.Print("\nAgent Name: ")
	newAgent.AgentName = utils.Scanner()
	if newAgent.AgentName == "" {
		red.Println("⛔ Agent name cannot be empty. Please enter a valid name.")
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
			cyan.Println(agent.Data.GetAgent[i].AgentName)
			isAgentExist = true
		}
	}

	if isAgentExist {
		PrintExistingAgents(agent)
		goto AGENT_NAME
	}

	// Get agent description as input
	cyan.Print("\nAgent Description: ")
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
		pems      [2]bool
		err       error
		resources [2]string
	)

	if mode == "cluster" {
		resources = [2]string{"clusterrole", "clusterrolebinding"}
	} else {
		resources = [2]string{"role", "rolebinding"}
	}

	for i, resource := range resources {
		pems[i], err = k8s.CheckSAPermissions(k8s.CheckSAPermissionsParams{Verb: "create", Resource: resource, Print: true}, kubeconfig)
		if err != nil {
			red.Println(err)
		}
	}

	for _, pem := range pems {
		if !pem {
			red.Println("\n🚫 You don't have sufficient permissions.\n🙄 Please use a service account with sufficient permissions.")
			os.Exit(1)
		}
	}

	cyan.Println("\n🌟 Sufficient permissions. Connecting Agent")
}

// Summary display the agent details based on input
func Summary(agent types.Agent, kubeconfig *string) {
	cyan.Printf("\n📌 Summary \nAgent Name: %s\nAgent Description: %s\nPlatform Name: %s\n", agent.AgentName, agent.Description, agent.PlatformName)
	if ok, _ := k8s.NsExists(agent.Namespace, kubeconfig); ok {
		cyan.Println("Namespace: ", agent.Namespace)
	} else {
		cyan.Println("Namespace: ", agent.Namespace, "(new)")
	}

	if k8s.SAExists(k8s.SAExistsParams{Namespace: agent.Namespace, Serviceaccount: agent.ServiceAccount}, kubeconfig) {
		cyan.Println("Service Account: ", agent.ServiceAccount)
	} else {
		cyan.Println("Service Account: ", agent.ServiceAccount, "(new)")
	}

	cyan.Printf("\nInstallation Mode: %s\n", agent.Mode)
}

func ConfirmInstallation() {
	var descision string
	cyan.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&descision)

	if strings.ToLower(descision) == "yes" || strings.ToLower(descision) == "y" {
		cyan.Println("👍 Continuing agent connection!!")
	} else {
		red.Println("✋ Exiting agent connection!!")
		os.Exit(1)
	}
}

func CreateRandomProject(cred types.Credentials) string {
	rand, err := utils.GenerateRandomString(10)
	utils.PrintError(err)

	projectName := cred.Username + "-" + rand

	project, err := apis.CreateProjectRequest(projectName, cred)
	utils.PrintError(err)

	return project.Data.CreateProject.ID
}
