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
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/apis/environment"
	"github.com/litmuschaos/litmusctl/pkg/apis/infrastructure"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/infra_ops"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/spf13/cobra"
)

// infraCmd represents the Chaos infra command
var infraCmd = &cobra.Command{
	Use: "chaos-infra",
	Short: `Connect an external Chaos infra.
	Example(s):
	#connect a Chaos infra
	litmusctl connect chaos-infra --name="new-chaos-infra" --non-interactive

	#connect a Chaos infra within a project
	litmusctl connect chaos-infra --name="new-chaos-infra" --project-id="d861b650-1549-4574-b2ba-ab754058dd04" --non-interactive

	Note: The default location of the config file is $HOME/.litmusconfig, and can be overridden by a --config flag
`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		nonInteractive, err := cmd.Flags().GetBool("non-interactive")
		utils.PrintError(err)

		kubeconfig, err := cmd.Flags().GetString("kubeconfig")
		utils.PrintError(err)

		var newInfra types.Infra

		newInfra.ProjectId, err = cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if newInfra.ProjectId == "" {
			userDetails, err := apis.GetProjectDetails(credentials)
			utils.PrintError(err)

			var (
				userID        = userDetails.Data.ID
				projectExists = false
			)

		outerloop:
			for _, project := range userDetails.Data.Projects {
				for _, member := range project.Members {
					if (member.UserID == userID) && (member.Role == "Owner" || member.Role == "Editor") {
						projectExists = true
						break outerloop
					}
				}
			}

			if !projectExists {
				utils.White_B.Print("Creating a random project...")
				newInfra.ProjectId = infra_ops.CreateRandomProject(credentials)
			}
		}

		if nonInteractive {

			newInfra.Mode, err = cmd.Flags().GetString("installation-mode")
			utils.PrintError(err)

			if newInfra.Mode == "" {
				utils.Red.Print("Error: --installation-mode flag is empty")
				os.Exit(1)
			}

			newInfra.InfraName, err = cmd.Flags().GetString("name")
			utils.PrintError(err)

			newInfra.SkipSSL, err = cmd.Flags().GetBool("skip-ssl")
			utils.PrintError(err)

			if newInfra.InfraName == "" {
				utils.Red.Print("Error: --name flag is empty")
				os.Exit(1)
			}

			newInfra.EnvironmentID, err = cmd.Flags().GetString("environmentID")
			if newInfra.EnvironmentID == "" {
				utils.Red.Print("Error: --environment flag is empty")
				os.Exit(1)
			}

			newInfra.Description, err = cmd.Flags().GetString("description")
			utils.PrintError(err)

			newInfra.PlatformName, err = cmd.Flags().GetString("platform-name")
			utils.PrintError(err)

			if newInfra.PlatformName == "" {
				utils.Red.Print("Error: --platform-name flag is empty")
				os.Exit(1)
			}

			newInfra.InfraType, err = cmd.Flags().GetString("chaos-infra-type")
			utils.PrintError(err)
			if newInfra.InfraType == "" {
				utils.Red.Print("Error: --chaos-infra-type flag is empty")
				os.Exit(1)
			}

			newInfra.NodeSelector, err = cmd.Flags().GetString("node-selector")
			utils.PrintError(err)
			if newInfra.NodeSelector != "" {
				if ok := utils.CheckKeyValueFormat(newInfra.NodeSelector); !ok {
					os.Exit(1)
				}
			}

			toleration, err := cmd.Flags().GetString("tolerations")
			utils.PrintError(err)

			if toleration != "" {
				newInfra.Tolerations = toleration
			}

			newInfra.Namespace, err = cmd.Flags().GetString("namespace")
			utils.PrintError(err)

			newInfra.ServiceAccount, err = cmd.Flags().GetString("service-account")
			utils.PrintError(err)

			newInfra.NsExists, err = cmd.Flags().GetBool("ns-exists")
			utils.PrintError(err)

			newInfra.SAExists, err = cmd.Flags().GetBool("sa-exists")
			utils.PrintError(err)

			if newInfra.Mode == "" {
				newInfra.Mode = utils.DefaultMode
			}

			if newInfra.ProjectId == "" {
				utils.Red.Println("Error: --project-id flag is empty")
				os.Exit(1)
			}

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			infra_ops.ValidateSAPermissions(newInfra.Namespace, newInfra.Mode, &kubeconfig)

			// Check if infra already exists
			isInfraExist, err, infraList := infra_ops.ValidateInfraNameExists(newInfra.InfraName, newInfra.ProjectId, credentials)
			utils.PrintError(err)

			if isInfraExist {
				infra_ops.PrintExistingInfra(infraList)
				os.Exit(1)
			}
			envIDs, err := environment.GetEnvironmentList(newInfra.ProjectId, credentials)
			utils.PrintError(err)

			// Check if Environment exists
			var isEnvExist = false
			for i := range envIDs.Data.ListEnvironmentDetails.Environments {
				if newInfra.EnvironmentID == envIDs.Data.ListEnvironmentDetails.Environments[i].EnvironmentID {
					utils.White_B.Print(envIDs.Data.ListEnvironmentDetails.Environments[i].EnvironmentID)
					isEnvExist = true
					break
				}
			}
			if !isEnvExist {
				utils.Red.Println("\nChaos Environment with the given ID doesn't exists.")
				infra_ops.PrintExistingEnvironments(envIDs)
				utils.White_B.Println("\n❗ Please enter a name from the List or Create a new environment using `litmusctl create chaos-environment`")
				os.Exit(1)
			}

		} else {
			userDetails, err := apis.GetProjectDetails(credentials)
			utils.PrintError(err)

			if newInfra.ProjectId == "" {
				// Fetch project id
				newInfra.ProjectId = infra_ops.GetProjectID(userDetails)
			}

			modeType := infra_ops.GetModeType()

			// Check if user has sufficient permissions based on mode
			utils.White_B.Print("\n🏃 Running prerequisites check....")
			infra_ops.ValidateSAPermissions(newInfra.Namespace, modeType, &kubeconfig)
			newInfra, err = infra_ops.GetInfraDetails(modeType, newInfra.ProjectId, credentials, &kubeconfig)
			utils.PrintError(err)

			newInfra.ServiceAccount, newInfra.SAExists = k8s.ValidSA(newInfra.Namespace, &kubeconfig)
			newInfra.Mode = modeType
		}

		infra_ops.Summary(newInfra, &kubeconfig)

		if !nonInteractive {
			infra_ops.ConfirmInstallation()
		}

		infra, err := infrastructure.ConnectInfra(newInfra, credentials)
		if err != nil {
			utils.Red.Println("\n❌ Chaos Infra connection failed: " + err.Error() + "\n")
			os.Exit(1)
		}

		if infra.Data.RegisterInfraDetails.Token == "" {
			utils.Red.Println("\n❌ failed to get the Infra registration token: " + "\n")
			os.Exit(1)
		}

		path := fmt.Sprintf("%s%s/%s.yaml", credentials.Endpoint, utils.ChaosYamlPath, infra.Data.RegisterInfraDetails.Token)
		utils.White_B.Print("Applying YAML:\n", path)

		// Print error message in case Data field is null in response
		if (infra.Data == infrastructure.RegisterInfra{}) {
			utils.White_B.Print("\n🚫 Chaos new infrastructure connection failed: " + infra.Errors[0].Message + "\n")
			os.Exit(1)
		}

		//Apply infra connection yaml
		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    infra.Data.RegisterInfraDetails.Token,
			Endpoint: credentials.Endpoint,
			YamlPath: utils.ChaosYamlPath,
		}, kubeconfig, false)
		if err != nil {
			utils.Red.Print("\n❌ Failed to apply connection yaml: \n" + err.Error() + "\n")
			utils.White_B.Print("\n Error:  \n" + err.Error())
			os.Exit(1)
		}

		utils.White_B.Print("\n", yamlOutput)

		// Watch subscriber pod status
		k8s.WatchPod(k8s.WatchPodParams{Namespace: newInfra.Namespace, Label: utils.ChaosInfraLabel}, &kubeconfig)

		utils.White_B.Println("\n🚀 Chaos new infrastructure connection successful!! 🎉")
		utils.White_B.Println("👉 Litmus Chaos Infrastructure can be accessed here: " + fmt.Sprintf("%s/%s", credentials.Endpoint, utils.ChaosInfraPath))
	},
}

