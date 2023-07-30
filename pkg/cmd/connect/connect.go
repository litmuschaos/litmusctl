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
package connect

import (
	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var ConnectCmd = &cobra.Command{
	Use: "connect",
	Short: `Connect resources for LitmusChaos Execution plane.
		Examples:
		#connect a Chaos Infrastructure
		litmusctl connect chaos-infra --name="new-chaos-infra" --non-interactive

		#connect a chaos-infrastructure within a project
		litmusctl connect chaos-infra --name="new-chaos-infra" --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --non-interactive
		Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
}
