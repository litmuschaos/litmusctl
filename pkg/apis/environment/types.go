package environment

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type CreateEnvironmentGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectId string                         `json:"projectID"`
		Request   model.CreateEnvironmentRequest `json:"request"`
	} `json:"variables"`
}

type CreateEnvironmentResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data CreateEnvironmentData `json:"data"`
}

type CreateEnvironmentData struct {
	EnvironmentDetails model.Environment `json:"createEnvironment"`
}

type ListEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data EnvironmentsList `json:"data"`
}

type EnvironmentsList struct {
	ListEnvironmentDetails model.ListEnvironmentResponse `json:"listEnvironments"`
}

type CreateEnvironmentListGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string                       `json:"projectID"`
		Request   model.ListEnvironmentRequest `json:"request"`
	}
}

type GetEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetEnvironment `json:"data"`
}

type GetEnvironment struct {
	EnvironmentDetails model.Environment `json:"getEnvironment"`
}

type CreateEnvironmentGetGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string       `json:"projectID"`
		EnvironmentID string   `json:"environmentID"`
	}
}

type CreateEnvironmentDeleteGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID     string `json:"projectID"`
		EnvironmentID string `json:"environmentID"`
	}
}

type DeleteChaosEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteChaosEnvironmentDetails `json:"data"`
}

type DeleteChaosEnvironmentDetails struct {
	DeleteChaosEnvironment string `json:"deleteChaosExperiment"`
}
