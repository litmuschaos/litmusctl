package register

import (
	"fmt"
	"os"

	common "github.com/mayadata-io/cli-utils/pkg/common"
	"github.com/mayadata-io/cli-utils/pkg/common/k8s"
	"github.com/mayadata-io/cli-utils/pkg/constants"
	util "github.com/mayadata-io/kuberactl/cmd/propel/util"
	"github.com/mayadata-io/kuberactl/core"

	//"github.com/mayadata-io/kuberactl/pkg/types/propel"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register Kubera Propel agent",
	Long:  `Register registers the agent to Kubera Propel`,
	Run: func(cmd *cobra.Command, args []string) {

		var c common.Credentials
		var pErr error
		fmt.Println("🔥 Registering Kubera Enterprise agent")
		fmt.Println("\n📶 Please enter Kubera Enterprise details --")
		// Get Kubera Enterprise URL as input
		c.Host, pErr = common.GetPortalURL()
		if pErr != nil {
			fmt.Printf("\n❌ URL parsing failed: [%s]", pErr.Error())
			os.Exit(1)
		}
		// Get username as input
		c.Username = common.GetUsername()
		// Get password as input
		c.Password = common.GetPassword()
		// Fetch authorization token
		t := common.Login(c)
		// Get LaunchProduct token
		productToken, err := core.LaunchProduct(t, c, "Propel")
		if err != nil {
			fmt.Printf("\n❌ Fetching LaunchProduct query failed: [%s]", err)
			os.Exit(1)
		}
		// Replace AccessToken with LaunchProduct token
		t.AccessToken = productToken.Data.LaunchProduct
		// Fetch project details
		user, uErr := util.ListPropelProjects(t, c)
		if uErr != nil {
			fmt.Printf("\n❌ Fetching project details failed: [%s]", uErr)
			os.Exit(1)
		}
		if len(user.Errors) != 0 {
			fmt.Printf("\n❌ Fetching project details failed: [%s]", user.Errors[0].Message)
			os.Exit(1)
		}
		// Fetch project id
		pid := util.SelectPropelProject(user)
		// Get agent details as input
		newAgent := util.GetPropelAgentDetails(pid, t, c)
		// Display details of agent to be connected
		common.Summary(newAgent, "propel")
		// Confirm before connecting the agent
		common.Confirm()
		// Register agent
		agent, cerror := util.RegisterPropelAgent(newAgent, t, c)
		if cerror != nil {
			fmt.Printf("\n❌ Agent registration failed: [%s]\n", cerror.Error())
			os.Exit(1)
		}
		// Print error message in case Data field is null in response
		if (agent.AgentData == util.AgentData{}) {
			fmt.Printf("\n🚫 Agent registration failed: [%s]\n", agent.Errors[0].Message)
			os.Exit(1)
		}
		// Apply agent registration yaml
		yamlOutput, yerror := common.ApplyYaml(agent.AgentData.AddCluster.ClusterToken, c, constants.PropelYamlPath)
		if yerror != nil {
			fmt.Printf("\n❌ Failed in applying registration yaml: [%s]\n", yerror.Error())
			os.Exit(1)
		}
		fmt.Println("\n", yamlOutput)
		// Watch subscriber pod status
		k8s.WatchPod(newAgent.Namespace, constants.PropelAgentLabel)
		fmt.Println("\n🚀 Agent Registration Successful!! 🎉")
		fmt.Println("👉 Kubera agents can be accessed here: " + fmt.Sprintf("%s/%s", c.Host, constants.PropelAgentPath))

	},
}
