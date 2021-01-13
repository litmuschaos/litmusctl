package chaos

import (
	"github.com/spf13/cobra"
)

// litmusCmd represents the litmus command
var ChaosCmd = &cobra.Command{
	Use:   "chaos",
	Short: "Kubera Chaos",
	Long:  `Kubera Chaos is used to run chaos tests`,
}
