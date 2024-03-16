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
		ProbeID, err := cmd.Flags().GetString("probe-id")

		if ProbeID == "" {
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

		} else {
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
				utils.Red.Println("Error converting UpdatedAt to int64:", err)
			}
			updatedTime := time.Unix(intUpdateTime, 0).String()
			intCreatedTime, err := strconv.ParseInt(probeGetData.CreatedAt, 10, 64)
			if err != nil {
				utils.Red.Println("Error converting CreatedAt to int64:", err)
			}
			createdTime := time.Unix(intCreatedTime, 0).String()
			utils.White_B.Fprintln(writer, "PROBE DETAILS")
			utils.White.Fprintln(writer, "PROBE ID \t", probeGetData.Name)
			utils.White.Fprintln(writer, "PROBE DESCRIPTION \t", *probeGetData.Description)
			utils.White.Fprintln(writer, "PROBE TYPE \t", probeGetData.Type)
			utils.White.Fprintln(writer, "PROBE INFRASTRUCTURE TYPE \t", probeGetData.InfrastructureType)
			if probeGetData.Type == "httpProbe" {
				utils.White.Fprintln(writer, "TIMEOUT \t", probeGetData.KubernetesHTTPProperties.ProbeTimeout)
				utils.White.Fprintln(writer, "INTERVAL \t", probeGetData.KubernetesHTTPProperties.Interval)
				if probeGetData.KubernetesHTTPProperties.Attempt != nil {
					utils.White.Fprintln(writer, "ATTEMPT \t", *probeGetData.KubernetesHTTPProperties.Attempt)
				}
				if probeGetData.KubernetesHTTPProperties.ProbePollingInterval != nil {
					utils.White.Fprintln(writer, "POLLING INTERVAL \t", *probeGetData.KubernetesHTTPProperties.ProbePollingInterval)
				}
				if probeGetData.KubernetesHTTPProperties.InitialDelay != nil {
					utils.White.Fprintln(writer, "INITIAL DELAY \t", *probeGetData.KubernetesHTTPProperties.InitialDelay)
				}
				if probeGetData.KubernetesHTTPProperties.EvaluationTimeout != nil {
					utils.White.Fprintln(writer, "EVALUATION TIMEOUT \t", *probeGetData.KubernetesHTTPProperties.EvaluationTimeout)
				}
				if probeGetData.KubernetesHTTPProperties.StopOnFailure != nil {
					utils.White.Fprintln(writer, "STOP ON FAILURE \t", *probeGetData.KubernetesHTTPProperties.StopOnFailure)
				}
				utils.White.Fprintln(writer, "URL \t", probeGetData.KubernetesHTTPProperties.URL)
				if probeGetData.KubernetesHTTPProperties.Method.Get != nil {
					utils.White.Fprintln(writer, "METHOD \t", "GET")
					utils.White.Fprintln(writer, "CRITERIA \t", probeGetData.KubernetesHTTPProperties.Method.Get.Criteria)
					utils.White.Fprintln(writer, "RESPONSE \t", probeGetData.KubernetesHTTPProperties.Method.Get.ResponseCode)
				}
				if probeGetData.KubernetesHTTPProperties.Method.Post != nil {
					utils.White.Fprintln(writer, "METHOD \t", "POST")
					utils.White.Fprintln(writer, "CONTENT-TYPE \t", probeGetData.KubernetesHTTPProperties.Method.Post.ContentType)
					if probeGetData.KubernetesHTTPProperties.Method.Post.Body != nil {
						utils.White.Fprintln(writer, "BODY \t", probeGetData.KubernetesHTTPProperties.Method.Post.Body)
					}
					if probeGetData.KubernetesHTTPProperties.Method.Post.BodyPath != nil {
						utils.White.Fprintln(writer, "BODYPATH \t", probeGetData.KubernetesHTTPProperties.Method.Post.BodyPath)
					}
					utils.White.Fprintln(writer, "CRITERIA \t", probeGetData.KubernetesHTTPProperties.Method.Post.Criteria)
					utils.White.Fprintln(writer, "RESPONSE CODE \t", probeGetData.KubernetesHTTPProperties.Method.Post.ResponseCode)
				}
				if probeGetData.KubernetesHTTPProperties.InsecureSkipVerify != nil {
					utils.White.Fprintln(writer, "INSECURE SKIP VERIFY \t", probeGetData.KubernetesHTTPProperties.InsecureSkipVerify)
				}

			}
			if probeGetData.Type == "cmdProbe" {
				utils.White.Fprintln(writer, "TIMEOUT \t", probeGetData.KubernetesCMDProperties.ProbeTimeout)
				utils.White.Fprintln(writer, "INTERVAL \t", probeGetData.KubernetesCMDProperties.Interval)
				if probeGetData.KubernetesCMDProperties.Attempt != nil {
					utils.White.Fprintln(writer, "ATTEMPT \t", *probeGetData.KubernetesCMDProperties.Attempt)
				}
				if probeGetData.KubernetesCMDProperties.ProbePollingInterval != nil {
					utils.White.Fprintln(writer, "POLLING INTERVAL \t", *probeGetData.KubernetesCMDProperties.ProbePollingInterval)
				}
				if probeGetData.KubernetesCMDProperties.InitialDelay != nil {
					utils.White.Fprintln(writer, "INITIAL DELAY \t", *probeGetData.KubernetesCMDProperties.InitialDelay)
				}
				if probeGetData.KubernetesCMDProperties.EvaluationTimeout != nil {
					utils.White.Fprintln(writer, "EVALUATION TIMEOUT \t", *probeGetData.KubernetesCMDProperties.EvaluationTimeout)
				}
				if probeGetData.KubernetesCMDProperties.StopOnFailure != nil {
					utils.White.Fprintln(writer, "STOP ON FAILURE \t", *probeGetData.KubernetesCMDProperties.StopOnFailure)
				}
				utils.White.Fprintln(writer, "Command \t", probeGetData.KubernetesCMDProperties.Command)
				if probeGetData.KubernetesCMDProperties.Comparator != nil {
					utils.White.Fprintln(writer, "COMPARATOR TYPE \t", probeGetData.KubernetesCMDProperties.Comparator.Type)
					utils.White.Fprintln(writer, "COMPARATOR VALUE \t", probeGetData.KubernetesCMDProperties.Comparator.Value)
					utils.White.Fprintln(writer, "COMPARATOR CRITERIA \t", probeGetData.KubernetesCMDProperties.Comparator.Criteria)
				}
				if probeGetData.KubernetesCMDProperties.Source != nil {

					utils.White.Fprintln(writer, "Source \t", *probeGetData.KubernetesCMDProperties.Source)
				}
			}
			if probeGetData.Type == "k8sProbe" {
				utils.White.Fprintln(writer, "TIMEOUT \t", probeGetData.K8sProperties.ProbeTimeout)
				utils.White.Fprintln(writer, "INTERVAL \t", probeGetData.K8sProperties.Interval)
				if probeGetData.K8sProperties.Attempt != nil {
					utils.White.Fprintln(writer, "ATTEMPT \t", *probeGetData.K8sProperties.Attempt)
				}
				if probeGetData.K8sProperties.ProbePollingInterval != nil {
					utils.White.Fprintln(writer, "POLLING INTERVAL \t", *probeGetData.K8sProperties.ProbePollingInterval)
				}
				if probeGetData.K8sProperties.InitialDelay != nil {
					utils.White.Fprintln(writer, "INITIAL DELAY \t", *probeGetData.K8sProperties.InitialDelay)
				}
				if probeGetData.K8sProperties.EvaluationTimeout != nil {
					utils.White.Fprintln(writer, "EVALUATION TIMEOUT \t", *probeGetData.K8sProperties.EvaluationTimeout)
				}
				if probeGetData.K8sProperties.StopOnFailure != nil {
					utils.White.Fprintln(writer, "STOP ON FAILURE \t", *probeGetData.K8sProperties.StopOnFailure)
				}
				if probeGetData.K8sProperties.Group != nil {
					utils.White.Fprintln(writer, "GROUP \t", *probeGetData.K8sProperties.Group)
				}
				utils.White.Fprintln(writer, "VERSION \t", probeGetData.K8sProperties.Version)
				utils.White.Fprintln(writer, "RESOURCE \t", probeGetData.K8sProperties.Resource)
				if probeGetData.K8sProperties.Namespace != nil {
					utils.White.Fprintln(writer, "NAMESPACE \t", *probeGetData.K8sProperties.Namespace)
				}
				if probeGetData.K8sProperties.ResourceNames != nil {
					utils.White.Fprintln(writer, "RESOURCES NAMES \t", *probeGetData.K8sProperties.ResourceNames)
				}
				if probeGetData.K8sProperties.FieldSelector != nil {
					utils.White.Fprintln(writer, "FIELD SELECTOR \t", *probeGetData.K8sProperties.FieldSelector)
				}
				if probeGetData.K8sProperties.LabelSelector != nil {
					utils.White.Fprintln(writer, "LABEL SELECTOR \t", *probeGetData.K8sProperties.LabelSelector)
				}
				utils.White.Fprintln(writer, "OPERATION \t", probeGetData.K8sProperties.Operation)

			}
			if probeGetData.Type == "promProbe" {
				utils.White.Fprintln(writer, "TIMEOUT \t", probeGetData.PromProperties.ProbeTimeout)
				utils.White.Fprintln(writer, "INTERVAL \t", probeGetData.PromProperties.Interval)
				if probeGetData.PromProperties.Attempt != nil {
					utils.White.Fprintln(writer, "ATTEMPT \t", *probeGetData.PromProperties.Attempt)
				}
				if probeGetData.PromProperties.ProbePollingInterval != nil {
					utils.White.Fprintln(writer, "POLLING INTERVAL \t", *probeGetData.PromProperties.ProbePollingInterval)
				}
				if probeGetData.PromProperties.InitialDelay != nil {
					utils.White.Fprintln(writer, "INITIAL DELAY \t", *probeGetData.PromProperties.InitialDelay)
				}
				if probeGetData.PromProperties.EvaluationTimeout != nil {
					utils.White.Fprintln(writer, "EVALUATION TIMEOUT \t", *probeGetData.PromProperties.EvaluationTimeout)
				}
				if probeGetData.PromProperties.StopOnFailure != nil {
					utils.White.Fprintln(writer, "STOP ON FAILURE \t", *probeGetData.PromProperties.StopOnFailure)
				}
				utils.White.Fprintln(writer, "Endpoint \t", probeGetData.PromProperties.Endpoint)
				utils.White.Fprintln(writer, "Comparator Type \t", probeGetData.PromProperties.Comparator.Type)
				utils.White.Fprintln(writer, "Comparator Criteria \t", probeGetData.PromProperties.Comparator.Criteria)
				utils.White.Fprintln(writer, "Comparator Value \t", probeGetData.PromProperties.Comparator.Value)
				if probeGetData.PromProperties.Query != nil {
					utils.White.Fprintln(writer, "Query \t", *probeGetData.PromProperties.Query)
				}
				if probeGetData.PromProperties.QueryPath != nil {
					utils.White.Fprintln(writer, "Querypath \t", *probeGetData.PromProperties.QueryPath)
				}

			}
			utils.White.Fprintln(writer, "CREATED AT\t", createdTime)
			utils.White.Fprintln(writer, "CREATED BY\t", probeGetData.CreatedBy.Username)
			utils.White.Fprintln(writer, "UPDATED AT\t", updatedTime)
			utils.White.Fprintln(writer, "UPDATED BY\t", probeGetData.UpdatedBy.Username)
			utils.White.Fprintln(writer, "TAGS\t", strings.Join(probeGetData.Tags, ", "))
			writer.Flush()
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
