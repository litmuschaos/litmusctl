package update

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

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
