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
package rootCmd

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/cmd/upgrade"
	"github.com/litmuschaos/litmusctl/pkg/cmd/version"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/litmuschaos/litmusctl/pkg/cmd/config"
	"github.com/litmuschaos/litmusctl/pkg/cmd/create"
	"github.com/litmuschaos/litmusctl/pkg/cmd/get"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string
var SkipSSLVerify bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "litmusctl",
	Short: "Litmusctl controls the litmuschaos agent plane",
	Long:  `Litmusctl controls the litmuschaos agent plane. ` + "\n" + ` Find more information at: https://github.com/litmuschaos/litmusctl`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(config.ConfigCmd)
	rootCmd.AddCommand(create.CreateCmd)
	rootCmd.AddCommand(get.GetCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(upgrade.UpgradeCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.litmusctl)")
	rootCmd.PersistentFlags().BoolVar(&SkipSSLVerify, "skipSSL", false, "skipSSL")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".litmusconfig" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(utils.DefaultFileName)
	}

	if SkipSSLVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
