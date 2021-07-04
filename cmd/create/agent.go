package create

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/agent"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
	"os"
)

// agentCmd represents the agent command
var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var configFilePath string
		configFilePath, err := cmd.Flags().GetString("config")
		utils.PrintError(err)

		if configFilePath == "" {
			configFilePath = types.DefaultFileName
		}

		obj, err := config.YamltoObject(configFilePath)
		utils.PrintError(err)

		fmt.Println(obj)

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
			Token:    token,
			Endpoint: obj.CurrentAccount,
		}

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		utils.PrintError(err)

		var newAgent types.Agent

		if nonInteractive {
			newAgent.ProjectId, err = cmd.Flags().GetString("project-id")
			utils.PrintError(err)

			newAgent.Mode, err = cmd.Flags().GetString("installation-mode")
			utils.PrintError(err)

			newAgent.AgentName, err = cmd.Flags().GetString("agent-name")
			utils.PrintError(err)

			newAgent.Description, err = cmd.Flags().GetString("agent-description")
			utils.PrintError(err)

			newAgent.PlatformName, err = cmd.Flags().GetString("platform-name")
			utils.PrintError(err)

			newAgent.ClusterType, err = cmd.Flags().GetString("cluster-type")
			utils.PrintError(err)

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

			// Check if user has sufficient permissions based on mode
			fmt.Println("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(newAgent.Mode, &kubeconfig)

			agent, err := apis.GetAgentList(credentials, newAgent.ProjectId)
			utils.PrintError(err)

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
			utils.PrintError(err)

			// Fetch project id
			projectID := agent.GetProjectID(userDetails)

			modeType := agent.GetModeType()

			// Check if user has sufficient permissions based on mode
			fmt.Println("\n🏃 Running prerequisites check....")
			agent.ValidateSAPermissions(modeType, &kubeconfig)
			newAgent, err = agent.GetAgentDetails(modeType, projectID, credentials, &kubeconfig)
			utils.PrintError(err)

			newAgent.ServiceAccount, newAgent.SAExists = k8s.ValidSA(newAgent.Namespace, &kubeconfig)
			newAgent.Mode = modeType
			agent.Summary(newAgent, &kubeconfig)

			agent.ConfirmInstallation()

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
	CreateCmd.AddCommand(agentCmd)

	agentCmd.Flags().BoolP("non-interactive", "n", false, "Help message for toggle")
	agentCmd.Flags().StringP("kubeconfig", "k", "", "Help message for toggle")
	agentCmd.Flags().String("project-id", "", "Help message for toggle")
	agentCmd.Flags().String("installation-mode", "", "Help message for toggle")
	agentCmd.Flags().String("agent-name", "", "Help message for toggle")
	agentCmd.Flags().String("agent-description", "", "Help message for toggle")
	agentCmd.Flags().String("platform-name", "", "Help message for toggle")
	agentCmd.Flags().String("cluster-type", "", "Help message for toggle")
	agentCmd.Flags().String("namespace", "", "Help message for toggle")
	agentCmd.Flags().String("service-account", "", "Help message for toggle")
	agentCmd.Flags().Bool("ns-exists", false, "Help message for toggle")
	agentCmd.Flags().Bool("sa-exists", false, "Help message for toggle")
}
