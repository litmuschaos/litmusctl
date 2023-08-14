package run

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var RunCmd = &cobra.Command{
	Use: "run",
	Short: `Runs Experiment for LitmusChaos Execution plane.
		Examples:

		#Run a Chaos Experiment
		litmusctl run chaos-experiment  --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --experiment-id="ab754058dd04"

		Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
}
