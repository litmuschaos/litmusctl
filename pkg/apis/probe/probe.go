package probe

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

func GetProbeRequest(pid string, probeID string, cred types.Credentials) (GetProbeResponse, error) {
	var gqlReq GetProbeGQLRequest
	gqlReq.Query = GetProbeQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.ProbeName = probeID

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return GetProbeResponse{}, errors.New("Error in getting requested probe" + err.Error())
	}

	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)

	if err != nil {
		return GetProbeResponse{}, errors.New("Error in getting requested probe" + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return GetProbeResponse{}, errors.New("Error in getting requested probe" + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var getProbeResponse GetProbeResponse
		err = json.Unmarshal(bodyBytes, &getProbeResponse)
		if err != nil {
			return GetProbeResponse{}, errors.New("Error in getting requested probe" + err.Error())
		}
		if len(getProbeResponse.Errors) > 0 {
			return GetProbeResponse{}, errors.New(getProbeResponse.Errors[0].Message)
		}
		return getProbeResponse, nil

	} else {
		return GetProbeResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}

}

func ListProbeRequest(pid string, probetypes []*models.ProbeType, cred types.Credentials) (ListProbeResponse, error) {
	var gqlReq ListProbeGQLRequest
	gqlReq.Query = ListProbeQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.Filter = models.ProbeFilterInput{
		Type: probetypes,
	}

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return ListProbeResponse{}, errors.New("Error in listing probes" + err.Error())
	}
	resp, err := apis.SendRequest(
		apis.SendRequestParams{
			Endpoint: cred.ServerEndpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)

	if err != nil {
		return ListProbeResponse{}, errors.New("Error in listing probes" + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ListProbeResponse{}, errors.New("Error in listing probes" + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var listProbeResponse ListProbeResponse
		err = json.Unmarshal(bodyBytes, &listProbeResponse)
		if err != nil {
			return ListProbeResponse{}, errors.New("Error in listing probes" + err.Error())
		}
		if len(listProbeResponse.Errors) > 0 {
			return ListProbeResponse{}, errors.New(listProbeResponse.Errors[0].Message)
		}
		return listProbeResponse, nil

	} else {
		return ListProbeResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}

}
