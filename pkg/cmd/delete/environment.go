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
package delete

import (
	"os"

	"github.com/litmuschaos/litmusctl/pkg/apis/environment"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"strings"
)

// experimentCmd represents the Chaos Experiment command
var environmentCmd = &cobra.Command{
	Use: "chaos-environment",
	Short: `Delete a Chaos environment
Example:
#delete a Chaos Environment
litmusctl delete chaos-environment --project-id=8adf62d5-64f8-4c66-ab53-63729db9dd9a --environment-id=environmentexample

Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
`,

	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		environmentID, err := cmd.Flags().GetString("environment-id")
		utils.PrintError(err)

		// Handle blank input for project ID

		if projectID == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Project ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.Red.Println("⛔ Error:", err)
				os.Exit(1)
			}
			projectID = result
		}

		if environmentID == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Environment ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.Red.Println("⛔ Error:", err)
				os.Exit(1)
			}
			environmentID = result
		}

		// Handle blank input for Chaos Environment ID
		if environmentID == "" {
			utils.Red.Println("⛔ Chaos Environment ID can't be empty!!")
			os.Exit(1)
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == projectID {
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

		environmentGet, err := environment.GetChaosEnvironment(projectID, environmentID, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ You don't have enough permissions to delete an environment.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}
		environmentGetData := environmentGet.Data.EnvironmentDetails
		if len(environmentGetData.InfraIDs) > 0 {
			utils.Red.Println("Chaos Infras present in the Chaos Environment" +
				", delete the Chaos Infras first to delete the Environment")
		    os.Exit(1)
		}

		// confirm before deletion
		prompt := promptui.Prompt{
			Label:     "Are you sure you want to delete this Chaos Environment? (y/n)",
			AllowEdit: true,
		}

		result, err := prompt.Run()
		if err != nil {
			utils.Red.Println("⛔ Error:", err)
			os.Exit(1)
		}

		if result != "y" {
			utils.White_B.Println("\n❌ Chaos Environment was not deleted.")
			os.Exit(0)
		}

		// Make API call
		_, err = environment.DeleteEnvironment(projectID, environmentID, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Error in deleting Chaos Environment: ", err.Error())
			os.Exit(1)
		}

		utils.White_B.Println("\n🚀 Chaos Environment successfully deleted.")

	},
}

func init() {
	DeleteCmd.AddCommand(environmentCmd)

	environmentCmd.Flags().String("project-id", "", "Set the project-id to delete Chaos Environment for the particular project. To see the projects, apply litmusctl get projects")
	environmentCmd.Flags().String("environment-id", "", "Set the environment-id to delete the particular Chaos Environment.")
}
