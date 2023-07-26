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
	"fmt"
	"os"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// experimentCmd represents the Chaos Experiment command
var experimentCmd = &cobra.Command{
	Use:   "chaos-experiment",
	Short: "Describe a Chaos Experiment within the project",
	Long:  `Describe a Chaos Experiment within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var describeExperimentRequest model.ListExperimentRequest

		pid, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if pid == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&pid)

			for pid == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		var experimentID string
		if len(args) == 0 {
			utils.White_B.Print("\nEnter the Chaos Experiment ID: ")
			fmt.Scanln(&experimentID)
		} else {
			experimentID = args[0]
		}

		// Handle blank input for Chaos Experiment ID
		if experimentID == "" {
			utils.Red.Println("⛔ Chaos Experiment ID can't be empty!!")
			os.Exit(1)
		}

		describeExperimentRequest.ExperimentIDs = append(describeExperimentRequest.ExperimentIDs, &experimentID)

		experiment, err := apis.GetExperimentList(pid, describeExperimentRequest, credentials)
		utils.PrintError(err)

		if len(experiment.Data.ListExperimentDetails.Experiments) == 0 {
			utils.Red.Println("⛔ No chaos experiment found with ID: ", experimentID)
			os.Exit(1)
		}

		yamlManifest, err := yaml.JSONToYAML([]byte(experiment.Data.ListExperimentDetails.Experiments[0].ExperimentManifest))
		if err != nil {
			utils.Red.Println("❌ Error parsing Chaos Experiment manifest: " + err.Error())
			os.Exit(1)
		}
		utils.PrintInYamlFormat(string(yamlManifest))
	},
}

func init() {
	DescribeCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Experiments from the particular project. To see the projects, apply litmusctl get projects")
}
