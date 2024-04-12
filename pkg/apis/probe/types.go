package probe

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type GetProbeGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string `json:"projectID"`
		ProbeName string `json:"probeName"`
	} `json:"variables"`
}

type GetProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetProbeResponseData `json:"data"`
}

type GetProbeResponseData struct {
	GetProbe model.Probe `json:"getProbe"`
}

type ListProbeGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID string                 `json:"projectID"`
		Filter    model.ProbeFilterInput `json:"filter"`
	} `json:"variables"`
}

type ListProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ListProbeResponseData `json:"data"`
}

type ListProbeResponseData struct {
	Probes []model.Probe `json:"listProbes"`
}
type DeleteProbeGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProbeName string `json:"probeName"`
		ProjectID string `json:"projectID"`
	} `json:"variables"`
}

type DeleteProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteProbeResponseData `json:"data"`
}

type DeleteProbeResponseData struct {
	DeleteProbe bool `json:"deleteProbe"`
}
