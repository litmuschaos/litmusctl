package experiment

import "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type SaveExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data SavedExperimentDetails `json:"data"`
}

type SavedExperimentDetails struct {
	Message string `json:"saveChaosExperiment"`
}

type SaveChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                  string                           `json:"projectID"`
		SaveChaosExperimentRequest model.SaveChaosExperimentRequest `json:"request"`
	} `json:"variables"`
}

type RunExperimentResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data RunExperimentData `json:"data"`
}

type RunExperimentData struct {
	RunExperimentDetails model.RunChaosExperimentResponse `json:"runChaosExperiment"`
}

type ExperimentListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentList `json:"data"`
}

type ExperimentList struct {
	ListExperimentDetails model.ListExperimentResponse `json:"listExperiment"`
}

type GetChaosExperimentsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		GetChaosExperimentRequest model.ListExperimentRequest `json:"request"`
		ProjectID                 string                      `json:"projectID"`
	} `json:"variables"`
}

type ExperimentRunData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentRun `json:"data"`
}

type ExperimentRun struct {
	ExperimentRunDetails model.ExperimentRun `json:"getExperimentRun"`
}

type ExperimentRunListData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ExperimentRunsList `json:"data"`
}

type ExperimentRunsList struct {
	ListExperimentRunDetails model.ListExperimentRunResponse `json:"listExperimentRun"`
}

type GetChaosExperimentRunsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                         `json:"projectID"`
		GetChaosExperimentRunRequest model.ListExperimentRunRequest `json:"request"`
	} `json:"variables"`
}

type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string `json:"projectID"`
		NotifyID  string `json:"notifyID"`
	} `json:"variables"`
}

type DeleteChaosExperimentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosExperimentDetails `json:"data"`
}

type DeleteChaosExperimentDetails struct {
	IsDeleted bool `json:"deleteChaosExperiment"`
}

type DeleteChaosExperimentGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID       string  `json:"projectID"`
		ExperimentID    *string `json:"experimentID"`
		ExperimentRunID *string `json:"experimentRunID"`
	} `json:"variables"`
}

type PodLogResponse struct {
	Log      string `json:"log"`
	Typename string `json:"__typename"`
}

type PodLogDetails struct {
	GetPodLog PodLogResponse `json:"getPodLog"`
}

type PodLogData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data PodLogDetails `json:"data"`
}

// Define the PodLogRequest structure
type PodLogRequest struct {
	InfraID         string `json:"infraID"`
	ExperimentRunID string `json:"experimentRunID"`
	PodName         string `json:"podName"`
	PodNamespace    string `json:"podNamespace"`
	PodType         string `json:"podType"`
	RunnerPod       string `json:"runnerPod"`
	ChaosNamespace  string `json:"chaosNamespace"`
}

type GetPodLogsGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		Request PodLogRequest `json:"request"`
	} `json:"variables"`
}
