package agent

import (
	"github.com/spf13/cobra"
)

// agentCmd represents the agent command
var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Kubera Agent",
	Long:  `agent is used to manage Kubera agents`,
}
