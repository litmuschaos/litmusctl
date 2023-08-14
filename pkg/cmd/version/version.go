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
package version

import (
	"os"

	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of litmusctl",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		compatibilityArr := utils.CompatibilityMatrix[os.Getenv("CLIVersion")]
		utils.White_B.Println("Litmusctl version: ", os.Getenv("CLIVersion"))
		utils.White_B.Println("Compatible ChaosCenter versions: ")
		utils.White_B.Print("[ ")
		for _, v := range compatibilityArr {
			utils.White_B.Print("'" + v + "' ")
		}
		utils.White_B.Print("]\n")
	},
}
