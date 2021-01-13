package propel

import (
	"github.com/spf13/cobra"
)

// propelCmd represents the propel command
var PropelCmd = &cobra.Command{
	Use:   "propel",
	Short: "Propel is used to manage Kubera propel agents",
	Long:  `Kubera Propel agents can be managed via this command`,
}
