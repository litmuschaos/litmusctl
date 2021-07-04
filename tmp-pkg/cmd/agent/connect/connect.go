package connect

import (
	"fmt"
	"os"

	utils "github.com/litmuschaos/litmusctl/tmp-pkg/common"
	"github.com/litmuschaos/litmusctl/tmp-pkg/common/chaos"
	"github.com/spf13/cobra"
)

// connectCmd represents the agent command
var ConnectCmd = &cobra.Command{
	Use:   "agent",
	Short: "Connect LitmusChaos agent",
	Long:  `Connect connects the agent to LitmusChaos`,
	Run: func(cmd *cobra.Command, args []string) {
		var c utils.Credentials
		var pErr error
		fmt.Println("🔥 Connecting LitmusChaos agent")
		fmt.Println("\n📶 Please enter LitmusChaos details --")
		// Get LitmusChaos URL as input
		c.Host, pErr = utils.GetPortalURL()
		if pErr != nil {
			fmt.Printf("\n❌ URL parsing failed: [%s]", pErr.Error())
			os.Exit(1)
		}
		// Get username as input
		c.Username = utils.GetUsername()
		// Get password as input
		c.Password = utils.GetPassword()
		// Fetch authorization token
		t := utils.Login(c, "auth/login")

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		if err != nil {
			fmt.Print(err)
		}

		chaos.Connect(t, c, kubeconfig)
	},
}
