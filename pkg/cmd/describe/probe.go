// /*
// Copyright © 2021 The LitmusChaos Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// */

package describe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	apis "github.com/litmuschaos/litmusctl/pkg/apis/probe"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

var probeCmd = &cobra.Command{
	Use:   "probe",
	Short: "Describe a Probe within the project",
	Long:  `Describe a Probe within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var getProbeYAMLRequest model.GetProbeYAMLRequest

		pid, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if pid == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Project ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.PrintError(err)
				os.Exit(1)
			}
			pid = result
		}

		var probeID string
		probeID, err = cmd.Flags().GetString("probe-id")
		utils.PrintError(err)
		// Handle blank input for Probe ID

		if probeID == "" {
			utils.White_B.Print("\nEnter the Probe ID: ")
			fmt.Scanln(&probeID)

			if probeID == "" {
				utils.Red.Println("⛔ Probe ID can't be empty!!")
				os.Exit(1)
			}
		}
		getProbeYAMLRequest.ProbeName = probeID

		probeMode, err := cmd.Flags().GetString("mode")
		utils.PrintError(err)

		if probeMode == "" {
			prompt := promptui.Select{
				Label: "Please select the probe mode ?",
				Items: []string{"SOT", "EOT", "Edge", "Continuous", "OnChaos"},
			}
			_, option, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}
			probeMode = option
			fmt.Printf("You chose %q\n", option)
		}
		getProbeYAMLRequest.Mode = model.Mode(probeMode)
		getProbeYAML, err := apis.GetProbeYAMLRequest(pid, getProbeYAMLRequest, credentials)
		if err != nil {
			utils.Red.Println(err)
			os.Exit(1)
		}
		getProbeYAMLData := getProbeYAML.Data.GetProbeYAML

		probeOutput, err := cmd.Flags().GetString("output")
		utils.PrintError(err)

		switch probeOutput {
		case "json":
			jsonData, _ := yaml.YAMLToJSON([]byte(getProbeYAMLData))
			// utils.PrintInJsonFormat(string(jsonData))
			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonData, "", "    ") // Adjust the indentation as needed
			if err != nil {
				utils.Red.Println("❌ Error formatting JSON: " + err.Error())
				os.Exit(1)
			}

			fmt.Println(prettyJSON.String())

		default:
			utils.PrintInYamlFormat(getProbeYAMLData)
		}

	},
}

func init() {
	DescribeCmd.AddCommand(probeCmd)
	probeCmd.Flags().String("project-id", "", "Set the project-id to get Probe details from the particular project. To see the projects, apply litmusctl get projects")
	probeCmd.Flags().String("probe-id", "", "Set the probe-id to get the Probe details in Yaml format")
	probeCmd.Flags().String("mode", "", "Set the mode for the probes from SOT/EOT/Edge/Continuous/OnChaos ")
	probeCmd.Flags().String("output", "", "Set the output format for the probe")
}
