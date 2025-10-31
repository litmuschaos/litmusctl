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
package get

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	apis "github.com/litmuschaos/litmusctl/pkg/apis/probe"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var probesCmd = &cobra.Command{
	Use:   "probes",
	Short: "Display list of probes",
	Long:  `Display list of probes`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var projectID string
		projectID, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if projectID == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&projectID)

			if projectID == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}
		ProbeID, err := cmd.Flags().GetString("probe-id")

		if ProbeID == "" {
			getProbeList(projectID, cmd, credentials)
		} else {
			getProbeDetails(projectID, ProbeID, credentials)

		}

	},
}

func init() {
	GetCmd.AddCommand(probesCmd)

	probesCmd.Flags().String("project-id", "", "Set the project-id to get Probe from a particular project.")
	probesCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	probesCmd.Flags().String("probe-types", "", "Set the probe-types as comma separated values to filter the probes")
	probesCmd.Flags().String("probe-id", "", "Set the probe-details to the ID of probe for getting all the details related to the probe.")
}

func getProbeList(projectID string, cmd *cobra.Command, credentials types.Credentials) {

	// calls the probeList endpoint for the list of probes
	var selectedItems []*models.ProbeType
	NoninteractiveMode, err := cmd.Flags().GetBool("non-interactive")
	utils.PrintError(err)

	if NoninteractiveMode == false {
		prompt := promptui.Select{
			Label: "Do you want to enable advance filter probes?",
			Items: []string{"Yes", "No"},
		}
		_, option, err := prompt.Run()
		if err != nil {
			utils.PrintFormattedError("Prompt failed", err)
			return
		}
		fmt.Printf("You chose %q\n", option)

		if option == "Yes" {
			items := []models.ProbeType{"httpProbe", "cmdProbe", "promProbe", "k8sProbe", "done"}
			for {
				prompt := promptui.Select{
					Label: "Select ProbeType",
					Items: items,
					Templates: &promptui.SelectTemplates{
						Active:   `▸ {{ . | cyan }}`,
						Inactive: `  {{ . | white }}`,
						Selected: `{{ "✔" | green }} {{ . | bold }}`,
					},
				}

				selectedIndex, result, err := prompt.Run()
				if err != nil {
					utils.PrintFormattedError("Prompt failed", err)
					os.Exit(1)
				}

				if items[selectedIndex] == "done" {
					break
				}

				final := models.ProbeType(result)
				selectedItems = append(selectedItems, &final)
				items = append(items[:selectedIndex], items[selectedIndex+1:]...)

			}

			fmt.Printf("Selected Probe Types: %v\n", selectedItems)
		}
	} else {
		var probeTypes string
		probeTypes, err = cmd.Flags().GetString("probe-types")
		values := strings.Split(probeTypes, ",")
		for _, value := range values {
			probeType := models.ProbeType(value)
			selectedItems = append(selectedItems, &probeType)
		}
	}

	probes_get, _ := apis.ListProbeRequest(projectID, selectedItems, credentials)
	probes_data := probes_get.Data.Probes

	itemsPerPage := 5
	page := 1
	totalProbes := len(probes_data)

	writer := tabwriter.NewWriter(os.Stdout, 8, 8, 8, '\t', tabwriter.AlignRight)
	utils.White_B.Fprintln(writer, "PROBE ID\t PROBE TYPE\t REFERENCED BY\t CREATED BY\t CREATED AT")

	for {
		writer.Flush()
		start := (page - 1) * itemsPerPage
		if start >= totalProbes {
			utils.Red.Println("No more probes to display")
			break
		}
		end := start + itemsPerPage
		if end > totalProbes {
			end = totalProbes
		}
		for _, probe := range probes_data[start:end] {
			intTime, err := strconv.ParseInt(probe.CreatedAt, 10, 64)
			if err != nil {
				utils.PrintFormattedError("Error converting CreatedAt to int64", err)
				continue
			}
			humanTime := time.Unix(intTime, 0)
			var probeReferencedBy int
			if probe.ReferencedBy != nil {
				probeReferencedBy = *probe.ReferencedBy
			}

			utils.White.Fprintln(writer, probe.Name+"\t"+fmt.Sprintf("%v", probe.Type)+"\t"+fmt.Sprintf("%d", probeReferencedBy)+"\t"+probe.CreatedBy.Username+"\t"+humanTime.String())
		}
		writer.Flush()

		paginationPrompt := promptui.Prompt{
			Label:     "Press Enter to show more probes (or type 'q' to quit)",
			AllowEdit: true,
			Default:   "",
		}

		userInput, err := paginationPrompt.Run()
		utils.PrintError(err)

		if userInput == "q" {
			break
		}
		page++
	}

}
func getProbeDetails(projectID, ProbeID string, credentials types.Credentials) {
	//call the probe get endpoint to get the probes details
	probeGet, err := apis.GetProbeRequest(projectID, ProbeID, credentials)
	if err != nil {
		if strings.Contains(err.Error(), "permission_denied") {
			utils.Red.Println("❌ You don't have enough permissions to access this resource.")
			os.Exit(1)
		} else {
			utils.PrintError(err)
			os.Exit(1)
		}
	}
	probeGetData := probeGet.Data.GetProbe
	writer := tabwriter.NewWriter(os.Stdout, 30, 8, 2, '\t', tabwriter.AlignRight)
	intUpdateTime, err := strconv.ParseInt(probeGetData.UpdatedAt, 10, 64)
	if err != nil {
		utils.PrintFormattedError("Error converting UpdatedAt to int64", err)
	}
	updatedTime := time.Unix(intUpdateTime, 0).String()
	intCreatedTime, err := strconv.ParseInt(probeGetData.CreatedAt, 10, 64)
	if err != nil {
		utils.PrintFormattedError("Error converting CreatedAt to int64", err)
	}
	createdTime := time.Unix(intCreatedTime, 0).String()
	utils.White_B.Fprintln(writer, "PROBE DETAILS")
	utils.White.Fprintln(writer, "PROBE ID \t", probeGetData.Name)
	utils.White.Fprintln(writer, "PROBE DESCRIPTION \t", *probeGetData.Description)
	utils.White.Fprintln(writer, "PROBE TYPE \t", probeGetData.Type)
	utils.White.Fprintln(writer, "PROBE INFRASTRUCTURE TYPE \t", probeGetData.InfrastructureType)

	switch probeGetData.Type {
	case "httpProbe":
		printHTTPProbeDetails(writer, probeGetData)
	case "cmdProbe":
		printCmdProbeDetails(writer, probeGetData)
	case "k8sProbe":
		printK8sProbeDetails(writer, probeGetData)
	case "promProbe":
		printPromProbeDetails(writer, probeGetData)
	}

	utils.White.Fprintln(writer, "CREATED AT\t", createdTime)
	utils.White.Fprintln(writer, "CREATED BY\t", probeGetData.CreatedBy.Username)
	utils.White.Fprintln(writer, "UPDATED AT\t", updatedTime)
	utils.White.Fprintln(writer, "UPDATED BY\t", probeGetData.UpdatedBy.Username)
	utils.White.Fprintln(writer, "TAGS\t", strings.Join(probeGetData.Tags, ", "))
	if probeGetData.ReferencedBy != nil {
		utils.White.Fprintln(writer, "REFERENCED BY\t", *probeGetData.ReferencedBy)
	}
	writer.Flush()
}

