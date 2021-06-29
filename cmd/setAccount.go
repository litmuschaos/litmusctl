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
package cmd

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"os"
)


// setAccountCmd represents the setAccount command
var setAccountCmd = &cobra.Command{
	Use:   "set-account",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var (
			authInput types.AuthInput
			defaultFileName = types.DefaultFileName
		)

		authInput.Endpoint, _ = cmd.Flags().GetString("endpoint")
		authInput.Username, _ = cmd.Flags().GetString("username")
		authInput.Password, _ = cmd.Flags().GetString("password")

		if authInput.Endpoint != "" && authInput.Username != "" && authInput.Password != "" {
			exists := utils.FileExists(defaultFileName)
			lgt, err := utils.GetFileLength(defaultFileName)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			resp, err := apis.Auth(authInput)
			if err != nil {
				fmt.Println(err)
			}

			var user = types.User{
				ExpiresIn: string(resp.ExpiresIn),
				Password: authInput.Password,
				Token: resp.AccessToken,
				Username: authInput.Username,
			}

			var users []types.User
			users = append(users, user)

			var account = types.Account{
				Endpoint: authInput.Endpoint,
				Users: users,
			}


			// If config file doesn't exist or length of the file is zero.
			if !exists || lgt == 0 {

				var accounts []types.Account
				accounts = append(accounts, account)

				var config = types.LitmuCtlConfig{
					APIVersion: "v1",
					Kind: "Config",
					CurrentAccount: authInput.Endpoint,
					CurrentUser: authInput.Username,
					Accounts: accounts,
				}

				err := utils.CreateNewLitmusCtlConfig(defaultFileName, config)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				os.Exit(0)
			} else {
				// checking syntax
				err = utils.ConfigSyntaxCheck(defaultFileName)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				err = utils.UpdateLitmusCtlConfig(account, defaultFileName)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

			}

		} else {
			fmt.Println("Error: some flags are missing. Run 'litmusctl config set-account --help' for usage. ")
		}


	},
}

func init() {
	configCmd.AddCommand(setAccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	setAccountCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non-interactive usages (Default false)")
	setAccountCmd.Flags().StringP("endpoint", "e", "" , "Account endpoint. Mandatory")
	setAccountCmd.Flags().StringP("username", "u", "", "Account username. Mandatory")
	setAccountCmd.Flags().StringP("password", "p", "", "Account password. Mandatory")

	//litmusctl config set-account --endpoint --username --password
}