/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/connect"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"os"

	"github.com/spf13/cobra"
)

func printError(err error)  {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("connect called")
		defaultFileName := types.DefaultFileName

		obj, err := config.YamltoObject(defaultFileName)
		printError(err)


		if obj.CurrentUser == "" || obj.CurrentAccount == "" {
			fmt.Println("Current user or current account is not set")
			os.Exit(1)
		}

		var token string
		for _, account := range obj.Accounts {
			if account.Endpoint == obj.CurrentAccount {
				for _, user := range account.Users {
					if user.Username == obj.CurrentUser {
						token = user.Token
					}
				}
			}
		}

		var credentials = types.Credentials{
			Username: obj.CurrentUser,
			Token: token,
			Endpoint: obj.CurrentAccount,
		}

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		printError(err)

		var newAgent types.Agent

		if nonInteractive {
			newAgent.ProjectId, err = cmd.Flags().GetString("project-id")
			printError(err)

			newAgent.Mode, err = cmd.Flags().GetString("installation-mode")
			printError(err)

			newAgent.AgentName, err = cmd.Flags().GetString("agent-name")
			printError(err)

			newAgent.Description, err = cmd.Flags().GetString("agent-description")
			printError(err)


			newAgent.PlatformName, err = cmd.Flags().GetString("platform-name")
			printError(err)


			newAgent.ClusterType, err = cmd.Flags().GetString("cluster-type")
			printError(err)


			newAgent.Namespace, err = cmd.Flags().GetString("namespace")
			printError(err)

			newAgent.ServiceAccount, err = cmd.Flags().GetString("service-account")
			printError(err)

			newAgent.NsExists, err = cmd.Flags().GetBool("ns-exists")
			printError(err)

			newAgent.SAExists, err = cmd.Flags().GetBool("sa-exists")
			printError(err)

			if newAgent.Mode == "" {
				newAgent.Mode = utils.DefaultMode
			}

			// Check if user has sufficient permissions based on mode
			fmt.Println("\n🏃 Running prerequisites check....")
			connect.ValidateSAPermissions(newAgent.Mode, &kubeconfig)

			agent, err := apis.GetAgentList(credentials, newAgent.ProjectId)
			printError(err)

			// Duplicate agent check
			var isAgentExist = false
			for i := range agent.Data.GetAgent {
				if newAgent.AgentName == agent.Data.GetAgent[i].AgentName {
					fmt.Println(agent.Data.GetAgent[i].AgentName)
					isAgentExist = true
				}
			}

			fmt.Println(isAgentExist)
			if isAgentExist {
				fmt.Println("🚫 Agent with the given name already exists.")
				// Print agent list if existing agent name is entered twice
				fmt.Print("\n📘 Connected agents list -----------\n\n")

				for i := range agent.Data.GetAgent {
					fmt.Println("-", agent.Data.GetAgent[i].AgentName)
				}

				fmt.Println("\n-------------------------------------")

				fmt.Println("❗ Please enter a different name.")
				os.Exit(1)
			}





		} else {

			userDetails, err := apis.GetProjectDetails(credentials)
			printError(err)

			// Fetch project id
			projectID := connect.GetProjectID(userDetails)

			modeType := connect.GetModeType()

			// Check if user has sufficient permissions based on mode
			fmt.Println("\n🏃 Running prerequisites check....")
			connect.ValidateSAPermissions(modeType, &kubeconfig)
			newAgent, err = connect.GetAgentDetails(modeType, projectID, credentials, &kubeconfig)
			printError(err)

			newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace, &kubeconfig)
			newAgent.Mode = modeType
			connect.Summary(newAgent, &kubeconfig)

			connect.ConfirmInstallation()

		}

		agent, cerror := apis.ConnectAgent(newAgent, credentials)
		if cerror != nil {
			fmt.Printf("\n❌ Agent connection failed: [%s]\n", cerror.Error())
			os.Exit(1)
		}

		path := fmt.Sprintf("%s/%s/%s.yaml", credentials.Endpoint, utils.ChaosYamlPath, agent.Data.UserAgentReg.Token)
		fmt.Println("Applying YAML:\n", path)

		// Print error message in case Data field is null in response
		if (agent.Data == apis.AgentConnect{}) {
			fmt.Printf("\n🚫 Agent connection failed: [%s]\n", agent.Errors[0].Message)
			os.Exit(1)
		}
		//Apply agent connection yaml
		yamlOutput, yerror := utils.ApplyYaml(agent.Data.UserAgentReg.Token, credentials.Endpoint, utils.ChaosYamlPath, kubeconfig)
		if yerror != nil {
			fmt.Printf("\n❌ Failed in applying connection yaml: [%s]\n", yerror.Error())
			os.Exit(1)
		}
		fmt.Println("\n", yamlOutput)
		// Watch subscriber pod status
		k8s.WatchPod(newAgent.Namespace, utils.ChaosAgentLabel, &kubeconfig)
		fmt.Println("\n🚀 Agent Connection Successful!! 🎉")
		fmt.Println("👉 Litmus agents can be accessed here: " + fmt.Sprintf("%s/%s", credentials.Endpoint, utils.ChaosAgentPath))
	},
}

func init() {
	agentCmd.AddCommand(connectCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// connectCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	connectCmd.Flags().BoolP("non-interactive", "n", false, "Help message for toggle")
	connectCmd.Flags().StringP("kubeconfig", "k", "", "Help message for toggle")
	connectCmd.Flags().String("project-id", "","Help message for toggle")
	connectCmd.Flags().String("installation-mode", "","Help message for toggle")
	connectCmd.Flags().String("agent-name", "","Help message for toggle")
	connectCmd.Flags().String("agent-description", "","Help message for toggle")
	connectCmd.Flags().String("platform-name", "","Help message for toggle")
	connectCmd.Flags().String("cluster-type", "","Help message for toggle")
	connectCmd.Flags().String("namespace", "","Help message for toggle")
	connectCmd.Flags().String("service-account", "","Help message for toggle")
	connectCmd.Flags().Bool("ns-exists", false,"Help message for toggle")
	connectCmd.Flags().Bool("sa-exists", false,"Help message for toggle")

}
