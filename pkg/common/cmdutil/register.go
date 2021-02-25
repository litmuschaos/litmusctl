package chaos

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/common"
	"github.com/litmuschaos/litmusctl/pkg/common/k8s"
	"github.com/litmuschaos/litmusctl/pkg/constants"
)

func Register(t common.Token, c common.Credentials) {
	// Fetch project details
	_, uErr := GetProjectDetails(t, c, "chaos")
	if uErr != nil {
		fmt.Printf("\n❌ Fetching project details failed: [%s]", uErr)
		os.Exit(1)
	}
	// Fetch project id
	//pid := GetProject(user)
	// TODO: change this
	pid := "991043fc-f53b-446a-862e-29be55ba0af7"
	// Get mode of installation as input
	mode := GetMode()
	// Check if user has sufficient permissions based on mode
	fmt.Println("\n🏃 Running prerequisites check....")
	k8s.ValidateSAPermissions(mode)
	// Get agent details as input
	newAgent := GetAgentDetails(pid, t, c)
	newAgent.Mode = mode
	// Get service account as input
	newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace)
	// Display details of agent to be connected
	common.Summary(newAgent, "chaos")
	// Confirm before connecting the agent
	common.Confirm()
	// Register agent
	agent, cerror := RegisterAgent(newAgent, t, c)
	if cerror != nil {
		fmt.Printf("\n❌ Agent registration failed: [%s]\n", cerror.Error())
		os.Exit(1)
	}
	path := fmt.Sprintf("%s/%s/%s.yaml", c.Host, constants.ChaosYamlPath, agent.Data.UserAgentReg.Token)
	fmt.Println(path)
	// Print error message in case Data field is null in response
	if (agent.Data == AgentRegister{}) {
		fmt.Printf("\n🚫 Agent registration failed: [%s]\n", agent.Errors[0].Message)
		os.Exit(1)
	}
	// Apply agent registration yaml
	yamlOutput, yerror := common.ApplyYaml(agent.Data.UserAgentReg.Token, c, constants.ChaosYamlPath)
	if yerror != nil {
		fmt.Printf("\n❌ Failed in applying registration yaml: [%s]\n", yerror.Error())
		os.Exit(1)
	}
	fmt.Println("\n", yamlOutput)
	// Watch subscriber pod status
	k8s.WatchPod(newAgent.Namespace, constants.ChaosAgentLabel)
	fmt.Println("\n🚀 Agent Registration Successful!! 🎉")
	fmt.Println("👉 Kubera agents can be accessed here: " + fmt.Sprintf("%s/%s", c.Host, constants.ChaosAgentPath))
}
