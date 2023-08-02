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

type GetChaosExperimentRunGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID                    string                         `json:"projectID"`
		GetChaosExperimentRunRequest model.ListExperimentRunRequest `json:"request"`
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

type ServerVersionResponse struct {
	Data   ServerVersionData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type ServerVersionData struct {
	GetServerVersion GetServerVersionData `json:"getServerVersion"`
}

type GetServerVersionData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
