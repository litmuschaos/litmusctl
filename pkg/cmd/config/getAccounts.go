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
	"text/tabwriter"

	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// getAccountsCmd represents the getAccounts command
var getAccountsCmd = &cobra.Command{
	Use:   "get-accounts",
	Short: "Display accounts defined in the litmusconfig",
	Long:  `Display accounts defined in the litmusconfig`,
	Run: func(cmd *cobra.Command, args []string) {
		configFilePath := utils.GetLitmusConfigPath(cmd)

		obj, err := config.YamltoObject(configFilePath)
		utils.PrintError(err)

		writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
		fmt.Fprintln(writer, "CURRENT\tENDPOINT\tUSERNAME\tEXPIRESIN")
		for _, account := range obj.Accounts {
			for _, user := range account.Users {
				if obj.CurrentUser == user.Username && obj.CurrentAccount == account.Endpoint {
					fmt.Fprintln(writer, "*"+"\t"+account.Endpoint+"\t"+user.Username+"\t"+user.ExpiresIn)
				} else {
					fmt.Fprintln(writer, ""+"\t"+account.Endpoint+"\t"+user.Username+"\t"+user.ExpiresIn)
				}
			}
		}
		writer.Flush()
	},
}

func init() {
	ConfigCmd.AddCommand(getAccountsCmd)
}
