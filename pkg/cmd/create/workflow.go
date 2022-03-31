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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"sigs.k8s.io/yaml"

	types "github.com/litmuschaos/litmusctl/pkg/types"

	"github.com/spf13/cobra"
)

// workflowCmd represents the project command
var workflowCmd = &cobra.Command{
	Use: "workflow",
	Short: `Create a workflow
	Example:
	#create a workflow
	litmusctl create workflow -f workflow.yaml

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var chaosWorkFlowInput types.ChaosWorkFlowInput
		var chaosExperiment chaosTypes.ChaosExperiment
		var weightages []types.WeightagesInput

		workflowManifest, err := cmd.Flags().GetString("file")
		utils.PrintError(err)

		chaosWorkFlowInput.ProjectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if chaosWorkFlowInput.ProjectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&chaosWorkFlowInput.ProjectID)

			if chaosWorkFlowInput.ProjectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		workflow := readManifestFile(workflowManifest)
		workflowStr, _ := json.Marshal(workflow)
		chaosWorkFlowInput.WorkflowManifest = string(workflowStr)
		chaosWorkFlowInput.WorkflowName = workflow.ObjectMeta.Name

		for _, t := range workflow.Spec.Templates {
			if len(t.Inputs.Artifacts) != 0 {
				err := yaml.Unmarshal([]byte(t.Inputs.Artifacts[0].Raw.Data), &chaosExperiment)
				if chaosExperiment.Kind == "ChaosEngine" || err != nil {
					continue // Tried to parse a ChaosEngine spec
				}
				weightages = append(weightages,
					types.WeightagesInput{
						ExperimentName: chaosExperiment.ObjectMeta.Name,
						Weightage:      10, // TODO: fetch from annotation
					},
				)
			}
		}

		chaosWorkFlowInput.Weightages = weightages
		chaosWorkFlowInput.IsCustomWorkflow = true

		chaosWorkFlowInput.ClusterID, err = cmd.Flags().GetString("cluster-id")
		utils.PrintError(err)

		if chaosWorkFlowInput.ClusterID == "" {
			utils.White_B.Print("\nEnter the Cluster ID: ")
			fmt.Scanln(&chaosWorkFlowInput.ClusterID)

			if chaosWorkFlowInput.ClusterID == "" {
				utils.Red.Println("⛔ Cluster ID can't be empty!!")
				os.Exit(1)
			}
		}

		apis.CreateWorkflow(chaosWorkFlowInput, credentials)
	},
}

// TODO: Move this to utils
func readManifestFile(file string) v1alpha1.Workflow {
	var body []byte
	var workflowManifest v1alpha1.Workflow

	body, _ = ioutil.ReadFile(file)
	_ = yaml.Unmarshal(body, &workflowManifest)
	return workflowManifest
}

func init() {
	CreateCmd.AddCommand(workflowCmd)

	workflowCmd.Flags().String("project-id", "", "Set the project-id to create workflow for the particular project. To see the projects, apply litmusctl get projects")
	workflowCmd.Flags().String("cluster-id", "", "Set the cluster-id to create workflow for the particular cluster. To see the projects, apply litmusctl get agents")

	workflowCmd.Flags().StringP("file", "f", "", "The manifest file for the workflow")
}
