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
package infra_ops

import (
	"fmt"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis/environment"
	"github.com/litmuschaos/litmusctl/pkg/apis/infrastructure"
	"os"
	"strconv"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

func PrintExistingInfra(infra infrastructure.InfraData) {
	utils.Red.Println("\nChaos Infrastructure with the given name already exists.")
	// Print Chaos Infra list if existing Chaos Infra name is entered twice
	utils.White_B.Println("\nConnected Chaos Infrastructure list:")

	for i := range infra.Data.ListInfraDetails.Infras {
		utils.White_B.Println("-", infra.Data.ListInfraDetails.Infras[i].Name)
	}

	utils.White_B.Println("\n❗ Please enter a different name.")
}

func PrintExistingEnvironments(env environment.ListEnvironmentData) {
	// Print Chaos EnvironmentID list if Given ID doesn't exist
	utils.White_B.Println("\nExisting Chaos Environments list:")

	for i := range env.Data.ListEnvironmentDetails.Environments {
		utils.White_B.Println("-", env.Data.ListEnvironmentDetails.Environments[i].EnvironmentID)
	}
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

// GetModeType gets mode of Chaos Infrastructure installation as input
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
func GetInfraDetails(mode string, pid string, c types.Credentials, kubeconfig *string) (types.Infra, error) {
	var newInfra types.Infra
	// Get Infra name as input
	utils.White_B.Println("\nEnter the details of the Chaos Infrastructure")
	// Label for goto statement in case of invalid Chaos Infra name

INFRA_NAME:
	utils.White_B.Print("\nChaos Infra Name: ")
	newInfra.InfraName = utils.Scanner()
	if newInfra.InfraName == "" {
		utils.Red.Println("⛔ Chaos Infra name cannot be empty. Please enter a valid name.")
		goto INFRA_NAME
	}

	// Check if Chaos Infra with the given name already exists
	isInfraExist, err, infra := ValidateInfraNameExists(newInfra.InfraName, pid, c)

	if isInfraExist {
		PrintExistingInfra(infra)
		goto INFRA_NAME
	}

	// Get Infra description as input
	utils.White_B.Print("\nChaos Infrastructure Description: ")
	newInfra.Description = utils.Scanner()

ENVIRONMENT:
	utils.White_B.Print("\nChaos EnvironmentID: ")
	newInfra.EnvironmentID = utils.Scanner()

	if newInfra.EnvironmentID == "" {
		utils.Red.Println("⛔ Chaos Environment ID cannot be empty. Please enter a valid Environment.")
		goto ENVIRONMENT
	}

	// Check if Chaos Environment with the given name exists
	Env, err := environment.GetEnvironmentList(pid, c)
	if err != nil {
		return types.Infra{}, err
	}

	var isEnvExist = false
	for i := range Env.Data.ListEnvironmentDetails.Environments {
		if newInfra.EnvironmentID == Env.Data.ListEnvironmentDetails.Environments[i].EnvironmentID {
			utils.White_B.Println(Env.Data.ListEnvironmentDetails.Environments[i].EnvironmentID)
			isEnvExist = true
			break
		}
	}

	if !isEnvExist {
		utils.Red.Println("\nChaos Environment with the given ID doesn't exists.")
		PrintExistingEnvironments(Env)
		utils.White_B.Println("\n❗ Please enter a name from the List.")
		goto ENVIRONMENT
	}

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
				str += "\"operator\" : \"" + operator + "\" ,"
			}

			utils.White_B.Print("\nEffect: ")
			effect := utils.Scanner()

			if effect != "" {
				str += "\"effect\": \"" + effect + "\" ,"
			}

			if ts != "" {
				// check whether if effect is "NoSchedule" then tolerationsSeconds should be 0
				if effect != "NoSchedule" {
					str += "\"tolerationSeconds\": " + ts + " ,"
				}
			}

			utils.White_B.Print("\nKey: ")
			key := utils.Scanner()
			if key != "" {
				str += "\"key\": \"" + key + "\" ,"
			}

			utils.White_B.Print("\nValue: ")
			value := utils.Scanner()
			if key != "" {
				str += "\"value\": \"" + value + "\""
			}

			str += " },"
		}
		if nts > 0 {
			str = str[:len(str)-1]
		}
		str += "]"

		newInfra.Tolerations = str
	}

	// Get platform name as input
	newInfra.PlatformName = GetPlatformName(kubeconfig)
	// Set Infra type
	newInfra.InfraType = utils.InfraTypeExternal
	// Set project id
	newInfra.ProjectId = pid
	// Get namespace
	newInfra.Namespace, newInfra.NsExists = k8s.ValidNs(mode, utils.ChaosInfraLabel, kubeconfig)

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

	utils.White_B.Println("\n🌟 Sufficient permissions. Installing the Chaos Infra...")
}

// ValidateInfraNameExists checks if an infrastructure already exists
func ValidateInfraNameExists(infraName string, pid string, c types.Credentials) (bool, error, infrastructure.InfraData) {
	infra, err := infrastructure.GetInfraList(c, pid, model.ListInfraRequest{})
	if err != nil {
		return false, err, infrastructure.InfraData{}
	}

	for i := range infra.Data.ListInfraDetails.Infras {
		if infraName == infra.Data.ListInfraDetails.Infras[i].Name {
			utils.White_B.Println(infra.Data.ListInfraDetails.Infras[i].Name)
			return true, nil, infra
		}
	}
	return false, nil, infrastructure.InfraData{}
}

// Summary display the Infra details based on input
func Summary(infra types.Infra, kubeconfig *string) {
	utils.White_B.Printf("\n📌 Summary \nChaos Infra Name: %s\nChaos EnvironmentID: %s\nChaos Infra Description: %s\nChaos Infra SSL/TLS Skip: %t\nPlatform Name: %s\n", infra.InfraName, infra.EnvironmentID, infra.Description, infra.SkipSSL, infra.PlatformName)
	if ok, _ := k8s.NsExists(infra.Namespace, kubeconfig); ok {
		utils.White_B.Println("Namespace: ", infra.Namespace)
	} else {
		utils.White_B.Println("Namespace: ", infra.Namespace, "(new)")
	}

	if k8s.SAExists(k8s.SAExistsParams{Namespace: infra.Namespace, Serviceaccount: infra.ServiceAccount}, kubeconfig) {
		utils.White_B.Println("Service Account: ", infra.ServiceAccount)
	} else {
		utils.White_B.Println("Service Account: ", infra.ServiceAccount, "(new)")
	}

	utils.White_B.Printf("\nInstallation Mode: %s\n", infra.Mode)
}

func ConfirmInstallation() {
	var descision string
	utils.White_B.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&descision)

	if strings.ToLower(descision) == "yes" || strings.ToLower(descision) == "y" {
		utils.White_B.Println("👍 Continuing Chaos Infrastructure connection!!")
	} else {
		utils.Red.Println("✋ Exiting Chaos Infrastructure connection!!")
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
