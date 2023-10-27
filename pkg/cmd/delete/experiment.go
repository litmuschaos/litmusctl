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
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/completion"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// experimentCmd represents the Chaos Experiment command
var experimentCmd = &cobra.Command{
	Use: "chaos-experiment",
	Short: `Delete a Chaos experiment
	Example:
	#delete a Chaos Experiment
	litmusctl delete chaos-experiment c520650e-7cb6-474c-b0f0-4df07b2b025b --project-id=c520650e-7cb6-474c-b0f0-4df07b2b025b

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	ValidArgsFunction: completion.ExperimentIDCompletion,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID
		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			if projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		experimentID := args[0]

		// Handle blank input for Chaos Experiment ID
		if experimentID == "" {
			utils.White_B.Print("\nEnter the Chaos Experiment ID: ")
			fmt.Scanln(&experimentID)

			if experimentID == "" {
				utils.Red.Println("⛔ Chaos Experiment ID can't be empty!!")
				os.Exit(1)
			}
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

		// Make API call
		deleteExperiment, err := experiment.DeleteChaosExperiment(projectID, &experimentID, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Error in deleting Chaos Experiment: ", err.Error())
			os.Exit(1)
		}

		if deleteExperiment.Data.IsDeleted {
			utils.White_B.Println("\n🚀 Chaos Experiment successfully deleted.")
		} else {
			utils.White_B.Println("\n❌ Failed to delete Chaos Experiment. Please check if the ID is correct or not.")
		}
	},
}

func init() {
	DeleteCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to create Chaos Experiment for the particular project. To see the projects, apply litmusctl get projects")
	experimentCmd.RegisterFlagCompletionFunc("project-id", completion.ProjectIDFlagCompletion)

}
