package infrastructure

import models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type InfraData struct {
	Data   InfraList `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type InfraList struct {
	ListInfraDetails models.ListInfraResponse `json:"listInfras"`
}

type ListInfraGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID        string                  `json:"projectID"`
		ListInfraRequest models.ListInfraRequest `json:"request"`
	} `json:"variables"`
}

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

type RegisterInfraGqlRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectId            string                      `json:"projectID"`
		RegisterInfraRequest models.RegisterInfraRequest `json:"request"`
	} `json:"variables"`
}

type InfraConnectionData struct {
	Data   RegisterInfra `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type RegisterInfra struct {
	RegisterInfraDetails models.RegisterInfraResponse `json:"registerInfra"`
}

type DisconnectInfraData struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DisconnectInfraDetails `json:"data"`
}

type DisconnectInfraDetails struct {
	Message string `json:"deleteInfra"`
}

type DisconnectInfraGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string `json:"projectID"`
		InfraID   string `json:"infraID"`
	} `json:"variables"`
}
