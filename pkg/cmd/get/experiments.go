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
	"strings"
	"text/tabwriter"
	"time"

	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"

	"github.com/gorhill/cronexpr"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// experimentsCmd represents the Chaos experiments command
var experimentsCmd = &cobra.Command{
	Use:   "chaos-experiments",
	Short: "Display list of Chaos Experiments within the project",
	Long:  `Display list of Chaos Experiments within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var listExperimentRequest model.ListExperimentRequest
		var pid string
		pid, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if pid == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&pid)

			for pid == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		listAllExperiments, _ := cmd.Flags().GetBool("all")
		if !listAllExperiments {
			listExperimentRequest.Pagination = &model.Pagination{}
			listExperimentRequest.Pagination.Limit, _ = cmd.Flags().GetInt("count")
		}

		listExperimentRequest.Filter = &model.ExperimentFilterInput{}
		infraName, err := cmd.Flags().GetString("chaos-infra")
		utils.PrintError(err)
		listExperimentRequest.Filter.InfraName = &infraName

		experiments, err := experiment.GetExperimentList(pid, listExperimentRequest, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ The specified Project ID doesn't exist.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}

		outputFormat := ""
		outputPrompt := promptui.Select{
			Label: "Select output format",
			Items: []string{"Table", "json", "yaml"},
		}
		_, outputFormat, err = outputPrompt.Run()
		utils.PrintError(err)

		switch outputFormat {
		case "json":
			utils.PrintInJsonFormat(experiments.Data)

		case "yaml":
			utils.PrintInYamlFormat(experiments.Data)

		case "Table":
			itemsPerPage := 5
			page := 1
			totalExperiments := len(experiments.Data.ListExperimentDetails.Experiments)

			for {
				// calculating the start and end indices for the current page
				start := (page - 1) * itemsPerPage
				end := start + itemsPerPage
				if end > totalExperiments {
					end = totalExperiments
				}
				writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
				utils.White_B.Fprintln(writer, "CHAOS EXPERIMENT ID\tCHAOS EXPERIMENT NAME\tCHAOS EXPERIMENT TYPE\tNEXT SCHEDULE\tCHAOS INFRASTRUCTURE ID\tCHAOS INFRASTRUCTURE NAME\tLAST UPDATED By")

				for _, experiment := range experiments.Data.ListExperimentDetails.Experiments[start:end] {
					if experiment.CronSyntax != "" {
						utils.White.Fprintln(
							writer,
							experiment.ExperimentID+"\t"+experiment.Name+"\tCron Chaos Experiment\t"+cronexpr.MustParse(experiment.CronSyntax).Next(time.Now()).Format("January 2 2006, 03:04:05 pm")+"\t"+experiment.Infra.InfraID+"\t"+experiment.Infra.Name+"\t"+experiment.UpdatedBy.Username)
					} else {
						utils.White.Fprintln(
							writer,
							experiment.ExperimentID+"\t"+experiment.Name+"\tNon Cron Chaos Experiment\tNone\t"+experiment.Infra.InfraID+"\t"+experiment.Infra.Name+"\t"+experiment.UpdatedBy.Username)
					}
				}
				writer.Flush()
				paginationPrompt := promptui.Prompt{
					Label:     "Press Enter to show the next page (or type 'q' to quit)",
					AllowEdit: true,
					Default:   "",
				}

				userInput, err := paginationPrompt.Run()
				utils.PrintError(err)

				if userInput == "q" {
					break
				} else {
					page++
				}
			}
		}
	},
}

func init() {
	GetCmd.AddCommand(experimentsCmd)

	experimentsCmd.Flags().String("project-id", "", "Set the project-id to list Chaos experiments from the particular project. To see the projects, apply litmusctl get projects")
	experimentsCmd.Flags().Int("count", 30, "Set the count of Chaos experiments to display. Default value is 30")
	experimentsCmd.Flags().Bool("all", false, "Set to true to display all Chaos experiments")
	experimentsCmd.Flags().StringP("chaos-infra", "A", "", "Set the Chaos Infrastructure name to display all Chaos experiments targeted towards that particular Chaos Infrastructure.")

	experimentsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
