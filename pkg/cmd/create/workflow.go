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
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"

	"github.com/spf13/cobra"
)

// workflowCmd represents the project command
var workflowCmd = &cobra.Command{
	Use: "chaos-scenario",
	Short: `Create a Chaos Scenario
	Example:
	#create a Chaos Scenario
	litmusctl create chaos-scenario -f chaos-scenario.yaml --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --chaos-delegate-id="1c9c5801-8789-4ac9-bf5f-32649b707a5c"

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var chaosWorkFlowRequest model.ChaosWorkFlowRequest

		workflowManifest, err := cmd.Flags().GetString("file")
		utils.PrintError(err)

		chaosWorkFlowRequest.ProjectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID
		if chaosWorkFlowRequest.ProjectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&chaosWorkFlowRequest.ProjectID)

			if chaosWorkFlowRequest.ProjectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		chaosWorkFlowRequest.ClusterID, err = cmd.Flags().GetString("chaos-delegate-id")
		utils.PrintError(err)

		// Handle blank input for Chaos Delegate ID
		if chaosWorkFlowRequest.ClusterID == "" {
			utils.White_B.Print("\nEnter the Chaos Delegate ID: ")
			fmt.Scanln(&chaosWorkFlowRequest.ClusterID)

			if chaosWorkFlowRequest.ClusterID == "" {
				utils.Red.Println("⛔ Chaos Delegate ID can't be empty!!")
				os.Exit(1)
			}
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == chaosWorkFlowRequest.ProjectID {
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

		// Parse workflow manifest and populate chaosWorkFlowInput
		err = utils.ParseWorkflowManifest(workflowManifest, &chaosWorkFlowRequest)
		if err != nil {
			utils.Red.Println("❌ Error parsing Chaos Scenario manifest: " + err.Error())
			os.Exit(1)
		}

		// Make API call
		createdWorkflow, err := apis.CreateWorkflow(chaosWorkFlowRequest, credentials)
		if err != nil {
			if (createdWorkflow.Data == apis.CreatedChaosWorkflow{}) {
				if strings.Contains(err.Error(), "multiple write errors") {
					utils.Red.Println("\n❌ Chaos Scenario/" + chaosWorkFlowRequest.WorkflowName + " already exists")
					os.Exit(1)
				} else {
					utils.White_B.Print("\n❌ Chaos Scenario/" + chaosWorkFlowRequest.WorkflowName + " failed to be created: " + err.Error())
					os.Exit(1)
				}
			}
		}

		// Successful creation
		utils.White_B.Println("\n🚀 Chaos Scenario/" + createdWorkflow.Data.CreateChaosWorkflow.WorkflowName + " successfully created 🎉")
		if createdWorkflow.Data.CreateChaosWorkflow.CronSyntax == "" {
			utils.White_B.Println("\nThe next run of this Chaos Scenario will be scheduled immediately.")
		} else {
			utils.White_B.Println(
				"\nThe next run of this Chaos Scenario will be scheduled at " +
					cronexpr.MustParse(createdWorkflow.Data.CreateChaosWorkflow.CronSyntax).Next(time.Now()).Format("January 2nd 2006, 03:04:05 pm"))
		}
	},
}

func init() {
	CreateCmd.AddCommand(workflowCmd)

	workflowCmd.Flags().String("project-id", "", "Set the project-id to create Chaos Scenario for the particular project. To see the projects, apply litmusctl get projects")
	workflowCmd.Flags().String("chaos-delegate-id", "", "Set the chaos-delegate-id to create Chaos Scenario for the particular Chaos Delegate. To see the Chaos Delegates, apply litmusctl get chaos-delegates")
	workflowCmd.Flags().StringP("file", "f", "", "The manifest file for the Chaos Scenario")
}
