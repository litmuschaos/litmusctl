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
package create

import (
	"fmt"
	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/environment"
	"github.com/litmuschaos/litmusctl/pkg/infra_ops"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)

// environmentCmd represents the Chaos infra command
var environmentCmd = &cobra.Command{
	Use: "chaos-environment",
	Short: `Create an Environment.
	Example(s):

	#create a Chaos Environment
	litmusctl create chaos-environment --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --name="new-chaos-environment"

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var newEnvironment models.CreateEnvironmentRequest

		pid, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID
		if pid == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&pid)

			if pid == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		envName, err := cmd.Flags().GetString("name")
		utils.PrintError(err)

		// Handle blank input for project ID
		if envName == "" {
			utils.White_B.Print("\nEnter the Environment Name: ")
			fmt.Scanln(&envName)

			if envName == "" {
				utils.Red.Println("⛔ Environment Name can't be empty!!")
				os.Exit(1)
			}
		}
		newEnvironment.Name = envName

		description, err := cmd.Flags().GetString("description")
		utils.PrintError(err)
		newEnvironment.Description = &description

		envType, err := cmd.Flags().GetString("type")
		utils.PrintError(err)

		// Handle blank input for project ID
		if envType == "" {
			utils.White_B.Print("\nEnter the Environment Type: ")
			fmt.Scanln(&envType)

			if envType == "" {
				utils.Red.Println("⛔ Environment Type can't be empty!!")
				os.Exit(1)
			}
		}
		newEnvironment.Type = models.EnvironmentType(envType)

		envs, err := environment.GetEnvironmentList(pid, credentials)
		utils.PrintError(err)

		// Generate EnvironmentID from Environment Name
		envID := utils.GenerateNameID(newEnvironment.Name)
		newEnvironment.EnvironmentID = envID

		// Check if Environment exists
		var isEnvExist = false
		for i := range envs.Data.ListEnvironmentDetails.Environments {
			if envID == envs.Data.ListEnvironmentDetails.Environments[i].EnvironmentID {
				utils.White_B.Print(envs.Data.ListEnvironmentDetails.Environments[i].EnvironmentID)
				isEnvExist = true
				break
			}
		}
		if isEnvExist {
			utils.Red.Println("\nChaos Environment with the given ID already exists, try with a different name")
			infra_ops.PrintExistingEnvironments(envs)
			os.Exit(1)
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == pid {
				project = p
			}
		}
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == "Owner" || member.Role == "Editor") {
				editAccess = true
			}
		}
		if !editAccess {
			utils.Red.Println("⛔ User doesn't have edit access to the project!!")
			os.Exit(1)
		}

		newEnv, err := environment.CreateEnvironment(pid, newEnvironment, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Chaos Environment connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}
		//TODO: add the nil checker for the response(newEnv.Data)
		//Print error message in case Data field is null in response
		//if (newEnv.Data == environment.CreateEnvironmentData{}) {
		//	utils.White_B.Print("\n🚫 Chaos newInfra connection failed: " + newEnv.Errors[0].Message + "\n")
		//	os.Exit(1)
		//}
		utils.White_B.Println("\n🚀 New Chaos Environment creation successful!! 🎉")
		utils.White_B.Println("EnvironmentID: " + newEnv.Data.EnvironmentDetails.EnvironmentID)
	},
}

func init() {
	CreateCmd.AddCommand(environmentCmd)
	environmentCmd.Flags().String("project-id", "", "Set the project-id to install Chaos infra for the particular project. To see the projects, apply litmusctl get projects")
	environmentCmd.Flags().String("type", "NON_PROD", "Set the installation mode for the kind of Chaos infra | Supported=cluster/namespace")
	environmentCmd.Flags().String("name", "", "Set the Chaos infra name")
	environmentCmd.Flags().String("description", "---", "Set the Chaos infra description")
}
