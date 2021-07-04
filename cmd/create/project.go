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
package create

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use:   "project",
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

		projectName, err := cmd.Flags().GetString("name")
		utils.PrintError(err)

		if projectName == "" {
			fmt.Print("\n📁 Enter the project name: ")
			fmt.Scanln(&projectName)
			for projectName == "" {
				fmt.Println("⛔ Project name can't be empty!!")
				os.Exit(1)
			}
		}

		err = apis.CreateProjectRequest(projectName, credentials)
		utils.PrintError(err)

	},
}

func init() {
	CreateCmd.AddCommand(projectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// projectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	projectCmd.Flags().String("name", "", "Help message for toggle")
}