func printHTTPProbeDetails(writer *tabwriter.Writer, probeData models.Probe) {
	utils.White.Fprintln(writer, "TIMEOUT \t", probeData.KubernetesHTTPProperties.ProbeTimeout)
	utils.White.Fprintln(writer, "INTERVAL \t", probeData.KubernetesHTTPProperties.Interval)
	printOptionalIntProperty(writer, "ATTEMPT", probeData.KubernetesHTTPProperties.Attempt)
	printOptionalIntProperty(writer, "RETRY", probeData.KubernetesHTTPProperties.Retry)
	printOptionalProperty(writer, "POLLING INTERVAL", probeData.KubernetesHTTPProperties.ProbePollingInterval)
	printOptionalProperty(writer, "INITIAL DELAY", probeData.KubernetesHTTPProperties.InitialDelay)
	printOptionalProperty(writer, "EVALUATION TIMEOUT", probeData.KubernetesHTTPProperties.EvaluationTimeout)
	printOptionalBoolProperty(writer, "STOP ON FAILURE", probeData.KubernetesHTTPProperties.StopOnFailure)
}
func printCmdProbeDetails(writer *tabwriter.Writer, probeData models.Probe) {
	utils.White.Fprintln(writer, "TIMEOUT \t", probeData.KubernetesCMDProperties.ProbeTimeout)
	utils.White.Fprintln(writer, "INTERVAL \t", probeData.KubernetesCMDProperties.Interval)
	printOptionalIntProperty(writer, "ATTEMPT", probeData.KubernetesCMDProperties.Attempt)
	printOptionalIntProperty(writer, "RETRY", probeData.KubernetesCMDProperties.Retry)
	printOptionalProperty(writer, "POLLING INTERVAL", probeData.KubernetesCMDProperties.ProbePollingInterval)
	printOptionalProperty(writer, "INITIAL DELAY", probeData.KubernetesCMDProperties.InitialDelay)
	printOptionalProperty(writer, "EVALUATION TIMEOUT", probeData.KubernetesCMDProperties.EvaluationTimeout)
	printOptionalBoolProperty(writer, "STOP ON FAILURE", probeData.KubernetesCMDProperties.StopOnFailure)
}
func printK8sProbeDetails(writer *tabwriter.Writer, probeData models.Probe) {
	utils.White.Fprintln(writer, "TIMEOUT \t", probeData.K8sProperties.ProbeTimeout)
	utils.White.Fprintln(writer, "INTERVAL \t", probeData.K8sProperties.Interval)
	printOptionalIntProperty(writer, "ATTEMPT", probeData.K8sProperties.Attempt)
	printOptionalIntProperty(writer, "RETRY", probeData.K8sProperties.Retry)
	printOptionalProperty(writer, "POLLING INTERVAL", probeData.K8sProperties.ProbePollingInterval)
	printOptionalProperty(writer, "INITIAL DELAY", probeData.K8sProperties.InitialDelay)
	printOptionalProperty(writer, "EVALUATION TIMEOUT", probeData.K8sProperties.EvaluationTimeout)
	printOptionalBoolProperty(writer, "STOP ON FAILURE", probeData.K8sProperties.StopOnFailure)
}
func printPromProbeDetails(writer *tabwriter.Writer, probeData models.Probe) {
	utils.White.Fprintln(writer, "TIMEOUT \t", probeData.PromProperties.ProbeTimeout)
	utils.White.Fprintln(writer, "INTERVAL \t", probeData.PromProperties.Interval)
	printOptionalIntProperty(writer, "ATTEMPT", probeData.PromProperties.Attempt)
	printOptionalIntProperty(writer, "RETRY", probeData.PromProperties.Retry)
	printOptionalProperty(writer, "POLLING INTERVAL", probeData.PromProperties.ProbePollingInterval)
	printOptionalProperty(writer, "INITIAL DELAY", probeData.PromProperties.InitialDelay)
	printOptionalProperty(writer, "EVALUATION TIMEOUT", probeData.PromProperties.EvaluationTimeout)
	printOptionalBoolProperty(writer, "STOP ON FAILURE", probeData.PromProperties.StopOnFailure)
}
func printOptionalProperty(writer *tabwriter.Writer, name string, value *string) {
	if value != nil {
		utils.White.Fprintln(writer, name+"\t", *value)
	}
}
func printOptionalIntProperty(writer *tabwriter.Writer, name string, value *int) {
	if value != nil {
		utils.White.Fprintln(writer, name+"\t", *value)
	}
}
func printOptionalBoolProperty(writer *tabwriter.Writer, name string, value *bool) {
	if value != nil {
		utils.White.Fprintln(writer, name+"\t", *value)
	}
}
