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

// projectCmd represents the project command
var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var configFilePath string
		configFilePath, err := cmd.Flags().GetString("config")
		utils.PrintError(err)

		if configFilePath == "" {
			configFilePath = types.DefaultFileName
		}

		obj, err := config.YamltoObject(configFilePath)
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

		projects, err := apis.ListProject(credentials)
		utils.PrintError(err)

		output, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch output {
		case "":
			writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
			fmt.Fprintln(writer, "PROJECT ID\tPROJECT NAME\tCREATEDAT")
			for _, project := range projects.Data.ListProjects {
				fmt.Fprintln(writer, project.ID+"\t"+project.Name+"\t"+project.CreatedAt+"\t")
			}
			writer.Flush()
			break

		case "json":
			var out bytes.Buffer
			byt, err := json.Marshal(projects.Data)
			utils.PrintError(err)

			err = json.Indent(&out, byt, "", "  ")
			utils.PrintError(err)

			fmt.Println(out.String())
			break

		case "yaml":
			byt, err := yaml.Marshal(projects.Data)
			utils.PrintError(err)

			fmt.Println(string(byt))
			break
		}
	},
}

func init() {
	GetCmd.AddCommand(projectsCmd)

	projectsCmd.Flags().StringP("output", "o", "", "Help message for toggle")
}
