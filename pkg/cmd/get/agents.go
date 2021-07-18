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
	"github.com/fatih/color"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
	"text/tabwriter"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "Display list of agents within the project",
	Long:  `Display list of agents within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		var (
			cyan_b = color.New(color.FgCyan, color.Bold)
			cyan = color.New(color.FgCyan)
			red = color.New(color.FgRed)
		)

		if projectID == "" {
			cyan_b.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			for projectID == "" {
				red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		agents, err := apis.GetAgentList(credentials, projectID)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "json":
			utils.PrintInJsonFormat(agents.Data)

		case "yaml":
			utils.PrintInYamlFormat(agents.Data)

		case "":

			writer := tabwriter.NewWriter(os.Stdout, 10, 10, 10, '\t', tabwriter.AlignRight)
			cyan_b.Fprintln(writer, "AGENT ID\tAGENT NAME\tSTATUS")

			for _, agent := range agents.Data.GetAgent {
				var status string
				if agent.IsActive {
					status = "ACTIVE"
				} else {
					status = "INACTIVE"
				}
				cyan.Fprintln(writer, agent.ClusterID + "\t" + agent.AgentName + "\t" + status + "\t")
			}
			writer.Flush()
		}
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)

	agentsCmd.Flags().String("project-id", "", "Set the project-id. To retrieve projects. Apply `litmusctl get projects`")

	agentsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
