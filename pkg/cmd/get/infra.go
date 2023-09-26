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

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis/infrastructure"

	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// InfraCmd represents the Infra command
var InfraCmd = &cobra.Command{
	Use:   "chaos-infra",
	Short: "Display list of Chaos Infrastructures within the project",
	Long:  `Display list of Chaos Infrastructures within the project`,
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

		infras, err := infrastructure.GetInfraList(credentials, projectID, models.ListInfraRequest{})
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

		switch output {
		case "json":
			utils.PrintInJsonFormat(infras.Data)

		case "yaml":
			utils.PrintInYamlFormat(infras.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 4, 8, 1, '\t', 0)
			utils.White_B.Fprintln(writer, "CHAOS INFRASTRUCTURE ID \tCHAOS INFRASTRUCTURE NAME\tSTATUS\tCHAOS ENVIRONMENT ID\t")

			for _, infra := range infras.Data.ListInfraDetails.Infras {
				var status string
				if infra.IsActive {
					status = "ACTIVE"
				} else {
					status = "INACTIVE"
				}
				utils.White.Fprintln(writer, infra.InfraID+"\t"+infra.Name+"\t"+status+"\t"+infra.EnvironmentID+"\t")
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(InfraCmd)

	InfraCmd.Flags().String("project-id", "", "Set the project-id. To retrieve projects. Apply `litmusctl get projects`")

	InfraCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
