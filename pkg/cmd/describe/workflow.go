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
package describe

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// workflowCmd represents the Chaos Scenario command
var workflowCmd = &cobra.Command{
	Use:   "chaos-scenario",
	Short: "Describe a Chaos Scenario within the project",
	Long:  `Describe a Chaos Scenario within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var describeWorkflowRequest model.ListWorkflowsRequest

		describeWorkflowRequest.ProjectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if describeWorkflowRequest.ProjectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&describeWorkflowRequest.ProjectID)

			for describeWorkflowRequest.ProjectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		var workflowID string
		if len(args) == 0 {
			utils.White_B.Print("\nEnter the Chaos Scenario ID: ")
			fmt.Scanln(&workflowID)
		} else {
			workflowID = args[0]
		}

		// Handle blank input for Chaos Scenario ID
		if workflowID == "" {
			utils.Red.Println("⛔ Chaos Scenario ID can't be empty!!")
			os.Exit(1)
		}

		describeWorkflowRequest.WorkflowIDs = append(describeWorkflowRequest.WorkflowIDs, &workflowID)

		workflow, err := apis.GetWorkflowList(describeWorkflowRequest, credentials)
		utils.PrintError(err)

		if len(workflow.Data.ListWorkflowDetails.Workflows) == 0 {
			utils.Red.Println("⛔ No chaos scenarios found with ID: ", workflowID)
			os.Exit(1)
		}

		yamlManifest, err := yaml.JSONToYAML([]byte(workflow.Data.ListWorkflowDetails.Workflows[0].WorkflowManifest))
		if err != nil {
			utils.Red.Println("❌ Error parsing Chaos Scenario manifest: " + err.Error())
			os.Exit(1)
		}
		utils.PrintInYamlFormat(string(yamlManifest))
	},
}

func init() {
	DescribeCmd.AddCommand(workflowCmd)

	workflowCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Scenarios from the particular project. To see the projects, apply litmusctl get projects")
}
