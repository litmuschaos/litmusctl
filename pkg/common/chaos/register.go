package chaos

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/common"
	"github.com/litmuschaos/litmusctl/pkg/common/k8s"
	"github.com/litmuschaos/litmusctl/pkg/constants"
)

func Connect(t common.Token, c common.Credentials, kubeconfig string) {
	// Fetch project details
	user, uErr := GetProjectDetails(t, c)

	if uErr != nil {
		fmt.Printf("\n❌ Fetching project details failed: [%s]", uErr)
		os.Exit(1)
	}
	// Fetch project id
	pid := GetProject(user)

	// Get mode of installation as input
	mode := GetMode()
	// Check if user has sufficient permissions based on mode
	fmt.Println("\n🏃 Running prerequisites check....")
	k8s.ValidateSAPermissions(mode, &kubeconfig)
	// Get agent details as input
	newAgent := GetAgentDetails(pid, t, c, &kubeconfig)
	newAgent.Mode = mode
	// Get service account as input
	newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace, &kubeconfig)
	// Display details of agent to be connected
	common.Summary(newAgent, "chaos", &kubeconfig)
	// Confirm before connecting the agent
	common.Confirm()
	// Connect agent
	agent, cerror := ConnectAgent(newAgent, t, c)
	if cerror != nil {
		fmt.Printf("\n❌ Agent connection failed: [%s]\n", cerror.Error())
		os.Exit(1)
	}

	path := fmt.Sprintf("%s/%s/%s.yaml", c.Host, constants.ChaosYamlPath, agent.Data.UserAgentReg.Token)
	fmt.Println("Applying YAML:\n", path)

	// Print error message in case Data field is null in response
	if (agent.Data == AgentConnect{}) {
		fmt.Printf("\n🚫 Agent connection failed: [%s]\n", agent.Errors[0].Message)
		os.Exit(1)
	}
	//Apply agent connection yaml
	yamlOutput, yerror := common.ApplyYaml(agent.Data.UserAgentReg.Token, c, constants.ChaosYamlPath, kubeconfig)
	if yerror != nil {
		fmt.Printf("\n❌ Failed in applying connection yaml: [%s]\n", yerror.Error())
		os.Exit(1)
	}
	fmt.Println("\n", yamlOutput)
	// Watch subscriber pod status
	k8s.WatchPod(newAgent.Namespace, constants.ChaosAgentLabel, &kubeconfig)
	fmt.Println("\n🚀 Agent Connection Successful!! 🎉")
	fmt.Println("👉 Litmus agents can be accessed here: " + fmt.Sprintf("%s/%s", c.Host, constants.ChaosAgentPath))
}
