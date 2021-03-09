package chaos

import (
	"fmt"
	"log"
	"net/url"

	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/go-resty/resty/v2"
	ymlparser "gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
)

type ListPkgData struct {
	Data struct {
		ListHubPkgData []PackageData `json:"ListHubPkgData"`
	} `json:"data"`
}
type PackageData struct {
	Experiments []string `json:"Experiments"`
	ChartName   string   `json:"chartName"`
}

type YAMLData struct {
	Data struct {
		GetYAMLData string `json:"getYAMLData"`
	} `json:"data"`
}

type GenerateWorkflowInputs struct {
	HubName        string
	ProjectID      string
	ChartName      string
	ExperimentName *string
	AccessToken    string
	FileType       *string
	URL            *url.URL
	WorkName       string
	WorkNamespace  string
	ClusterID      string
	Packages       []*PackageData
}

type GetClusters struct {
	Data struct {
		GetCluster []struct {
			ClusterID   string `json:"cluster_id"`
			ClusterName string `json:"cluster_name"`
		} `json:"getCluster"`
	} `json:"data"`
}

type GetHubStatus struct {
	Data struct {
		GetHubStatus []struct {
			ID      string `json:"id"`
			HubName string `json:"HubName"`
		} `json:"getHubStatus"`
	} `json:"data"`
}

func GetYamlData(inputs GenerateWorkflowInputs) (YAMLData, error) {
	client := resty.New()

	var yamlDataResponse YAMLData
	gql_query := `{"query":"query {\n  getYAMLData(experimentInput: {\n    ProjectID: \"` + inputs.ProjectID + `\"\n    HubName: \"` + inputs.HubName + `\"\n    ChartName: \"` + inputs.ChartName + `\"\n    ExperimentName: \"` + *inputs.ExperimentName + `\"\n    FileType: \"` + *inputs.FileType + `\"\n    \n  })\n}"}`
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", inputs.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(gql_query).
		SetResult(&yamlDataResponse).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				inputs.URL,
			),
		)

	if err != nil || !response.IsSuccess() {
		return YAMLData{}, err
	}

	return yamlDataResponse, nil
}

