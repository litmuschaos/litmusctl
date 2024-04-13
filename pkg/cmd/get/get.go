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
package get

import (
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var GetCmd = &cobra.Command{
	Use: "get",
	Short: `Examples:
		#get list of projects accessed by the user
		litmusctl get projects

		#get list of Chaos Infrastructure within the project
		litmusctl get chaos-infra --project-id=""

		#get list of chaos Chaos Experiments
		litmusctl get chaos-experiments --project-id=""

		#get list of Chaos Experiment runs
		litmusctl get chaos-experiment-runs --project-id=""

		#get list of Chaos Environments
		litmusctl get chaos-environments --project-id=""

		#get list of Probes in a Project
		litmusctl get probes --project-id=""

		Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
}
