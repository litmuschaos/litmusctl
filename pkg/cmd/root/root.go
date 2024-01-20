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
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/cmd/run"
	"github.com/litmuschaos/litmusctl/pkg/cmd/save"

	"github.com/litmuschaos/litmusctl/pkg/cmd/connect"
	"github.com/litmuschaos/litmusctl/pkg/cmd/delete"
	"github.com/litmuschaos/litmusctl/pkg/cmd/describe"
	"github.com/litmuschaos/litmusctl/pkg/cmd/disconnect"
	"github.com/litmuschaos/litmusctl/pkg/cmd/upgrade"
	"github.com/litmuschaos/litmusctl/pkg/cmd/version"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/litmuschaos/litmusctl/pkg/cmd/config"
	"github.com/litmuschaos/litmusctl/pkg/cmd/create"
	"github.com/litmuschaos/litmusctl/pkg/cmd/get"
	"github.com/litmuschaos/litmusctl/pkg/cmd/list"
	config2 "github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/spf13/cobra"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

//var kubeconfig string

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
	rootCmd.AddCommand(connect.ConnectCmd)
	rootCmd.AddCommand(disconnect.DisconnectCmd)
	rootCmd.AddCommand(delete.DeleteCmd)
	rootCmd.AddCommand(describe.DescribeCmd)
	rootCmd.AddCommand(version.VersionCmd)
	rootCmd.AddCommand(upgrade.UpgradeCmd)
	rootCmd.AddCommand(save.SaveCmd)
	rootCmd.AddCommand(run.RunCmd)
	rootCmd.AddCommand(list.ListCmd)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.litmusctl)")
	//rootCmd.PersistentFlags().StringVar(&kubeconfig, "kubeconfig", "", "kubeconfig file (default is $HOME/.kube/config")
	rootCmd.PersistentFlags().BoolVar(&config2.SkipSSLVerify, "skipSSL", false, "skipSSL, litmusctl will skip ssl/tls verification while communicating with portal")
	rootCmd.PersistentFlags().StringVar(&config2.CACert, "cacert", "", "cacert <path_to_crt_file> , custom ca certificate used for communicating with portal")
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

	viper.AutomaticEnv() // read in environment variables that match

	if config2.SkipSSLVerify {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	} else if config2.CACert != "" {
		caCert, err := ioutil.ReadFile(config2.CACert)
		cobra.CheckErr(err)
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{RootCAs: caCertPool}
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
