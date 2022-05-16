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

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: `Upgrades the LitmusChaos agent plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the project ID: ")
			fmt.Scanln(&projectID)
		}

		cluster_id, err := cmd.Flags().GetString("cluster-id")
		utils.PrintError(err)

		if cluster_id == "" {
			utils.White_B.Print("\nEnter the cluster ID: ")
			fmt.Scanln(&cluster_id)
		}

		output, err := apis.UpgradeAgent(context.Background(), credentials, projectID, cluster_id)
		if err != nil {
			utils.Red.Print(output)
			utils.PrintError(err)
		} else {
			utils.White.Print(output)
		}
	},
}

func init() {
	UpgradeCmd.AddCommand(agentCmd)
	agentCmd.Flags().String("project-id", "", "Enter the project ID")
	agentCmd.Flags().String("cluster-id", "", "Enter the cluster ID")
}
