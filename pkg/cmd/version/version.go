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
package version

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/litmuschaos/litmusctl/pkg/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of litmusctl",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cliVersion := os.Getenv("CLIVersion")
		if cliVersion == "" {
			utils.Red.Println("Error: CLIVersion environment variable is not set.")
			return
		}
		compatibilityArr := utils.CompatibilityMatrix[cliVersion]
		utils.White_B.Println("Litmusctl version: ", os.Getenv("CLIVersion"))
		utils.White_B.Println("Compatible ChaosCenter versions: ")
		utils.White_B.Print("[ ")
		for _, v := range compatibilityArr {
			utils.White_B.Print("'" + v + "' ")
		}
		utils.White_B.Print("]\n")
	},
}

var UpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Changes the version of litmusctl",
	Args:  cobra.ExactArgs(1),
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		updateVersion := args[0]
		homeDir, err := homedir.Dir()
		if err != nil {
			utils.PrintError(err)
		}

		var assetURL string = "https://litmusctl-production-bucket.s3.amazonaws.com/"
		var binaryName string = "litmusctl"
		switch runtime.GOOS {
		case "windows":
			if runtime.GOARCH == "386" {
				binaryName += "-windows-386-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			} else if runtime.GOARCH == "amd64" {
				binaryName += "-windows-amd64-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			} else {
				binaryName += "-windows-arm64-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			}
		case "linux":
			if runtime.GOARCH == "arm64" {
				binaryName += "-linux-arm64-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			} else if runtime.GOARCH == "amd64" {
				binaryName += "-linux-amd64-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			} else if runtime.GOARCH == "arm" {
				binaryName += "-linux-arm-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			} else {
				binaryName += "-linux-386-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			}
		case "darwin":
			if runtime.GOARCH == "amd64" {
				binaryName += "-darwin-amd64-" + updateVersion + ".tar.gz"
				assetURL += binaryName
			}
		}

		utils.White.Print("Downloading:\n")

		resp2, err := http.Get(assetURL)
		if err != nil {
			utils.PrintError(err)
		}

		tempFile, err := os.CreateTemp("", binaryName)
		if err != nil {
			utils.PrintError(err)
			return
		}
		defer os.Remove(tempFile.Name())

		_, err = io.Copy(tempFile, resp2.Body)
		if err != nil {
			utils.PrintError(err)
			return
		}
		utils.White_B.Print("OK\n")

		tempFile.Close()

		tarCmd := exec.Command("tar", "xzf", tempFile.Name(), "-C", homeDir)
		tarCmd.Stdout = os.Stdout
		tarCmd.Stderr = os.Stderr

		utils.White_B.Print("Extracting binary...\n")
		err = tarCmd.Run()
		if err != nil {
			utils.PrintError(err)
			return
		}

		utils.White_B.Print("Binary extracted successfully\n")
	},
}

func init() {
	VersionCmd.AddCommand(UpdateCmd)
}
