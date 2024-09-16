package update

import (
	"github.com/spf13/cobra"
)


var UpdateCmd = &cobra.Command{
	Use: "update",
	Short: `It updates ChaosCenter account's details. 
		Examples
		
		#update password
		litmusctl update password
		`,
}