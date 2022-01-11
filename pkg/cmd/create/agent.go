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
	"os"

	"github.com/litmuschaos/litmusctl/pkg/agent"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use: "agent",
	Short: `Create an external agent.
	Example(s):
	#create an agent
	litmusctl create agent --agent-name="new-agent" --non-interactive

	#create an agent within a project
	litmusctl create agent --agent-name="new-agent" --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --non-interactive
	
	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		utils.PrintError(err)

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		utils.PrintError(err)

		var newAgent types.Agent

		newAgent.ProjectId, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if newAgent.ProjectId == "" {
			userDetails, err := apis.GetProjectDetails(credentials)
			utils.PrintError(err)

			var (
				userID        = userDetails.Data.ID
				projectExists = false
			)

		outerloop:
			for _, project := range userDetails.Data.Projects {
				for _, member := range project.Members {
					if (member.UserID == userID) && (member.Role == "Owner" || member.Role == "Editor") {
						projectExists = true
						break outerloop
					}
				}
			}

			if !projectExists {
				utils.White_B.Print("Creating a random project...")
				newAgent.ProjectId = agent.CreateRandomProject(credentials)
			}
		}

		if nonInteractive {

			newAgent.Mode, err = cmd.Flags().GetString("installation-mode")
			utils.PrintError(err)

			if newAgent.Mode == "" {
				utils.Red.Print("Error: --installation-mode flag is empty")
				os.Exit(1)
			}

			newAgent.AgentName, err = cmd.Flags().GetString("agent-name")
			utils.PrintError(err)

			newAgent.SkipSSL, err = cmd.Flags().GetBool("skip-agent-ssl")
			utils.PrintError(err)

			if newAgent.AgentName == "" {
				utils.Red.Print("Error: --agent-name flag is empty")
				os.Exit(1)
			}

			newAgent.Description, err = cmd.Flags().GetString("agent-description")
			utils.PrintError(err)

			newAgent.PlatformName, err = cmd.Flags().GetString("platform-name")
			utils.PrintError(err)

			if newAgent.PlatformName == "" {
				utils.Red.Print("Error: --platform-name flag is empty")
				os.Exit(1)
			}

			newAgent.ClusterType, err = cmd.Flags().GetString("cluster-type")
			utils.PrintError(err)
			if newAgent.ClusterType == "" {
				utils.Red.Print("Error: --cluster-type flag is empty")
				os.Exit(1)
			}

			newAgent.NodeSelector, err = cmd.Flags().GetString("node-selector")
			utils.PrintError(err)
			if newAgent.NodeSelector != "" {
				if ok := utils.CheckKeyValueFormat(newAgent.NodeSelector); !ok {
					os.Exit(1)
				}
			}

			toleration, err := cmd.Flags().GetString("tolerations")
			utils.PrintError(err)

			if toleration != "" {
				var tolerations []types.Toleration
				err := json.Unmarshal([]byte(toleration), &tolerations)
				utils.PrintError(err)

				str := "["
				for _, tol := range tolerations {
					str += "{"
					if tol.TolerationSeconds > 0 {
						str += "tolerationSeconds: " + fmt.Sprint(tol.TolerationSeconds) + " "
					}
					if tol.Effect != "" {
						str += "effect: \\\"" + tol.Effect + "\\\" "
					}
					if tol.Key != "" {
						str += "key: \\\"" + tol.Key + "\\\" "
					}

					if tol.Value != "" {
						str += "value: \\\"" + tol.Value + "\\\" "
					}

					if tol.Operator != "" {
						str += "operator : \\\"" + tol.Operator + "\\\" "
					}

					str += " }"
				}
				str += "]"

				newAgent.Tolerations = str
			}

			newAgent.Namespace, err = cmd.Flags().GetString("namespace")
			utils.PrintError(err)

			newAgent.ServiceAccount, err = cmd.Flags().GetString("service-account")
			utils.PrintError(err)

			newAgent.NsExists, err = cmd.Flags().GetBool("ns-exists")
			utils.PrintError(err)

			newAgent.SAExists, err = cmd.Flags().GetBool("sa-exists")
			utils.PrintError(err)

			if newAgent.Mode == "" {
				newAgent.Mode = utils.DefaultMode
			}

			if newAgent.ProjectId == "" {
				utils.Red.Println("Error: --project-id flag is empty")
				os.Exit(1)
			}

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(newAgent.Mode, &kubeconfig)

			agents, err := apis.GetAgentList(credentials, newAgent.ProjectId)
			utils.PrintError(err)

			// Duplicate agent check
			var isAgentExist = false
			for i := range agents.Data.GetAgent {
				if newAgent.AgentName == agents.Data.GetAgent[i].AgentName {
					utils.White_B.Print(agents.Data.GetAgent[i].AgentName)
					isAgentExist = true
				}
			}

			if isAgentExist {
				agent.PrintExistingAgents(agents)
				os.Exit(1)
			}

		} else {
			userDetails, err := apis.GetProjectDetails(credentials)
			utils.PrintError(err)

			if newAgent.ProjectId == "" {
				// Fetch project id
				newAgent.ProjectId = agent.GetProjectID(userDetails)
			}

			modeType := agent.GetModeType()

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(modeType, &kubeconfig)
			newAgent, err = agent.GetAgentDetails(modeType, newAgent.ProjectId, credentials, &kubeconfig)
			utils.PrintError(err)

			newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace, &kubeconfig)
			newAgent.Mode = modeType
		}

		agent.Summary(newAgent, &kubeconfig)

		if !nonInteractive {
			agent.ConfirmInstallation()
		}

		agent, err := apis.ConnectAgent(newAgent, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Agent connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}

		path := fmt.Sprintf("%s/%s/%s.yaml", credentials.Endpoint, utils.ChaosYamlPath, agent.Data.UserAgentReg.Token)
		utils.White_B.Print("Applying YAML:\n", path)

		// Print error message in case Data field is null in response
		if (agent.Data == apis.AgentConnect{}) {
			utils.White_B.Print("\n🚫 Agent connection failed: " + agent.Errors[0].Message + "\n")
			os.Exit(1)
		}

		//Apply agent connection yaml
		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    agent.Data.UserAgentReg.Token,
			Endpoint: credentials.Endpoint,
			YamlPath: utils.ChaosYamlPath,
		}, kubeconfig, false)
		if err != nil {
			utils.White_B.Print("\n❌ Failed in applying connection yaml: \n" + yamlOutput)
			os.Exit(1)
		}

		utils.White_B.Print("\n", yamlOutput)

		// Watch subscriber pod status
		k8s.WatchPod(k8s.WatchPodParams{Namespace: newAgent.Namespace, Label: utils.ChaosAgentLabel}, &kubeconfig)

		utils.White_B.Println("\n🚀 Agent Connection Successful!! 🎉")
		utils.White_B.Println("👉 Litmus agents can be accessed here: " + fmt.Sprintf("%s/%s", credentials.Endpoint, utils.ChaosAgentPath))
	},
}

func init() {
	CreateCmd.AddCommand(agentCmd)

	agentCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	agentCmd.Flags().StringP("kubeconfig", "k", "", "Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)")
	agentCmd.Flags().String("tolerations", "", "Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)")

	agentCmd.Flags().String("project-id", "", "Set the project-id to install agent for the particular project. To see the projects, apply litmusctl get projects")
	agentCmd.Flags().String("installation-mode", "cluster", "Set the installation mode for the kind of agent | Supported=cluster/namespace")
	agentCmd.Flags().String("agent-name", "", "Set the agent name")
	agentCmd.Flags().String("agent-description", "---", "Set the agent description")
	agentCmd.Flags().String("platform-name", "Others", "Set the platform name. Supported- AWS/GKE/Openshift/Rancher/Others")
	agentCmd.Flags().String("cluster-type", "external", "Set the cluster-type to external for external agents | Supported=external/internal")
	agentCmd.Flags().String("node-selector", "", "Set the node-selector for agent components | Format: \"key1=value1,key2=value2\")")
	agentCmd.Flags().String("namespace", "litmus", "Set the namespace for the agent installation")
	agentCmd.Flags().String("service-account", "litmus", "Set the service account to be used by the agent")
	agentCmd.Flags().Bool("skip-agent-ssl", false, "Set whether agent will skip ssl/tls check (can be used for self-signed certs, if cert is not provided in portal)")
	agentCmd.Flags().Bool("ns-exists", false, "Set the --ns-exists=false if the namespace mentioned in the --namespace flag is not existed else set it to --ns-exists=true | Note: Always set the boolean flag as --ns-exists=Boolean")
	agentCmd.Flags().Bool("sa-exists", false, "Set the --sa-exists=false if the service-account mentioned in the --service-account flag is not existed else set it to --sa-exists=true | Note: Always set the boolean flag as --sa-exists=Boolean\"\n")
}
