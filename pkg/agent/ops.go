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
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"os"
	"strconv"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

func PrintExistingAgents(agent apis.InfraData) {
	utils.Red.Println("\nChaos Infrastructure with the given name already exists.")
	// Print Chaos Delegate list if existing Chaos Delegate name is entered twice
	utils.White_B.Println("\nConnected Chaos Infrastructure list:")

	for i := range agent.Data.ListInfraDetails.Infras {
		utils.White_B.Println("-", agent.Data.ListInfraDetails.Infras[i].Name)
	}

	utils.White_B.Println("\n❗ Please enter a different name.")
}

// GetProjectID display list of projects and returns the project id based on input
func GetProjectID(u apis.ProjectDetails) string {
	var pid int
	utils.White_B.Println("Project list:")
	for index := range u.Data.Projects {
		utils.White_B.Printf("%d.  %s\n", index+1, u.Data.Projects[index].Name)
	}

repeat:
	utils.White_B.Printf("\nSelect a project [Range: 1-%s]: ", fmt.Sprint(len(u.Data.Projects)))
	fmt.Scanln(&pid)

	for pid < 1 || pid > len(u.Data.Projects) {
		utils.Red.Println("❗ Invalid Project. Please select a correct one.")
		goto repeat
	}

	return u.Data.Projects[pid-1].ID
}

// GetModeType gets mode of Chaos Delegate installation as input
func GetModeType() string {
repeat:
	var (
		cluster_no   = 1
		namespace_no = 2
		mode         = cluster_no
	)

	utils.White_B.Println("\nInstallation Modes:\n1. Cluster\n2. Namespace")
	utils.White_B.Print("\nSelect Mode [Default: ", utils.DefaultMode, "] [Range: 1-2]: ")
	fmt.Scanln(&mode)

	if mode == 1 {
		return "cluster"
	}

	if mode == 2 {
		return "namespace"
	}

	if (mode != cluster_no) || (mode != namespace_no) {
		utils.Red.Println("🚫 Invalid mode. Please enter the correct mode")
		goto repeat
	}

	return utils.DefaultMode
}

// GetInfraDetails take details of Chaos Infrastructure as input
func GetInfraDetails(mode string, pid string, c types.Credentials, kubeconfig *string) (types.Agent, error) {
	var newInfra types.Agent
	// Get agent name as input
	utils.White_B.Println("\nEnter the details of the Chaos Delegate")
	// Label for goto statement in case of invalid Chaos Delegate name

INFRA_NAME:
	utils.White_B.Print("\nChaos Infra Name: ")
	newInfra.AgentName = utils.Scanner()
	if newInfra.AgentName == "" {
		utils.Red.Println("⛔ Chaos Infra name cannot be empty. Please enter a valid name.")
		goto INFRA_NAME
	}

	// Check if Chaos Delegate with the given name already exists
	Infra, err := apis.GetInfraList(c, pid, model.ListInfraRequest{})
	if err != nil {
		return types.Agent{}, err
	}

	var isAgentExist = false
	for i := range Infra.Data.ListInfraDetails.Infras {
		if newInfra.AgentName == Infra.Data.ListInfraDetails.Infras[i].Name {
			utils.White_B.Println(Infra.Data.ListInfraDetails.Infras[i].Name)
			isAgentExist = true
		}
	}

	if isAgentExist {
		PrintExistingAgents(Infra)
		goto INFRA_NAME
	}

	// Get agent description as input
	utils.White_B.Print("\nChaos Infrastructure Description: ")
	newInfra.Description = utils.Scanner()

	utils.White_B.Print("\nDo you want Chaos Infrastructure to skip SSL/TLS check (Y/N) (Default: N): ")
	skipSSLDescision := utils.Scanner()

	if strings.ToLower(skipSSLDescision) == "y" {
		newInfra.SkipSSL = true
	} else {
		newInfra.SkipSSL = false
	}

	utils.White_B.Print("\nDo you want NodeSelector to be added in the Chaos Infrastructure deployments (Y/N) (Default: N): ")
	nodeSelectorDescision := utils.Scanner()

	if strings.ToLower(nodeSelectorDescision) == "y" {
		utils.White_B.Print("\nEnter the NodeSelector (Format: key1=value1,key2=value2): ")
		newInfra.NodeSelector = utils.Scanner()

		if ok := utils.CheckKeyValueFormat(newInfra.NodeSelector); !ok {
			os.Exit(1)
		}
	}

	utils.White_B.Print("\nDo you want Tolerations to be added in the Chaos Infrastructure deployments? (Y/N) (Default: N): ")
	tolerationDescision := utils.Scanner()

	if strings.ToLower(tolerationDescision) == "y" {
		utils.White_B.Print("\nHow many tolerations? ")
		no_of_tolerations := utils.Scanner()

		nts, err := strconv.Atoi(no_of_tolerations)
		utils.PrintError(err)

		str := "["
		for tol := 0; tol < nts; tol++ {
			str += "{"

			utils.White_B.Print("\nToleration count: ", tol+1)

			utils.White_B.Print("\nTolerationSeconds: (Press Enter to ignore)")
			ts := utils.Scanner()

			utils.White_B.Print("\nOperator: ")
			operator := utils.Scanner()
			if operator != "" {
				str += "operator : \\\"" + operator + "\\\" "
			}

			utils.White_B.Print("\nEffect: ")
			effect := utils.Scanner()

			if effect != "" {
				str += "effect: \\\"" + effect + "\\\" "
			}

			if ts != "" {
				str += "tolerationSeconds: " + ts + " "
			}

			utils.White_B.Print("\nKey: ")
			key := utils.Scanner()
			if key != "" {
				str += "key: \\\"" + key + "\\\" "
			}

			utils.White_B.Print("\nValue: ")
			value := utils.Scanner()
			if key != "" {
				str += "value: \\\"" + value + "\\\" "
			}

			str += " }"
		}
		str += "]"

		newInfra.Tolerations = str
	}

	// Get platform name as input
	newInfra.PlatformName = GetPlatformName(kubeconfig)
	// Set agent type
	newInfra.ClusterType = utils.AgentType
	// Set project id
	newInfra.ProjectId = pid
	// Get namespace
	newInfra.Namespace, newInfra.NsExists = k8s.ValidNs(mode, utils.ChaosAgentLabel, kubeconfig)

	return newInfra, nil
}

