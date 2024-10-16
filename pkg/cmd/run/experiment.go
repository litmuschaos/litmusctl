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
package run

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

// Define the necessary struct to capture the nested runnerPod field
type Execution struct {
	Namespace string          `json:"namespace"`
	Nodes     map[string]Node `json:"nodes"`
}

type Node struct {
	ChaosData  *ChaosData `json:"chaosData,omitempty"`
	Name       string     `json:"name"`
	Type       string     `json:"type"`
	Phase      string     `json:"phase"`
	StartedAt  string     `json:"startedAt"`
	FinishedAt string     `json:"finishedAt"`
}

type ChaosData struct {
	Namespace string `json:"namespace"`
	RunnerPod string `json:"runnerPod"`
}

func getEmojiForPhase(phase string) string {
	switch phase {
	case "Pending":
		return "⏳"
	case "Running":
		return "🏃"
	case "Succeeded":
		return "✅"
	case "Skipped":
		return "⤵️"
	case "Failed":
		return "❗"
	case "Error":
		return "❌"
	case "Omitted":
		return "🚫"
	case "Completed":
		return "🏁"
	default:
		return "❓"
	}
}

// logNodeDetails logs the node details in a standardized format
func logNodeDetails(node Node, prefix string) {
	utils.Cyan.Printf("\n🚀 %s: %s", prefix, node.Type)
	info := fmt.Sprintf(" %s - Phase: %s %s | ⏰ Started At: %s",
		node.Name, node.Phase, getEmojiForPhase(node.Phase), utils.FormatTimeStamp(node.StartedAt))

	// Log finished time if available
	if node.FinishedAt != "" {
		info += fmt.Sprintf(" | ⏰ Finished At: %s", utils.FormatTimeStamp(node.FinishedAt))
	}

	utils.Green.Println(info)
}

// experimentCmd represents the project command
var experimentCmd = &cobra.Command{
	Use: "chaos-experiment",
	Short: `Create a Chaos Experiment
	Example:
	#Save a Chaos Experiment
	litmusctl run chaos-experiment --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --experiment-id="1c9c5801-8789-4ac9-bf5f-32649b707a5c"

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// Fetch user credentials
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		pid, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		// Handle blank input for project ID
		if pid == "" {
			utils.White_B.Print("\nEnter the Project ID: ")
			fmt.Scanln(&pid)

			if pid == "" {
				utils.Red.Println("⛔ Project ID can't be empty!!")
				os.Exit(1)
			}
		}

		eid, err := cmd.Flags().GetString("experiment-id")
		utils.PrintError(err)

		// Handle blank input for Chaos Experiment ID
		if eid == "" {
			utils.White_B.Print("\nEnter the Chaos Experiment ID: ")
			fmt.Scanln(&eid)

			if eid == "" {
				utils.Red.Println("⛔ Chaos Experiment ID can't be empty!!")
				os.Exit(1)
			}
		}

		// Perform authorization
		userDetails, err := apis.GetProjectDetails(credentials)
		utils.PrintError(err)
		var editAccess = false
		var project apis.Project
		for _, p := range userDetails.Data.Projects {
			if p.ID == pid {
				project = p
			}
		}
		for _, member := range project.Members {
			if (member.UserID == userDetails.Data.ID) && (member.Role == utils.MemberOwnerRole || member.Role == utils.MemberEditorRole) {
				editAccess = true
			}
		}
		if !editAccess {
			utils.Red.Println("⛔ User doesn't have edit access to the project!!")
			os.Exit(1)
		}

		// Make API call
		runExperiment, err := experiment.RunExperiment(pid, eid, credentials)
		if err != nil {
			if (runExperiment.Data == experiment.RunExperimentData{}) {
				if strings.Contains(err.Error(), "multiple run errors") {
					utils.Red.Println("\n❌ Chaos Experiment already exists")
					os.Exit(1)
				}
				if strings.Contains(err.Error(), "no documents in result") {
					utils.Red.Println("❌ The specified Project ID or Chaos Infrastructure ID doesn't exist.")
					os.Exit(1)
				} else {
					utils.White_B.Print("\n❌ Failed to run chaos experiment: " + err.Error())
					os.Exit(1)
				}
			}
		}

		//Successful run
		utils.White_B.Println("\n🚀 Chaos Experiment running successfully 🎉")

		// Check if we need to stream pod logs
		streamLogs, err := cmd.Flags().GetBool("stream-logs")
		utils.PrintError(err)

		if streamLogs {
			exp, err := experiment.GetExperimentRun(project.ID, runExperiment.Data.RunExperimentDetails.NotifyID, credentials)
			if err != nil {
				utils.Red.Print("\n❌ Failed to fetch experiment: " + err.Error())
			}
			// Create a map to keep track of seen nodes
			seenNodes := make(map[string]string)

			for {
				exp, err = experiment.GetExperimentRun(project.ID, runExperiment.Data.RunExperimentDetails.NotifyID, credentials)
				if err != nil {
					utils.Red.Print("\n❌ Failed to fetch experiment: " + err.Error())
				}

				var execution Execution
				err = json.Unmarshal([]byte(exp.Data.ExperimentRunDetails.ExecutionData), &execution)
				if err != nil {
					utils.Red.Print("\n❌ Error unmarshalling JSON " + err.Error())
				}

				for nodeName, node := range execution.Nodes {
					// Check if node is new
					if _, found := seenNodes[nodeName]; !found {
						// Log new node found
						logNodeDetails(node, "New node detected")

						// Mark the node as seen
						seenNodes[nodeName] = node.Phase

					} else if node.Phase != seenNodes[nodeName] {
						// Node already seen, but phase has changed
						logNodeDetails(node, "Node phase updated")

						// Check and log runnerPod if available
						if node.ChaosData != nil && node.ChaosData.RunnerPod != "" && node.Phase == "Completed" {
							podLogReq := experiment.PodLogRequest{
								InfraID:         exp.Data.ExperimentRunDetails.Infra.InfraID,
								ExperimentRunID: exp.Data.ExperimentRunDetails.ExperimentRunID,
								PodNamespace:    execution.Namespace,
								PodType:         node.Type,
								RunnerPod:       node.ChaosData.RunnerPod,
								ChaosNamespace:  node.ChaosData.Namespace,
							}
							podLogRes, err := experiment.GetPodLogs(podLogReq, credentials)
							if err != nil {
								utils.White_B.Print("\n❌ Failed to fetch logs: " + err.Error())
							}
							utils.White_B.Println("\n🚀PodLogs: \n " + podLogRes.Data.GetPodLog.Log)
						}

						// Update the node phase
						seenNodes[nodeName] = node.Phase
					}
				}

				if exp.Data.ExperimentRunDetails.Phase != "Running" {
					break
				}
				time.Sleep(time.Second * 1)
			}
		}
	},
}

func init() {
	RunCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to create Chaos Experiment for the particular project. To see the projects, apply litmusctl get projects")
	experimentCmd.Flags().String("experiment-id", "", "Set the environment-id to create Chaos Experiment for the particular Chaos Infrastructure. To see the Chaos Infrastructures, apply litmusctl get chaos-infra")
	experimentCmd.Flags().Bool("stream-logs", false, "Set the --stream-logs=true if you want to fetch and stream logs from the Pod\"\n")
}
