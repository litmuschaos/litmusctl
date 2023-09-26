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
	"os"
	"strings"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// experimentCmd represents the project command
var experimentCmd = &cobra.Command{
	Use: "chaos-experiment",
	Short: `Create a Chaos Experiment
	Example:
	#create a Chaos Experiment
	litmusctl create chaos-experiment -f chaos-experiment.yaml --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --chaos-infra-id="1c9c5801-8789-4ac9-bf5f-32649b707a5c"

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var chaosExperimentRequest models.SaveChaosExperimentRequest

		workflowManifest, err := cmd.Flags().GetString("file")
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

		chaosExperimentRequest.InfraID, err = cmd.Flags().GetString("chaos-infra-id")
		utils.PrintError(err)

		// Handle blank input for Chaos Infra ID
		if chaosExperimentRequest.InfraID == "" {
			utils.White_B.Print("\nEnter the Chaos Infra ID: ")
			fmt.Scanln(&chaosExperimentRequest.InfraID)

			if chaosExperimentRequest.InfraID == "" {
				utils.Red.Println("⛔ Chaos Infra ID can't be empty!!")
				os.Exit(1)
			}
		}

		chaosExperimentRequest.Description, err = cmd.Flags().GetString("description")
		utils.PrintError(err)
		if chaosExperimentRequest.Description == "" {
			utils.White_B.Print("\nExperiment Description: ")
			fmt.Scanln(&chaosExperimentRequest.Description)
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

		// Parse experiment manifest and populate chaosExperimentInput
		err = utils.ParseExperimentManifest(workflowManifest, &chaosExperimentRequest)
		if err != nil {
			utils.Red.Println("❌ Error parsing Chaos Experiment manifest: " + err.Error())
			os.Exit(1)
		}
		// Generate ExperimentID from ExperimentName
		chaosExperimentRequest.ID = utils.GenerateNameID(chaosExperimentRequest.Name)
		// Make API call
		createExperiment, err := experiment.CreateExperiment(pid, chaosExperimentRequest, credentials)
		if err != nil {
			if (createExperiment.Data == experiment.RunExperimentData{}) {
				if strings.Contains(err.Error(), "multiple write errors") {
					utils.Red.Println("\n❌ Chaos Experiment/" + chaosExperimentRequest.Name + " already exists")
					os.Exit(1)
				}
				if strings.Contains(err.Error(), "no documents in result") {
					utils.Red.Println("❌ The specified Project ID or Chaos Infrastructure ID doesn't exist.")
					os.Exit(1)
				}
				if strings.Contains(err.Error(), "multiple run errors") {
					utils.Red.Println("\n❌ Chaos Experiment already exists")
					os.Exit(1)
				}
			} else {
				utils.White_B.Print("\n❌ Chaos Experiment/" + chaosExperimentRequest.Name + " failed to be created: " + err.Error())
				os.Exit(1)
			}
		}

		//Successful creation
		utils.White_B.Println("\n🚀 Chaos Experiment successfully created and experiment-run is scheduled 🎉")
	},
}

func init() {
	CreateCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to create Chaos Experiment for the particular project. To see the projects, apply litmusctl get projects")
	experimentCmd.Flags().String("chaos-infra-id", "", "Set the chaos-infra-id to create Chaos Experiment for the particular Chaos Infrastructure. To see the Chaos Infrastructures, apply litmusctl get chaos-infra")
	experimentCmd.Flags().StringP("file", "f", "", "The manifest file for the Chaos Experiment")
	experimentCmd.Flags().StringP("description", "d", "", "The Description for the Chaos Experiment")
}
