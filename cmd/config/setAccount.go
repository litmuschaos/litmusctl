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
package config

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	config2 "github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/litmuschaos/litmusctl/tmp-pkg/constants"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"net/url"
	"os"
	"strings"
	"time"
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
	Run: func(command *cobra.Command, args []string) {
		var (
			authInput types.AuthInput
			defaultFileName = types.DefaultFileName
			err error
		)

		authInput.Endpoint, err = command.Flags().GetString("endpoint")
		utils.PrintError(err)

		authInput.Username, err = command.Flags().GetString("username")
		utils.PrintError(err)

		authInput.Password, err = command.Flags().GetString("password")
		utils.PrintError(err)

		if authInput.Endpoint == "" {
			fmt.Print("\n👉 Host endpoint where litmus is installed: ")
			fmt.Scanln(&authInput.Endpoint)
			for authInput.Endpoint == "" {
				fmt.Println("⛔ Host URL can't be empty!!")
				os.Exit(1)
			}

			ep := strings.TrimRight(authInput.Endpoint, "/")
			newUrl, err := url.Parse(ep)
			if err != nil {
				utils.PrintError(err)
			}

			authInput.Endpoint = newUrl.String()
		}

		if authInput.Username == "" {
			fmt.Print("🤔 Username [", constants.DefaultUsername, "]: ")
			fmt.Scanln(&authInput.Username)
			if authInput.Username == "" {
				authInput.Username = constants.DefaultUsername
			}
		}

		if authInput.Password == "" {
			fmt.Print("🙈 Password: ")
			pass, err := terminal.ReadPassword(0)
			utils.PrintError(err)

			if pass == nil {
				fmt.Println("\n⛔ Password cannot be empty!")
				os.Exit(1)
			}

			authInput.Password = string(pass)
		}

		if authInput.Endpoint != "" && authInput.Username != "" && authInput.Password != "" {
			exists := config2.FileExists(defaultFileName)
			lgt, err := config2.GetFileLength(defaultFileName)
			utils.PrintError(err)


			resp, err := apis.Auth(authInput)
			utils.PrintError(err)

			var (
				timeNow =  time.Now()
				newTime = timeNow.Add(time.Second * time.Duration(resp.ExpiresIn))
			)

			var user = types.User{
				ExpiresIn: fmt.Sprint(newTime.Unix()),
				Token:     resp.AccessToken,
				Username:  authInput.Username,
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

				err := config2.CreateNewLitmusCtlConfig(defaultFileName, config)
				utils.PrintError(err)

				os.Exit(0)
			} else {
				// checking syntax
				err = config2.ConfigSyntaxCheck(defaultFileName)
				utils.PrintError(err)


				var updateLitmusCtlConfig = types.UpdateLitmusCtlConfig{
					Account: account,
					CurrentAccount: authInput.Endpoint,
					CurrentUser: authInput.Username,
				}

				err = config2.UpdateLitmusCtlConfig(updateLitmusCtlConfig, defaultFileName)
				utils.PrintError(err)

			}

		} else {
			fmt.Println("Error: some flags are missing. Run 'litmusctl config set-account --help' for usage. ")
		}


	},
}

func init() {
	ConfigCmd.AddCommand(setAccountCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setAccountCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setAccountCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	//setAccountCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non-interactive usages (Default false)")
	setAccountCmd.Flags().StringP("endpoint", "e", "" , "Account endpoint. Mandatory")
	setAccountCmd.Flags().StringP("username", "u", "", "Account username. Mandatory")
	setAccountCmd.Flags().StringP("password", "p", "", "Account password. Mandatory")

	//litmusctl config set-account --endpoint --username --password
}
