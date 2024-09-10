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
	"os"
	"text/tabwriter"
	"time"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "Display list of projects",
	Long:  `Display list of projects`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		outputFormat, _ := cmd.Flags().GetString("output")

		projects, err := apis.ListProject(credentials)
		utils.PrintError(err)

		switch outputFormat {
		case "json":
			utils.PrintInJsonFormat(projects.Data)

		case "yaml":
			utils.PrintInYamlFormat(projects.Data.Projects)

		case "":
			itemsPerPage := 5
			page := 1
			totalProjects := len(projects.Data.Projects)

			for {
				// calculating the start and end indices for the current page
				start := (page - 1) * itemsPerPage
				end := start + itemsPerPage
				if end > totalProjects {
					end = totalProjects

				}
				// check if there are no more projects to display
				if start >= totalProjects {
					utils.Red.Println("No more projects to display")
					break
				}

				// displaying the projects for the current page
				writer := tabwriter.NewWriter(os.Stdout, 8, 8, 8, '\t', tabwriter.AlignRight)
				utils.White_B.Fprintln(writer, "PROJECT ID\tPROJECT NAME\tCREATED AT")
				for _, project := range projects.Data.Projects[start:end] {
					intTime := project.CreatedAt
					humanTime := time.Unix(intTime/1000, 0) // Convert milliseconds to second
					utils.White.Fprintln(writer, project.ID+"\t"+project.Name+"\t"+humanTime.String()+"\t")
				}
				writer.Flush()

				// pagination prompt
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
	GetCmd.AddCommand(projectsCmd)

	projectsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
