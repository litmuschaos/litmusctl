package environment

import models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type CreateEnvironmentGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectId string                          `json:"projectID"`
		Request   models.CreateEnvironmentRequest `json:"request"`
	} `json:"variables"`
}

type CreateEnvironmentResponse struct {
	Data   CreateEnvironmentData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type CreateEnvironmentData struct {
	EnvironmentDetails models.Environment `json:"createEnvironment"`
}

type ListEnvironmentData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data EnvironmentsList `json:"data"`
}

type EnvironmentsList struct {
	ListEnvironmentDetails models.ListEnvironmentResponse `json:"listEnvironments"`
}

type CreateEnvironmentListGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string                        `json:"projectID"`
		Request   models.ListEnvironmentRequest `json:"request"`
	}
}
