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
package delete

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use: "delete",
	Short: `Delete resources for LitmusChaos Execution plane.
		Examples:
		#delete a Chaos Experiment
		litmusctl delete chaos-experiment c520650e-7cb6-474c-b0f0-4df07b2b025b --project-id=c520650e-7cb6-474c-b0f0-4df07b2b025b

		#delete a Chaos Environment
		litmusctl delete chaos-environment --project-id=8adf62d5-64f8-4c66-ab53-63729db9dd9a --environment-id=environmentexample

		Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
	`,
}
