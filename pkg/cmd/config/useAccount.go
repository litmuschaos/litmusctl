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
package config

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// useAccountCmd represents the useAccount command
var useAccountCmd = &cobra.Command{
	Use:   "use-account",
	Short: "Sets the current-account and current-username in a litmusconfig file",
	Long:  `Sets the current-account and current-username in a litmusconfig file`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := utils.GetLitmusConfigPath(cmd)

		endpoint, err := cmd.Flags().GetString("endpoint")
		utils.PrintError(err)

		if endpoint == "" {
			utils.White_B.Print("\nHost endpoint where litmus is installed: ")
			fmt.Scanln(&endpoint)

			for endpoint == "" {
				utils.Red.Println("\n⛔ Host URL can't be empty!!")
				os.Exit(1)
			}
		}

		username, err := cmd.Flags().GetString("username")
		utils.PrintError(err)

		if username == "" {
			utils.White_B.Print("\nUsername: ")
			fmt.Scanln(&username)

			for username == "" {
				utils.Red.Println("\n⛔ Username cannot be empty!!")
				os.Exit(1)
			}
		}

		if username == "" || endpoint == "" {
			utils.Red.Println("endpoint or username is not set")
			os.Exit(1)
		}

		exists := config.FileExists(configFilePath)

		err = config.ConfigSyntaxCheck(configFilePath)
		utils.PrintError(err)

		if exists {
			litmusconfig, err := config.YamltoObject(configFilePath)
			utils.PrintError(err)

			isAccountExist := config.IsAccountExists(litmusconfig, username, endpoint)
			if isAccountExist {
				err = config.UpdateCurrent(types.Current{
					CurrentAccount: endpoint,
					CurrentUser:    username,
				}, configFilePath)

				utils.PrintError(err)
				fmt.Printf("\n✅ Successfully set the current account to '%s' at '%s'\n", username, endpoint)
			} else {
				utils.Red.Println("\n⛔ Account not exists")
				os.Exit(1)
			}
		} else {
			utils.Red.Println("\n⛔ File not exists")
			os.Exit(1)
		}
	},
}

func init() {
	ConfigCmd.AddCommand(useAccountCmd)
	useAccountCmd.Flags().StringP("username", "u", "", "Help message for toggle")
	useAccountCmd.Flags().StringP("endpoint", "e", "", "Help message for toggle")
}
