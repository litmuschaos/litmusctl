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
package config

import (
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"io/ioutil"

	"os"

	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/spf13/cobra"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "Display litmusconfig settings or a specified litmusconfig file",
	Long: `Display litmusconfig settings or a specified litmusconfig file. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var configFilePath string
		configFilePath, err := cmd.Flags().GetString("config")
		utils.PrintError(err)

		if configFilePath == "" {
			configFilePath = utils.DefaultFileName
		}

		exists := config.FileExists(configFilePath)
		if !exists {
			fmt.Println("File reading error open ", configFilePath, ": no such file or directory. Use --config or -c flag to point the configfile")
			os.Exit(1)
		}

		data, err := ioutil.ReadFile(configFilePath)
		utils.PrintError(err)

		//Printing the config map
		fmt.Print(string(data))
	},
}

func init() {
	ConfigCmd.AddCommand(viewCmd)
}
