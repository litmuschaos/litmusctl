/*
Copyright © 2021 The LitmusChaos Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.W
See the License for the specific language governing permissions and
limitations under the License.
*/
package delete

import (
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/probe"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var probeCmd = &cobra.Command{
	Use: "probe",
	Short: `Delete a Probe
	Example:
	#delete a Probe
	litmusctl delete probe --project-id=c520650e-7cb6-474c-b0f0-4df07b2b025b --probe-id="example"
	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,

	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		projectID, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID

		if projectID == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Project ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.Red.Println("⛔ Error:", err)
				os.Exit(1)
			}
			projectID = result
		}
		probeID, err := cmd.Flags().GetString("probe-id")
		// Handle blank input for Probe ID
		if probeID == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Probe ID",
			}
			IDinput, err := prompt.Run()
			if err != nil {
				utils.Red.Println("⛔ Error:", err)
				os.Exit(1)
			}
			probeID = IDinput
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == projectID {
				project = p
			}
		}
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == "Owner" || member.Role == "Editor") {
				editAccess = true
			}
		}
		if !editAccess {
			utils.Red.Println("⛔ User doesn't have edit access to the project!!")
			os.Exit(1)
		}

		// confirm before deletion

		prompt := promptui.Prompt{
			Label:     "Are you sure you want to delete this probe and all the associations with experiment runs from the chaos control plane (y/n)",
			AllowEdit: true,
		}
		result, err := prompt.Run()
		if err != nil {
			utils.Red.Println("⛔ Error:", err)
			os.Exit(1)
		}

		if result != "y" {
			utils.White_B.Println("\n❌ Probe was not deleted.")
			os.Exit(0)
		}

		// Make API call
		deleteProbe, err := probe.DeleteProbeRequest(projectID, probeID, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Error in deleting Probe: ", err.Error())
			os.Exit(1)
		}

		if deleteProbe.Data.DeleteProbe {
			utils.White_B.Println("\n🚀 Probe was successfully deleted.")
		} else {
			utils.White_B.Println("\n❌ Failed to delete Probe. Please check if the ID is correct or not.")
		}
	},
}

func init() {
	DeleteCmd.AddCommand(probeCmd)

	probeCmd.Flags().String("project-id", "", "Set the project-id to delete Probe for the particular project. To see the projects, apply litmusctl get projects")
	probeCmd.Flags().String("probe-id", "", "Set the probe-id to delete that particular probe. To see the probes, apply litmusctl get probes")
}
