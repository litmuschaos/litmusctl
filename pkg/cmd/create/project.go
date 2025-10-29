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

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/manifoldco/promptui"
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

		utils.Debug("Starting project creation process...")

		projectName, err := cmd.Flags().GetString("name")
		utils.PrintError(err)
		if projectName == "" {
			utils.Debug("Project name not provided via flag, prompting user...")
			// prompt to ask project name
			prompt := promptui.Prompt{
				Label:     "Enter a project name",
				AllowEdit: true,
			}

			result, err := prompt.Run()
			if err != nil {
				utils.Red.Printf("Error: %v\n", err)
				return
			}

			projectName = result
		}
		utils.Debugf("Creating project with name: %s", projectName)
		var response apis.CreateProjectResponse
		response, err = apis.CreateProjectRequest(projectName, credentials)
		if err != nil {
			utils.Red.Printf("❌ Error creating project: %v\n", err)
		} else {
			fmt.Printf("Response: %+v\n", response)
			projectID := response.Data.ID
			utils.Debugf("Project created with ID: %s", projectID)
			utils.White_B.Printf("Project '%s' created successfully with project ID - '%s'!🎉\n", projectName, projectID)
		}
	},
}

func init() {
	CreateCmd.AddCommand(projectCmd)
	projectCmd.Flags().String("name", "", "Set the project name to create it")
}
