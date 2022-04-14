/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package create

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var CreateCmd = &cobra.Command{
	Use: "create",
	Short: `Create resources for LitmusChaos agent plane.
		Examples:
		#create a project
		litmusctl create project --name new-proj

		#create an agent
		litmusctl create agent --agent-name="new-agent" --non-interactive

		#create an agent within a project
		litmusctl create agent --agent-name="new-agent" --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --non-interactive

		#create a workflow from a file
		litmusctl create workflow -f workflow.yaml --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --cluster-id="d861b650-1549-4574-b2ba-ab754058dd04"

		Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
}
