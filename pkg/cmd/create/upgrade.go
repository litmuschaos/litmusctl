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
	"context"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: `Upgrades the LitmusChaos agent plane.`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)

		utils.PrintError(err)

		c := context.Background()

		// TAKE PROJECT_ID AND NAMESPACE AS INPUT AND PASS IT THROUGH GetManifests FUNCTION
		apis.GetManifest(c, credentials)
	},
}

func init() {
	CreateCmd.AddCommand(UpgradeCmd)
}
