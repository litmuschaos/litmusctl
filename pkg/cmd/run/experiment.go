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
package run

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

// experimentCmd represents the project command
var experimentCmd = &cobra.Command{
	Use: "chaos-experiment",
	Short: `Create a Chaos Experiment
	Example:
	#Save a Chaos Experiment
	litmusctl run chaos-experiment --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --experiment-id="1c9c5801-8789-4ac9-bf5f-32649b707a5c"

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

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

		eid, err := cmd.Flags().GetString("experiment-id")
		utils.PrintError(err)

		// Handle blank input for Chaos Experiment ID
		if eid == "" {
			utils.White_B.Print("\nEnter the Chaos Experiment ID: ")
			fmt.Scanln(&eid)

			if eid == "" {
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
			if p.ID == pid {
				project = p
			}
		}
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == utils.MemberOwnerRole || member.Role == utils.MemberEditorRole) {
				editAccess = true
			}
		}
		if !editAccess {
			utils.Red.Println("⛔ User doesn't have edit access to the project!!")
			os.Exit(1)
		}

		// Make API call
		runExperiment, err := experiment.RunExperiment(pid, eid, credentials)
		if err != nil {
			if (runExperiment.Data == experiment.RunExperimentData{}) {
				if strings.Contains(err.Error(), "multiple run errors") {
					utils.Red.Println("\n❌ Chaos Experiment already exists")
					os.Exit(1)
				}
				if strings.Contains(err.Error(), "no documents in result") {
					utils.Red.Println("❌ The specified Project ID or Chaos Infrastructure ID doesn't exist.")
					os.Exit(1)
				} else {
					utils.White_B.Print("\n❌ Failed to run chaos experiment: " + err.Error())
					os.Exit(1)
				}
			}
		}

		//Successful run
		utils.White_B.Println("\n🚀 Chaos Experiment running successfully 🎉")

	},
}

func init() {
	RunCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to create Chaos Experiment for the particular project. To see the projects, apply litmusctl get projects")
	experimentCmd.Flags().String("experiment-id", "", "Set the environment-id to create Chaos Experiment for the particular Chaos Infrastructure. To see the Chaos Infrastructures, apply litmusctl get chaos-infra")
}
