package completion

import (
	"strings"

	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/apis/experiment"
	"github.com/litmuschaos/litmusctl/pkg/apis/infrastructure"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/spf13/cobra"
)

func ProjectIDFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	credentials, err := utils.GetCredentials(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	projects, err := apis.ListProject(credentials)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	descriptions := make(map[string]string)

	for _, project := range projects.Data {
		if strings.HasPrefix(project.ID, toComplete) {
			completions = append(completions, project.ID)
			descriptions[project.ID] = project.Name
		}
	}

	var result []string
	for _, c := range completions {
		result = append(result, c+"\t"+descriptions[c])
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func ExperimentIDCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

	if len(args) > 0 {
		return nil, cobra.ShellCompDirectiveDefault
	}

	credentials, err := utils.GetCredentials(cmd)
	if err != nil {
		// Handle the error here if needed
		return nil, cobra.ShellCompDirectiveError
	}

	pid := cmd.Flag("project-id").Value.String()
	if pid == "" {
		return nil, cobra.ShellCompDirectiveError
	}

	var listExperimentRequest model.ListExperimentRequest
	listExperimentRequest.Filter = &model.ExperimentFilterInput{}

	experiments, err := experiment.GetExperimentList(pid, listExperimentRequest, credentials)
	if err != nil {
		// Handle the error here if needed
		return nil, cobra.ShellCompDirectiveError
	}

	completions := make([]string, 0)
	descriptions := make(map[string]string)

	for _, experiment := range experiments.Data.ListExperimentDetails.Experiments {
		if strings.HasPrefix(experiment.ExperimentID, toComplete) {
			completions = append(completions, experiment.ExperimentID)
			descriptions[experiment.ExperimentID] = experiment.Infra.Name + "/" + experiment.Name
		}
	}

	var result []string
	for _, c := range completions {
		result = append(result, c+"\t"+descriptions[c])
	}

	return result, cobra.ShellCompDirectiveNoFileComp
}

func ChaosInfraFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	credentials, err := utils.GetCredentials(cmd)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	pid := cmd.Flag("project-id").Value.String()
	if pid == "" {
		return nil, cobra.ShellCompDirectiveError
	}

	var listExperimentRequest model.ListExperimentRequest
	listExperimentRequest.Filter = &model.ExperimentFilterInput{}

	infras, err := infrastructure.GetInfraList(credentials, pid, model.ListInfraRequest{})
	if err != nil {
		if strings.Contains(err.Error(), "permission_denied") {
			return nil, cobra.ShellCompDirectiveError

		} else {
			return nil, cobra.ShellCompDirectiveError

		}
	}

	completions := make([]string, 0)
	descriptions := make(map[string]string)

	for _, infra := range infras.Data.ListInfraDetails.Infras {
		if strings.HasPrefix(infra.InfraID, toComplete) {
			completions = append(completions, infra.InfraID)
			descriptions[infra.InfraID] = infra.Name
		}
	}

	var result []string
	for _, c := range completions {
		result = append(result, c+"\t"+descriptions[c])
	}
	return result, cobra.ShellCompDirectiveNoFileComp
}

func OutputFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"json", "yaml"}, cobra.ShellCompDirectiveNoFileComp
}

func InstallModeTypeFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"cluster", "namespace"}, cobra.ShellCompDirectiveNoFileComp
}

func ChaosInfraTypeFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"external", "internal"}, cobra.ShellCompDirectiveNoFileComp
}

func PlatformNameFlagCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"AWS", "GKE", "Openshift", "Rancher", "Others"}, cobra.ShellCompDirectiveNoFileComp
}