func init() {
	ConnectCmd.AddCommand(infraCmd)

	infraCmd.Flags().BoolP("non-interactive", "n", false, "Set it to true for non interactive mode | Note: Always set the boolean flag as --non-interactive=Boolean")
	infraCmd.Flags().StringP("kubeconfig", "k", "", "Set to pass kubeconfig file if it is not in the default location ($HOME/.kube/config)")
	infraCmd.Flags().String("tolerations", "", "Set the tolerations for Chaos infra components | Format: '[{\"key\":\"key1\",\"value\":\"value1\",\"operator\":\"Exist\",\"effect\":\"NoSchedule\",\"tolerationSeconds\":30}]'")

	infraCmd.Flags().String("project-id", "", "Set the project-id to install Chaos infra for the particular project. To see the projects, apply litmusctl get projects")
	infraCmd.Flags().String("installation-mode", "cluster", "Set the installation mode for the kind of Chaos infra | Supported=cluster/namespace")
	infraCmd.Flags().String("name", "", "Set the Chaos infra name")
	infraCmd.Flags().String("description", "---", "Set the Chaos infra description")
	infraCmd.Flags().String("platform-name", "Others", "Set the platform name. Supported- AWS/GKE/Openshift/Rancher/Others")
	infraCmd.Flags().String("chaos-infra-type", "external", "Set the chaos-infra-type to external for external Chaos infras | Supported=external/internal")
	infraCmd.Flags().String("node-selector", "", "Set the node-selector for Chaos infra components | Format: \"key1=value1,key2=value2\")")
	infraCmd.Flags().String("namespace", "litmus", "Set the namespace for the Chaos infra installation")
	infraCmd.Flags().String("service-account", "litmus", "Set the service account to be used by the Chaos infra")
	infraCmd.Flags().Bool("skip-ssl", false, "Set whether Chaos infra will skip ssl/tls check (can be used for self-signed certs, if cert is not provided in portal)")
	infraCmd.Flags().Bool("ns-exists", false, "Set the --ns-exists=false if the namespace mentioned in the --namespace flag is not existed else set it to --ns-exists=true | Note: Always set the boolean flag as --ns-exists=Boolean")
	infraCmd.Flags().Bool("sa-exists", false, "Set the --sa-exists=false if the service-account mentioned in the --service-account flag is not existed else set it to --sa-exists=true | Note: Always set the boolean flag as --sa-exists=Boolean\"\n")
}
