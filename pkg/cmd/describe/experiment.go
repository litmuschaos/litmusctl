package describe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// experimentCmd represents the Chaos Experiment command
var experimentCmd = &cobra.Command{
	Use:   "chaos-experiment",
	Short: "Describe a Chaos Experiment within the project",
	Long:  `Describe a Chaos Experiment within the project`,
	Run: func(cmd *cobra.Command, args []string) {
		credentials, err := utils.GetCredentials(cmd)
		utils.PrintError(err)

		var describeExperimentRequest model.ListExperimentRequest

		pid, err := cmd.Flags().GetString("project-id")
		utils.PrintError(err)

		if pid == "" {
			prompt := promptui.Prompt{
				Label: "Enter the Project ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.PrintError(err)
				os.Exit(1)
			}
			pid = result
		}

		var experimentID string
		if len(args) == 0 {
			prompt := promptui.Prompt{
				Label: "Enter the Chaos Experiment ID",
			}
			result, err := prompt.Run()
			if err != nil {
				utils.PrintError(err)
				os.Exit(1)
			}
			experimentID = result
		} else {
			experimentID = args[0]
		}

		// Handle blank input for Chaos Experiment ID
		if experimentID == "" {
			utils.Red.Println("⛔ Chaos Experiment ID can't be empty!!")
			os.Exit(1)
		}

		describeExperimentRequest.ExperimentIDs = append(describeExperimentRequest.ExperimentIDs, &experimentID)

		experiment, err := experiment.GetExperimentList(pid, describeExperimentRequest, credentials)
		if err != nil {
			if strings.Contains(err.Error(), "permission_denied") {
				utils.Red.Println("❌ The specified Project ID doesn't exist.")
				os.Exit(1)
			} else {
				utils.PrintError(err)
				os.Exit(1)
			}
		}

		if len(experiment.Data.ListExperimentDetails.Experiments) == 0 {
			utils.Red.Println("⛔ No chaos experiment found with ID: ", experimentID)
			os.Exit(1)
		}

		yamlManifest, err := yaml.JSONToYAML([]byte(experiment.Data.ListExperimentDetails.Experiments[0].ExperimentManifest))
		if err != nil {
			utils.Red.Println("❌ Error parsing Chaos Experiment manifest: " + err.Error())
			os.Exit(1)
		}

		// Add an output format prompt
		prompt := promptui.Select{
			Label: "Select an output format",
			Items: []string{"YAML", "JSON"},
		}
		i, _, err := prompt.Run()
		if err != nil {
			utils.PrintError(err)
			os.Exit(1)
		}

		switch i {
		case 0:
			// Output as YAML (default)
			utils.PrintInYamlFormat(string(yamlManifest))
		case 1:
			// Output as JSON
			jsonData, err := yaml.YAMLToJSON(yamlManifest)
			if err != nil {
				utils.Red.Println("❌ Error converting YAML to JSON: " + err.Error())
				os.Exit(1)
			}

			var prettyJSON bytes.Buffer
			err = json.Indent(&prettyJSON, jsonData, "", "    ") // Adjust the indentation as needed
			if err != nil {
				utils.Red.Println("❌ Error formatting JSON: " + err.Error())
				os.Exit(1)
			}

			fmt.Println(prettyJSON.String())
		default:
			utils.Red.Println("❌ Invalid output format selected")
			os.Exit(1)
		}
	},
}

func init() {
	DescribeCmd.AddCommand(experimentCmd)

	experimentCmd.Flags().String("project-id", "", "Set the project-id to list Chaos Experiments from the particular project. To see the projects, apply litmusctl get projects")
}
