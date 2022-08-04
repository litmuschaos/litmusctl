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
package get

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// workflowsCmd represents the Chaos Scenarios command
var workflowsCmd = &cobra.Command{
	Use:   "chaos-scenarios",
	Short: "Display list of Chaos Scenarios within the project",
	Long:  `Display list of Chaos Scenarios within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var listWorkflowsRequest model.ListWorkflowsRequest

		listWorkflowsRequest.ProjectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if listWorkflowsRequest.ProjectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&listWorkflowsRequest.ProjectID)

			for listWorkflowsRequest.ProjectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		listAllWorkflows, _ := cmd.Flags().GetBool("all")
		if !listAllWorkflows {
			listWorkflowsRequest.Pagination = &model.Pagination{}
			listWorkflowsRequest.Pagination.Limit, _ = cmd.Flags().GetInt("count")
		}

		listWorkflowsRequest.Filter = &model.WorkflowFilterInput{}
		agentName, err := cmd.Flags().GetString("chaos-delegate")
		utils.PrintError(err)
		listWorkflowsRequest.Filter.ClusterName = &agentName

		workflows, err := apis.GetWorkflowList(listWorkflowsRequest, credentials)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(workflows.Data)

		case "yaml":
			utils.PrintInYamlFormat(workflows.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "CHAOS SCENARIO ID\tCHAOS SCENARIO NAME\tCHAOS SCENARIO TYPE\tNEXT SCHEDULE\tCHAOS DELEGATE ID\tCHAOS DELEGATE NAME\tLAST UPDATED BY")

			for _, workflow := range workflows.Data.ListWorkflowDetails.Workflows {
				if workflow.CronSyntax != "" {
					utils.White.Fprintln(
						writer,
						workflow.WorkflowID+"\t"+workflow.WorkflowName+"\tCron Chaos Scenario\t"+cronexpr.MustParse(workflow.CronSyntax).Next(time.Now()).Format("January 2 2006, 03:04:05 pm")+"\t"+workflow.ClusterID+"\t"+workflow.ClusterName+"\t"+*workflow.LastUpdatedBy)
				} else {
					utils.White.Fprintln(
						writer,
						workflow.WorkflowID+"\t"+workflow.WorkflowName+"\tNon Cron Chaos Scenario\tNone\t"+workflow.ClusterID+"\t"+workflow.ClusterName+"\t"+*workflow.LastUpdatedBy)
				}
			}

			if listAllWorkflows || (workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows <= listWorkflowsRequest.Pagination.Limit) {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Scenarios", workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows, workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows))
			} else {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Scenarios", listWorkflowsRequest.Pagination.Limit, workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows))
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(workflowsCmd)

	workflowsCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Scenarios from the particular project. To see the projects, apply litmusctl get projects")
	workflowsCmd.Flags().Int("count", 30, "Set the count of Chaos Scenarios to display. Default value is 30")
	workflowsCmd.Flags().Bool("all", false, "Set to true to display all Chaos Scenarios")
	workflowsCmd.Flags().StringP("chaos-delegate", "A", "", "Set the Chaos Delegate name to display all Chaos Scenarios targeted towards that particular Chaos Delegate.")

	workflowsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
