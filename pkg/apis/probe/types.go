package probe

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type ListProbeGQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		ProjectID          string                    `json:"projectID"`
		InfrastructureType *model.InfrastructureType `json:"infrastructureType"`
		ProbeNames         []string                  `json:"probeNames"`
		Filter             model.ProbeFilterInput    `json:"probeFilterInput"`
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
