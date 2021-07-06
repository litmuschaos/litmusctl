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

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
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

		projects, err := apis.ListProject(credentials)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(projects.Data)
			break

		case "yaml":
			utils.PrintInYamlFormat(projects.Data)
			break

		case "":
			writer := tabwriter.NewWriter(os.Stdout, 8, 8, 8, '\t', tabwriter.AlignRight)
			fmt.Fprintln(writer, "PROJECT ID\tPROJECT NAME\tCREATEDAT")
			for _, project := range projects.Data.ListProjects {
				fmt.Fprintln(writer, project.ID+"\t"+project.Name+"\t"+project.CreatedAt+"\t")
			}
			writer.Flush()
			break
		}
	},
}

func init() {
	GetCmd.AddCommand(projectsCmd)

	projectsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
