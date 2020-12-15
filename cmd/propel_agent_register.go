/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"os"

	"github.com/mayadata-io/kuberactl/pkg/types/propel"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var propelAgentRegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register Kubera Propel agent",
	Long:  `Register registers the agent to Kubera Propel`,
	Run: func(cmd *cobra.Command, args []string) {

		var c Credentials
		var pErr error
		fmt.Println("🔥 Registering Kubera Enterprise agent")
		fmt.Println("\n📶 Please enter Kubera Enterprise details --")
		// Get Kubera Enterprise URL as input
		c.Host, pErr = getPortalURL()
		if pErr != nil {
			fmt.Printf("\n❌ URL parsing failed: [%s]", pErr.Error())
			os.Exit(1)
		}
		// Get username as input
		c.Username = getUsername()
		// Get password as input
		c.Password = getPassword()
		// Fetch authorization token
		t := login(c)
		// Get LaunchProduct token
		productToken, err := LaunchProduct(t, c, "Propel")
		if err != nil {
			fmt.Printf("\n❌ Fetching LaunchProduct query failed: [%s]", err)
			os.Exit(1)
		}
		// Replace AccessToken with LaunchProduct token
		t.AccessToken = productToken.Data.LaunchProduct
		// Fetch project details
		user, uErr := ListPropelProjects(t, c)
		if uErr != nil {
			fmt.Printf("\n❌ Fetching project details failed: [%s]", uErr)
			os.Exit(1)
		}
		if len(user.Errors) != 0 {
			fmt.Printf("\n❌ Fetching project details failed: [%s]", user.Errors[0].Message)
			os.Exit(1)
		}
		// Fetch project id
		pid := SelectPropelProject(user)
		// Get agent details as input
		newAgent := GetPropelAgentDetails(pid, t, c)
		// Display details of agent to be connected
		Summary(newAgent, "propel")
		// Confirm before connecting the agent
		confirm()
		// Register agent
		agent, cerror := RegisterPropelAgent(newAgent, t, c)
		if cerror != nil {
			fmt.Printf("\n❌ Agent registration failed: [%s]\n", cerror.Error())
			os.Exit(1)
		}
		// Print error message in case Data field is null in response
		if (agent.AgentData == propel.AgentData{}) {
			fmt.Printf("\n🚫 Agent registration failed: [%s]\n", agent.Errors[0].Message)
			os.Exit(1)
		}
		// Apply agent registration yaml
		yamlOutput, yerror := ApplyYaml(agent.AgentData.AddCluster.ClusterToken, c, propelYamlPath)
		if yerror != nil {
			fmt.Printf("\n❌ Failed in applying registration yaml: [%s]\n", yerror.Error())
			os.Exit(1)
		}
		fmt.Println("\n", yamlOutput)
		// Watch subscriber pod status
		WatchPod(newAgent.Namespace, propelAgentLabel)
		fmt.Println("\n🚀 Agent Registration Successful!! 🎉")
		fmt.Println("👉 Kubera agents can be accessed here: " + fmt.Sprintf("%s/%s", c.Host, propelAgentPath))

	},
}

func init() {
	propelAgentCmd.AddCommand(propelAgentRegisterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// registerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// registerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
