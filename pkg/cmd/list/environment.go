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

package list

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis/environment"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

var ListChaosEnvironmentCmd = &cobra.Command{
	Use:   "chaos-environments",
	Short: "Get Chaos Environments within the project",
	Long:  `Display the Chaos Environments within the project with the targeted id `,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			for projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		environmentList, err := environment.GetEnvironmentList(projectID, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ You don't have enough permissions to access this resource.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}
		environmentListData := environmentList.Data.ListEnvironmentDetails.Environments

		itemsPerPage := 5
		page := 1
		totalEnvironments := len(environmentListData)

		writer := tabwriter.NewWriter(os.Stdout, 30, 8, 0, '\t', tabwriter.AlignRight)
		utils.White_B.Fprintln(writer, "CHAOS ENVIRONMENT ID\tCHAOS ENVIRONMENT NAME\tCREATED AT\tCREATED BY")
		for {
			writer.Flush()
			// calculating the start and end indices for the current page
			start := (page - 1) * itemsPerPage
			if start >= totalEnvironments {
				writer.Flush()
				utils.Red.Println("No more environments to display")
				break
			}
			end := start + itemsPerPage
			if end > totalEnvironments {
				end = totalEnvironments

			}
			for _, environment := range environmentListData[start:end] {
				intTime, err := strconv.ParseInt(environment.CreatedAt, 10, 64)
				if err != nil {
					fmt.Println("Error converting CreatedAt to int64:", err)
					continue
				}
				humanTime := time.Unix(intTime, 0)
				utils.White.Fprintln(
					writer,
					environment.EnvironmentID+"\t"+environment.Name+"\t"+humanTime.String()+"\t"+environment.CreatedBy.Username,
				)
			}
			writer.Flush()
			// Check if it's the last item or if user wants to see more
			paginationPrompt := promptui.Prompt{
				Label:     "Press Enter to show more environments (or type 'q' to quit)",
				AllowEdit: true,
				Default:   "",
			}

			userInput, err := paginationPrompt.Run()
			utils.PrintError(err)

			if userInput == "q" {
				break
			}
			// Move to the next page
			page++
		}
	},
}

func init() {
	ListCmd.AddCommand(ListChaosEnvironmentCmd)
	ListChaosEnvironmentCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Environments from a particular project.")
}
