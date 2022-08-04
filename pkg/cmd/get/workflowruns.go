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
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// workflowRunsCmd represents the Chaos Scenario runs command
var workflowRunsCmd = &cobra.Command{
	Use:   "chaos-scenario-runs",
	Short: "Display list of Chaos Scenario runs within the project",
	Long:  `Display list of Chaos Scenario runs within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var listWorkflowRunsRequest model.ListWorkflowRunsRequest

		listWorkflowRunsRequest.ProjectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if listWorkflowRunsRequest.ProjectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&listWorkflowRunsRequest.ProjectID)

			for listWorkflowRunsRequest.ProjectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		listAllWorkflowRuns, _ := cmd.Flags().GetBool("all")
		if !listAllWorkflowRuns {
			listWorkflowRunsRequest.Pagination = &model.Pagination{}
			listWorkflowRunsRequest.Pagination.Limit, _ = cmd.Flags().GetInt("count")
		}

		workflowRuns, err := apis.GetWorkflowRunsList(listWorkflowRunsRequest, credentials)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(workflowRuns.Data)

		case "yaml":
			utils.PrintInYamlFormat(workflowRuns.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "CHAOS SCENARIO RUN ID\tSTATUS\tRESILIENCY SCORE\tCHAOS SCENARIO ID\tCHAOS SCENARIO NAME\tTARGET CHAOS DELEGATE\tLAST RUN\tEXECUTED BY")

			for _, workflowRun := range workflowRuns.Data.ListWorkflowRunsDetails.WorkflowRuns {

				var lastUpdated string
				unixSecondsInt, err := strconv.ParseInt(workflowRun.LastUpdated, 10, 64)
				if err != nil {
					lastUpdated = "None"
				} else {
					lastUpdated = time.Unix(unixSecondsInt, 0).Format("January 2 2006, 03:04:05 pm")
				}

				utils.White.Fprintln(
					writer,
					workflowRun.WorkflowRunID+"\t"+workflowRun.Phase+"\t"+strconv.FormatFloat(*workflowRun.ResiliencyScore, 'f', 2, 64)+"\t"+workflowRun.WorkflowID+"\t"+workflowRun.WorkflowName+"\t"+workflowRun.ClusterName+"\t"+lastUpdated+"\t"+workflowRun.ExecutedBy)
			}

			if listAllWorkflowRuns || (workflowRuns.Data.ListWorkflowRunsDetails.TotalNoOfWorkflowRuns <= listWorkflowRunsRequest.Pagination.Limit) {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Scenario runs", workflowRuns.Data.ListWorkflowRunsDetails.TotalNoOfWorkflowRuns, workflowRuns.Data.ListWorkflowRunsDetails.TotalNoOfWorkflowRuns))
			} else {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Scenario runs", listWorkflowRunsRequest.Pagination.Limit, workflowRuns.Data.ListWorkflowRunsDetails.TotalNoOfWorkflowRuns))
			}

			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(workflowRunsCmd)

	workflowRunsCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Scenarios from the particular project. To see the projects, apply litmusctl get projects")
	workflowRunsCmd.Flags().Int("count", 30, "Set the count of Chaos Scenario runs to display. Default value is 30")
	workflowRunsCmd.Flags().BoolP("all", "A", false, "Set to true to display all Chaos Scenario runs")

	workflowRunsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
