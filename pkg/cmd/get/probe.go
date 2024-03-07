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

		interactiveMode, err := cmd.Flags().GetBool("non-interactive")

		var selectedItems []*models.ProbeType

		if interactiveMode == true {
			prompt := promptui.Select{
				Label: "Do you want to enable advance filter probes?",
				Items: []string{"Yes", "No"},
			}
			_, option, err := prompt.Run()
			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
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
						fmt.Printf("Prompt failed %v\n", err)
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
		utils.White_B.Fprintln(writer, "PROBE ID\t PROBE TYPE\t CREATED AT\t CREATED BY")

		for {
			writer.Flush()
			// calculating the start and end indices for the current page
			start := (page - 1) * itemsPerPage
			if start >= totalProbes {
				utils.Red.Println("No more probes to display")
				writer.Flush()
				break
			}
			end := start + itemsPerPage
			if end > totalProbes {
				end = totalProbes
			}
			for _, probe := range probes_data[start:end] {
				intTime, err := strconv.ParseInt(probe.CreatedAt, 10, 64)
				if err != nil {
					fmt.Println("Error converting CreatedAt to int64:", err)
					continue
				}
				humanTime := time.Unix(intTime, 0)
				utils.White.Fprintln(writer, probe.Name+"\t"+fmt.Sprintf("%v", probe.Type)+"\t"+probe.CreatedBy.Username+"\t"+humanTime.String())
			}
			writer.Flush()

			// Check if it's the last item or if user wants to see more
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
			// Move to the next page
			page++
		}
	},
}

func init() {
	GetCmd.AddCommand(probesCmd)

	probesCmd.Flags().String("project-id", "", "Set the project-id to get Probe from a particular project.")
	probesCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	probesCmd.Flags().String("probe-types", "", "Set the probe-types as comma separated values to filter the probes")
}
