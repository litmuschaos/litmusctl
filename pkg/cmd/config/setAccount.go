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
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

// setAccountCmd represents the setAccount command
var setAccountCmd = &cobra.Command{
	Use: "set-account",
	Short: `Sets an account entry in litmusconfig.
		Examples(s)
		#set a new account
		litmusctl config set-account  --endpoint "" --password "" --username ""
		`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := utils.GetLitmusConfigPath(cmd)

		var (
			authInput types.AuthInput
			err       error
		)

		authInput.Endpoint, err = cmd.Flags().GetString("endpoint")
		utils.PrintError(err)

		authInput.Username, err = cmd.Flags().GetString("username")
		utils.PrintError(err)

		authInput.Password, err = cmd.Flags().GetString("password")
		utils.PrintError(err)

		if authInput.Endpoint == "" {
			utils.White_B.Print("\nHost endpoint where litmus is installed: ")
			fmt.Scanln(&authInput.Endpoint)

			for authInput.Endpoint == "" {
				utils.Red.Println("\n⛔ Host URL can't be empty!!")
				os.Exit(1)
			}

			ep := strings.TrimRight(authInput.Endpoint, "/")
			newUrl, err := url.Parse(ep)
			utils.PrintError(err)

			authInput.Endpoint = newUrl.String()
		}

		if authInput.Username == "" {
			utils.White_B.Print("\nUsername [Default: ", utils.DefaultUsername, "]: ")
			fmt.Scanln(&authInput.Username)
			if authInput.Username == "" {
				authInput.Username = utils.DefaultUsername
			}
		}

		if authInput.Password == "" {
			utils.White_B.Print("\nPassword: ")
			pass, err := term.ReadPassword(0)
			utils.PrintError(err)

			if pass == nil {
				utils.Red.Println("\n⛔ Password cannot be empty!")
				os.Exit(1)
			}

			authInput.Password = string(pass)
		}

		if authInput.Endpoint != "" && authInput.Username != "" && authInput.Password != "" {
			exists := config.FileExists(configFilePath)
			var lgt int
			if exists {
				lgt, err = config.GetFileLength(configFilePath)
				utils.PrintError(err)
			}

			resp, err := apis.Auth(authInput)
			utils.PrintError(err)
			// Decoding token
			token, _ := jwt.Parse(resp.AccessToken, nil)
			if token == nil {
				os.Exit(1)
			}
			claims, _ := token.Claims.(jwt.MapClaims)

			var user = types.User{
				ExpiresIn: fmt.Sprint(time.Now().Add(time.Second * time.Duration(resp.ExpiresIn)).Unix()),
				Token:     resp.AccessToken,
				Username:  claims["username"].(string),
			}

			var users []types.User
			users = append(users, user)

			var account = types.Account{
				Endpoint: authInput.Endpoint,
				Users:    users,
			}

			// If config file doesn't exist or length of the file is zero.
			if !exists || lgt == 0 {

				var accounts []types.Account
				accounts = append(accounts, account)

				var litmuCtlConfig = types.LitmuCtlConfig{
					APIVersion:     "v1",
					Kind:           "Config",
					CurrentAccount: authInput.Endpoint,
					CurrentUser:    claims["username"].(string),
					Accounts:       accounts,
				}

				err := config.CreateNewLitmusCtlConfig(configFilePath, litmuCtlConfig)
				utils.PrintError(err)

				os.Exit(0)
			} else {
				// checking syntax
				err = config.ConfigSyntaxCheck(configFilePath)
				utils.PrintError(err)

				var updateLitmusCtlConfig = types.UpdateLitmusCtlConfig{
					Account:        account,
					CurrentAccount: authInput.Endpoint,
					CurrentUser:    claims["username"].(string),
				}

				err = config.UpdateLitmusCtlConfig(updateLitmusCtlConfig, configFilePath)
				utils.PrintError(err)
			}
			utils.White_B.Printf("\naccount.username/%s configured", claims["username"].(string))

			serverResp, err := apis.GetServerVersion()
			var isCompatible bool
			if err != nil {
				utils.Red.Println("\nError: ", err)
			} else {
				compatibilityArr := utils.CompatibilityMatrix[os.Getenv("CLIVersion")]
				for _, v := range compatibilityArr {
					if v == serverResp.Data.GetServerVersion.Value {
						isCompatible = true
					}
				}

				if isCompatible == false {
					utils.Red.Println("\n🚫 ChaosCenter version: " + serverResp.Data.GetServerVersion.Value + " is not compatible with the installed LitmusCTL version: " + os.Getenv("CLIVersion"))
					utils.White_B.Println("Compatible ChaosCenter versions are: ")
					utils.White_B.Print("[ ")
					for _, v := range compatibilityArr {
						utils.White_B.Print("'" + v + "' ")
					}
					utils.White_B.Print("]")
					utils.White_B.Println("")
				} else {
					utils.White_B.Println("\n✅  Installed versions of ChaosCenter and LitmusCTL are compatible! ")
				}
			}

		} else {
			utils.Red.Println("\nError: some flags are missing. Run 'litmusctl config set-account --help' for usage. ")
		}
	},
}

func init() {
	ConfigCmd.AddCommand(setAccountCmd)

	setAccountCmd.Flags().StringP("endpoint", "e", "", "Account endpoint. Mandatory")
	setAccountCmd.Flags().StringP("username", "u", "", "Account username. Mandatory")
	setAccountCmd.Flags().StringP("password", "p", "", "Account password. Mandatory")
}