func GenerateWorkflow(wf_inputs GenerateWorkflowInputs) ([]byte, error) {

	var yaml v1alpha1.Workflow

	yaml.APIVersion = "argoproj.io/v1alpha1"
	yaml.Kind = "Workflow"
	yaml.ObjectMeta.Name = wf_inputs.WorkName
	yaml.ObjectMeta.Namespace = wf_inputs.WorkNamespace
	yaml.ObjectMeta.Labels = map[string]string{
		"cluster_id": wf_inputs.ClusterID,
	}

	var pram v1alpha1.Parameter
	pram.Name = "adminModeNamespace"
	pram.Value = &wf_inputs.WorkNamespace
	yaml.Spec.Arguments.Parameters = append(yaml.Spec.Arguments.Parameters, pram)
	//
	yaml.Spec.Entrypoint = "custom-chaos"
	var b = true
	var i int64 = 1000
	yaml.Spec.SecurityContext = &v1.PodSecurityContext{
		RunAsNonRoot: &b,
		RunAsUser:    &i,
	}

	var (
		custom_chaos        v1alpha1.Template
		install_experiments v1alpha1.Template
		engines             []v1alpha1.Template
		revert_chaos        v1alpha1.Template
	)

	custom_chaos.Name = "custom-chaos"
	custom_chaos.Steps = append(custom_chaos.Steps, v1alpha1.ParallelSteps{Steps: []v1alpha1.WorkflowStep{
		{
			Name:     "install-chaos-experiments",
			Template: "install-chaos-experiments",
		},
	}})

	install_experiments.Name = "install-chaos-experiments"
	install_experiments.Container = &v1.Container{
		Image:   "lachlanevenson/k8s-kubectl",
		Command: []string{"sh", "-c"},
		Args:    []string{""},
	}

	revert_chaos.Name = "revert-chaos"
	revert_chaos.Container = &v1.Container{
		Image:   "lachlanevenson/k8s-kubectl",
		Command: []string{"sh", "-c"},
		Args:    []string{"kubectl delete chaosengine "},
	}

	for _, pkg := range wf_inputs.Packages {

		for _, experiment := range pkg.Experiments {

			wf_inputs.ExperimentName = &experiment

			custom_chaos.Steps = append(custom_chaos.Steps, v1alpha1.ParallelSteps{Steps: []v1alpha1.WorkflowStep{
				{
					Name:     experiment,
					Template: experiment,
				},
			}})
			//
			var file_type = "experiment"
			wf_inputs.FileType = &file_type
			yamlData, err := GetYamlData(wf_inputs)
			if err != nil {
				log.Print(err)
			}

			install_experiments.Inputs.Artifacts = append(install_experiments.Inputs.Artifacts,
				v1alpha1.Artifact{
					Name: experiment,
					Path: "/tmp/" + experiment + ".yaml",
					ArtifactLocation: v1alpha1.ArtifactLocation{
						Raw: &v1alpha1.RawArtifact{
							Data: fmt.Sprint(yamlData.Data.GetYAMLData),
						},
					},
				})

			install_experiments.Container.Args[0] += "kubectl apply -f /tmp/" + experiment + ".yaml" + "-n {{workflow.parameters.adminModeNamespace}} | "

			revert_chaos.Container.Args[0] += experiment + " "

			file_type = "engine"
			wf_inputs.FileType = &file_type

			yamlData, err = GetYamlData(wf_inputs)
			if err != nil {
				log.Print(err)
			}

			var engine v1alpha1.Template
			engine.Name = experiment
			engine.Container = &v1.Container{
				Args: []string{
					`-file=/tmp/chaosengine-` + experiment + `.yaml`,
					"-saveName=/tmp/engine-name",
				},
				Image: "litmuschaos/litmus-checker:latest",
			}

			engine.Inputs.Artifacts = append(engine.Inputs.Artifacts, v1alpha1.Artifact{
				Name: experiment,
				Path: "/tmp/chaosengine-" + experiment + ".yaml",
				ArtifactLocation: v1alpha1.ArtifactLocation{
					Raw: &v1alpha1.RawArtifact{
						Data: fmt.Sprint(yamlData.Data.GetYAMLData),
					},
				},
			})

			engines = append(engines, engine)
		}
	}

	// Custom chaos
	custom_chaos.Steps = append(custom_chaos.Steps, v1alpha1.ParallelSteps{Steps: []v1alpha1.WorkflowStep{
		{
			Name:     "revert-chaos",
			Template: "revert-chaos",
		},
	}})

	yaml.Spec.Templates = append(yaml.Spec.Templates, custom_chaos)

	// Install experiments
	install_experiments.Container.Args[0] += "sleep 30"
	yaml.Spec.Templates = append(yaml.Spec.Templates, install_experiments)
	//
	// Install engines
	yaml.Spec.Templates = append(yaml.Spec.Templates, engines...)
	//
	// Revert Chaos
	revert_chaos.Container.Args[0] += "-n {{workflow.parameters.adminModeNamespace}}"
	yaml.Spec.Templates = append(yaml.Spec.Templates, revert_chaos)

	yamlByte, err := ymlparser.Marshal(yaml)
	if err != nil {
		return nil, err
	}

	return yamlByte, nil
}

func GetClustersQuery(project_id string, access_token string, url *url.URL) (GetClusters, error) {
	client := resty.New()

	var getClusters GetClusters
	gql_query := `{"query":"query {\n  getCluster(project_id: \"` + project_id + `\"){\n    cluster_id\n cluster_name\n  }\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", access_token)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(gql_query).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&getClusters).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				url,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return GetClusters{}, err
	}

	return getClusters, nil
}

func GetHubStatusQuery(project_id string, access_token string, url *url.URL) (GetHubStatus, error) {
	client := resty.New()

	var getHubStatus GetHubStatus
	gql_query := `{"query":"query {\n  getHubStatus(projectID: \"` + project_id + `\"){\n    id\n HubName \n  }\n}"}`
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", access_token)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(gql_query).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&getHubStatus).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				url,
			),
		)
	if err != nil || !response.IsSuccess() {
		return GetHubStatus{}, nil
	}

	return getHubStatus, nil
}

func ListPkgDataQuery(project_id string, hub_id string, access_token string, url *url.URL) (ListPkgData, error) {
	var pkgdata ListPkgData

	client := resty.New()

	gql_query := `{"query":"query {\n  ListHubPkgData(projectID: \"` + project_id + `\", hubID: \"` + hub_id + `\"){\n    Experiments\n    chartName\n  }\n}"}`
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", access_token)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(gql_query).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&pkgdata).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				url,
			),
		)
	if err != nil || !response.IsSuccess() {
		log.Print(err)
	}

	return pkgdata, nil
}
