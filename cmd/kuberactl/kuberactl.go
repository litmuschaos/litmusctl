package kuberactl

import (
	"fmt"
	"os"

	"github.com/mayadata-io/kuberactl/cmd/chaos"
	chaosAgent "github.com/mayadata-io/kuberactl/cmd/chaos/agent"
	chaosRegister "github.com/mayadata-io/kuberactl/cmd/chaos/agent/register"
	"github.com/mayadata-io/kuberactl/cmd/propel"
	propelAgent "github.com/mayadata-io/kuberactl/cmd/propel/agent"
	propelRegister "github.com/mayadata-io/kuberactl/cmd/propel/agent/register"
	"github.com/mayadata-io/kuberactl/cmd/version"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kuberactl",
	Short: "A brief description of your application",
	Long:  `Kuberactl is a cli for managing Kubera resources.`,
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
	RootCmd.AddCommand(chaos.ChaosCmd)
	RootCmd.AddCommand(propel.PropelCmd)
	chaos.ChaosCmd.AddCommand(chaosAgent.AgentCmd)
	propel.PropelCmd.AddCommand(propelAgent.AgentCmd)
	chaosAgent.AgentCmd.AddCommand(chaosRegister.RegisterCmd)
	propelAgent.AgentCmd.AddCommand(propelRegister.RegisterCmd)
	RootCmd.AddCommand(version.VersionCmd)

	// Create a persistent flag for kubeconfig
	// RootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "absolute path to the kubeconfig file")
}
