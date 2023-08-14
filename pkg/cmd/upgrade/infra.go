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
package upgrade

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var infraCmd = &cobra.Command{
	Use:   "chaos-infra",
	Short: `Upgrades the LitmusChaos Execution plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the project ID: ")
			fmt.Scanln(&projectID)
		}

		infraID, err := cmd.Flags().GetString("chaos-infra-id")
		utils.PrintError(err)

		if infraID == "" {
			utils.White_B.Print("\nEnter the Chaos Infra ID: ")
			fmt.Scanln(&infraID)
		}

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		utils.PrintError(err)

		output, err := apis.UpgradeInfra(context.Background(), credentials, projectID, infraID, kubeconfig)
		if err != nil {
			if strings.Contains(err.Error(), "no documents in result") {
				utils.Red.Println("❌ The specified Project ID or Chaos Infrastructure ID doesn't exist.")
				os.Exit(1)
			}
			utils.Red.Print("\n❌ Failed upgrading Chaos Infrastructure: \n" + err.Error() + "\n")
			os.Exit(1)
		}
		utils.White_B.Print("\n", output)
	},
}

func init() {
	UpgradeCmd.AddCommand(infraCmd)
	infraCmd.Flags().String("project-id", "", "Enter the project ID")
	infraCmd.Flags().String("kubeconfig", "", "Enter the kubeconfig path(default: $HOME/.kube/config))")
	infraCmd.Flags().String("chaos-infra-id", "", "Enter the Chaos Infrastructure ID")
}
