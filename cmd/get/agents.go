/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"os"
	"text/tabwriter"
)

// agentsCmd represents the agents command
var agentsCmd = &cobra.Command{
	Use:   "agents",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		obj, err := config.YamltoObject(types.DefaultFileName)
		utils.PrintError(err)

		var token string
		for _, account := range obj.Accounts {
			if account.Endpoint == obj.CurrentAccount {
				for _, user := range account.Users {
					if user.Username == obj.CurrentUser {
						token = user.Token
					}
				}
			}
		}

		var credentials = types.Credentials{
			Username: obj.CurrentUser,
			Token:    token,
			Endpoint: obj.CurrentAccount,
		}

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

		case "json":
			var out bytes.Buffer
			byt, err := json.Marshal(agents.Data)
			utils.PrintError(err)

			err = json.Indent(&out, byt, "", "  ")
			utils.PrintError(err)

			fmt.Println(out.String())
			break

		case "yaml":
			byt, err := yaml.Marshal(agents.Data)
			utils.PrintError(err)
			fmt.Println(string(byt))

			break
		}

	},
}

func init() {
	GetCmd.AddCommand(agentsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// agentsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	agentsCmd.Flags().String("project-id", "", "Help message for toggle")
	agentsCmd.Flags().StringP("output", "o", "", "Help message for toggle")

}
