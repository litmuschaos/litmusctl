package register

import (
	"fmt"
	"os"

	utils "github.com/litmuschaos/litmusctl/pkg/common"
	chaos "github.com/litmuschaos/litmusctl/pkg/common/chaos"
	"github.com/spf13/cobra"
)

// registerCmd represents the register command
var RegisterCmd = &cobra.Command{
	Use:   "register",
	Short: "Register LitmusChaos agent",
	Long:  `Register registers the agent to LitmusChaos`,
	Run: func(cmd *cobra.Command, args []string) {

		var c utils.Credentials
		var pErr error
		fmt.Println("🔥 Registering LitmusChaos agent")
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

		chaos.Register(t, c)
	},
}
