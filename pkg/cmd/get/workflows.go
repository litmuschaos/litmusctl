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

// workflowsCmd represents the workflows command
var workflowsCmd = &cobra.Command{
	Use:   "workflows",
	Short: "Display list of workflows within the project",
	Long:  `Display list of workflows within the project`,
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
		agentName, err := cmd.Flags().GetString("agent")
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
			utils.White_B.Fprintln(writer, "WORKFLOW ID\tWORKFLOW NAME\tWORKFLOW TYPE\tNEXT SCHEDULE\tAGENT ID\tAGENT NAME\tLAST UPDATED BY")

			for _, workflow := range workflows.Data.ListWorkflowDetails.Workflows {
				if workflow.CronSyntax != "" {
					utils.White.Fprintln(
						writer,
						workflow.WorkflowID+"\t"+workflow.WorkflowName+"\tCron Workflow\t"+cronexpr.MustParse(workflow.CronSyntax).Next(time.Now()).Format("January 2 2006, 03:04:05 pm")+"\t"+workflow.ClusterID+"\t"+workflow.ClusterName+"\t"+*workflow.LastUpdatedBy)
				} else {
					utils.White.Fprintln(
						writer,
						workflow.WorkflowID+"\t"+workflow.WorkflowName+"\tNon Cron Workflow\tNone\t"+workflow.ClusterID+"\t"+workflow.ClusterName+"\t"+*workflow.LastUpdatedBy)
				}
			}

			if listAllWorkflows || (workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows <= listWorkflowsRequest.Pagination.Limit) {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d workflows", workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows, workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows))
			} else {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d workflows", listWorkflowsRequest.Pagination.Limit, workflows.Data.ListWorkflowDetails.TotalNoOfWorkflows))
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(workflowsCmd)

	workflowsCmd.Flags().String("project-id", "", "Set the project-id to list workflows from the particular project. To see the projects, apply litmusctl get projects")
	workflowsCmd.Flags().Int("count", 30, "Set the count of workflows to display. Default value is 30")
	workflowsCmd.Flags().Bool("all", false, "Set to true to display all workflows")
	workflowsCmd.Flags().StringP("agent", "A", "", "Set the agent name to display all workflows targeted towards that particular agent.")

	workflowsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