func ValidateSAPermissions(namespace string, mode string, kubeconfig *string) {
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
		pems[i], err = k8s.CheckSAPermissions(k8s.CheckSAPermissionsParams{Verb: "create", Resource: resource, Print: true, Namespace: namespace}, kubeconfig)
		if err != nil {
			utils.Red.Println(err)
		}
	}

	for _, pem := range pems {
		if !pem {
			utils.Red.Println("\n🚫 You don't have sufficient permissions.\n🙄 Please use a service account with sufficient permissions.")
			os.Exit(1)
		}
	}

	utils.White_B.Println("\n🌟 Sufficient permissions. Installing the Chaos Delegate...")
}

// Summary display the agent details based on input
func Summary(agent types.Agent, kubeconfig *string) {
	utils.White_B.Printf("\n📌 Summary \nChaos Delegate Name: %s\nChaos Delegate Description: %s\nChaos Delegate SSL/TLS Skip: %t\nPlatform Name: %s\n", agent.AgentName, agent.Description, agent.SkipSSL, agent.PlatformName)
	if ok, _ := k8s.NsExists(agent.Namespace, kubeconfig); ok {
		utils.White_B.Println("Namespace: ", agent.Namespace)
	} else {
		utils.White_B.Println("Namespace: ", agent.Namespace, "(new)")
	}

	if k8s.SAExists(k8s.SAExistsParams{Namespace: agent.Namespace, Serviceaccount: agent.ServiceAccount}, kubeconfig) {
		utils.White_B.Println("Service Account: ", agent.ServiceAccount)
	} else {
		utils.White_B.Println("Service Account: ", agent.ServiceAccount, "(new)")
	}

	utils.White_B.Printf("\nInstallation Mode: %s\n", agent.Mode)
}

func ConfirmInstallation() {
	var descision string
	utils.White_B.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&descision)

	if strings.ToLower(descision) == "yes" || strings.ToLower(descision) == "y" {
		utils.White_B.Println("👍 Continuing Chaos Delegate connection!!")
	} else {
		utils.Red.Println("✋ Exiting Chaos Delegate connection!!")
		os.Exit(1)
	}
}

func CreateRandomProject(cred types.Credentials) string {
	rand, err := utils.GenerateRandomString(10)
	utils.PrintError(err)

	projectName := cred.Username + "-" + rand

	project, err := apis.CreateProjectRequest(projectName, cred)
	utils.PrintError(err)

	return project.Data.ID
}
