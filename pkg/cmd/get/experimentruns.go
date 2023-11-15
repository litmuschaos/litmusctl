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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// experimentRunsCmd represents the Chaos Experiments runs command
var experimentRunsCmd = &cobra.Command{
	Use:   "chaos-experiment-runs",
	Short: "Display list of Chaos Experiments runs within the project",
	Long:  `Display list of Chaos Experiments runs within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var listExperimentRunsRequest model.ListExperimentRunRequest
		var projectID string

		projectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			for projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		//  experiment ID flag
		experimentID, err := cmd.Flags().GetString("experiment-id")
		utils.PrintError(err)
		if experimentID != "" {
			listExperimentRunsRequest.ExperimentIDs = []*string{&experimentID}

		}

		// experiment run ID flag
		experimentRunID, err := cmd.Flags().GetString("experiment-run-id")
		utils.PrintError(err)
		if experimentRunID != "" {
			listExperimentRunsRequest.ExperimentRunIDs = []*string{&experimentRunID}

		}

		listAllExperimentRuns, _ := cmd.Flags().GetBool("all")
		if !listAllExperimentRuns {
			listExperimentRunsRequest.Pagination = &model.Pagination{}
			listExperimentRunsRequest.Pagination.Limit, _ = cmd.Flags().GetInt("count")
		}

		experimentRuns, err := experiment.GetExperimentRunsList(projectID, listExperimentRunsRequest, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ The specified Project ID doesn't exist.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(experimentRuns.Data)

		case "yaml":
			utils.PrintInYamlFormat(experimentRuns.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "CHAOS EXPERIMENT RUN ID\tSTATUS\tRESILIENCY SCORE\tCHAOS EXPERIMENT ID\tCHAOS EXPERIMENT NAME\tTARGET CHAOS INFRA\tUPDATED AT\tUPDATED BY")

			for _, experimentRun := range experimentRuns.Data.ListExperimentRunDetails.ExperimentRuns {

				var lastUpdated string
				unixSecondsInt, err := strconv.ParseInt(experimentRun.UpdatedAt, 10, 64)
				if err != nil {
					lastUpdated = "None"
				} else {
					lastUpdated = time.Unix(unixSecondsInt, 0).Format("January 2 2006, 03:04:05 pm")
				}

				utils.White.Fprintln(
					writer,
					experimentRun.ExperimentRunID+"\t"+experimentRun.Phase.String()+"\t"+strconv.FormatFloat(*experimentRun.ResiliencyScore, 'f', 2, 64)+"\t"+experimentRun.ExperimentID+"\t"+experimentRun.ExperimentName+"\t"+experimentRun.Infra.Name+"\t"+lastUpdated+"\t"+experimentRun.UpdatedBy.Username)
			}

			if listAllExperimentRuns || (experimentRuns.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns <= listExperimentRunsRequest.Pagination.Limit) {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Experiment runs", experimentRuns.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns, experimentRuns.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns))
			} else {
				utils.White_B.Fprintln(writer, fmt.Sprintf("\nShowing %d of %d Chaos Experiment runs", listExperimentRunsRequest.Pagination.Limit, experimentRuns.Data.ListExperimentRunDetails.TotalNoOfExperimentRuns))
			}

			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(experimentRunsCmd)

	experimentRunsCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Experiments from the particular project. To see the projects, apply litmusctl get projects")
	experimentRunsCmd.Flags().Int("count", 30, "Set the count of Chaos Experiments runs to display. Default value is 30")
	experimentRunsCmd.Flags().BoolP("all", "A", false, "Set to true to display all Chaos Experiments runs")

	experimentRunsCmd.Flags().String("experiment-id", "", "Set the experiment ID to list experiment runs within a specific experiment")
	experimentRunsCmd.Flags().String("experiment-run-id", "", "Set the experiment run ID to list a specific experiment run")

	experimentRunsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
