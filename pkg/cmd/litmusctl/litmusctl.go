package litmusctl

import (
	"fmt"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/cmd/version"

	chaosAgent "github.com/litmuschaos/litmusctl/pkg/cmd/agent"
	chaosConnect "github.com/litmuschaos/litmusctl/pkg/cmd/agent/connect"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "litmusctl",
	Short: "A brief description of your application",
	Long:  `Litmusctl is a cli for managing Litmus resources.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	RootCmd.AddCommand(version.VersionCmd)
	RootCmd.AddCommand(chaosAgent.AgentCmd)
	chaosAgent.AgentCmd.AddCommand(chaosConnect.ConnectCmd)

	// Create a persistent flag for kubeconfig
	chaosConnect.ConnectCmd.PersistentFlags().StringP("kubeconfig", "k", "", "absolute path to the kubeconfig file")
}
