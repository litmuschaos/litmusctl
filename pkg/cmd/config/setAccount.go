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

	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"

	"github.com/golang-jwt/jwt"
	"github.com/litmuschaos/litmusctl/pkg/apis"

	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
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

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		utils.PrintError(err)

		if nonInteractive {
			authInput.Endpoint, err = cmd.Flags().GetString("endpoint")
			utils.PrintError(err)

			authInput.Username, err = cmd.Flags().GetString("username")
			utils.PrintError(err)

			authInput.Password, err = cmd.Flags().GetString("password")
			utils.PrintError(err)

		} else {
			// prompts for account details
			promptEndpoint := promptui.Prompt{
				Label: "Host endpoint where litmus is installed",
			}
			authInput.Endpoint, err = promptEndpoint.Run()
			utils.PrintError(err)

			promptUsername := promptui.Prompt{
				Label:   "Username [Default: " + utils.DefaultUsername + "]",
				Default: utils.DefaultUsername,
			}
			authInput.Username, err = promptUsername.Run()
			utils.PrintError(err)

			promptPassword := promptui.Prompt{
				Label: "Password",
				Mask:  '*',
			}
			pass, err := promptPassword.Run()
			utils.PrintError(err)
			authInput.Password = pass
		}

		// Validate and format the endpoint URL
		ep := strings.TrimRight(authInput.Endpoint, "/")
		newURL, err := url.Parse(ep)
		utils.PrintError(err)
		authInput.Endpoint = newURL.String()

		if authInput.Endpoint == "" {
			utils.White_B.Print("\nHost endpoint where litmus is installed: ")
			fmt.Scanln(&authInput.Endpoint)

			if authInput.Endpoint == "" {
				utils.Red.Println("\n⛔ Host URL can't be empty!!")
				os.Exit(1)
			}

			ep := strings.TrimRight(authInput.Endpoint, "/")
			newUrl, err := url.Parse(ep)
			utils.PrintError(err)

			authInput.Endpoint = newUrl.String()
		}

		if authInput.Endpoint != "" && authInput.Username != "" && authInput.Password != "" {
			exists := config.FileExists(configFilePath)
			var lgt int
			if exists {
				lgt, err = config.GetFileLength(configFilePath)
				utils.PrintError(err)
			}

			resp, err := apis.Auth(authInput, apis.Client)
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
				Endpoint:       authInput.Endpoint,
				Users:          users,
				ServerEndpoint: authInput.Endpoint,
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

			} else {
				// checking syntax
				err = config.ConfigSyntaxCheck(configFilePath)
				utils.PrintError(err)

				var updateLitmusCtlConfig = types.UpdateLitmusCtlConfig{
					Account:        account,
					CurrentAccount: authInput.Endpoint,
					CurrentUser:    claims["username"].(string),
					ServerEndpoint: authInput.Endpoint,
				}

				err = config.UpdateLitmusCtlConfig(updateLitmusCtlConfig, configFilePath)
				utils.PrintError(err)
			}
			utils.White_B.Printf("\naccount.username/%s configured", claims["username"].(string))
		} else {
			utils.Red.Println("\nError: some flags are missing. Run 'litmusctl config set-account --help' for usage. ")
		}

		serverResp, err := experiment.GetServerVersion(authInput.Endpoint, apis.Client)
		var isCompatible bool
		if err != nil {
			utils.Red.Println("\nError: ", err)
		} else {
			compatibilityArr := utils.CompatibilityMatrix[os.Getenv("CLIVersion")]
			for _, v := range compatibilityArr {
				if v == serverResp.Data.GetServerVersion.Value {
					isCompatible = true
					break
				}
			}

			if !isCompatible {
				utils.Red.Println("\n🚫 ChaosCenter version: " + serverResp.Data.GetServerVersion.Value + " is not compatible with the installed LitmusCTL version: " + os.Getenv("CLIVersion"))
				utils.White_B.Println("Compatible ChaosCenter versions are: ")
				utils.White_B.Print("[ ")
				for _, v := range compatibilityArr {
					utils.White_B.Print("'" + v + "' ")
				}
				utils.White_B.Print("]\n")
			} else {
				utils.White_B.Println("\n✅  Installed versions of ChaosCenter and LitmusCTL are compatible! ")
			}
		}
	},
}

func init() {
	ConfigCmd.AddCommand(setAccountCmd)
	setAccountCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	setAccountCmd.Flags().StringP("endpoint", "e", "", "Account endpoint. Mandatory")
	setAccountCmd.Flags().StringP("username", "u", "", "Account username. Mandatory")
	setAccountCmd.Flags().StringP("password", "p", "", "Account password. Mandatory")
}
