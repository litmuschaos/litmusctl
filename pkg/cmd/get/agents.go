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

		if projectID == "" {
			fmt.Print("\n📁 Enter the Project ID: ")
			fmt.Scanln(&projectID)

			for projectID == "" {
				fmt.Println("⛔ Project ID can't be empty!!")
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
			break

		case "yaml":
			utils.PrintInYamlFormat(agents.Data)
			break

		case "":
			writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
			fmt.Fprintln(writer, "AGENT ID\tAGENT NAME\tSTATUS")
			for _, agent := range agents.Data.GetAgent {
				var status string
				if agent.IsActive {
					status = "ACTIVE"
				} else {
					status = "INACTIVE"
				}
				fmt.Fprintln(writer, agent.ClusterID+"\t"+agent.AgentName+"\t"+status)
			}
			writer.Flush()
			break
		}
	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)

	agentsCmd.Flags().String("project-id", "", "Help message for toggle")

	agentsCmd.Flags().StringP("output", "o", "", "Output format. One of:\njson|yaml")
}
