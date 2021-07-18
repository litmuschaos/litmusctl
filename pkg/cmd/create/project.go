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
package create

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// projectCmd represents the project command
var projectCmd = &cobra.Command{
	Use: "project",
	Short: `Create a project
	Example:
	#create a project
	litmusctl create project --name new-proj

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectName, err := cmd.Flags().GetString("name")
		utils.PrintError(err)

		var cyan = color.New(color.FgCyan, color.Bold)

		if projectName == "" {
			cyan.Print("\nEnter a project name: ")
			fmt.Scanln(&projectName)
		}

		_, err = apis.CreateProjectRequest(projectName, credentials)
		utils.PrintError(err)
	},
}

func init() {
	CreateCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("name", "", "Set the project name to create it")
}
