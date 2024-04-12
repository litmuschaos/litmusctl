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

func DeleteProbeRequest(pid string, probeid string, cred types.Credentials) (DeleteProbeResponse, error) {
	var gqlReq DeleteProbeGQLRequest
	gqlReq.Query = DeleteProbeQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.ProbeName = probeid

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DeleteProbeResponse{}, errors.New("Error in deleting probe" + err.Error())
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
		return DeleteProbeResponse{}, errors.New("Error in deleting probe" + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DeleteProbeResponse{}, errors.New("Error in deleting probe" + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var deleteProbeResponse DeleteProbeResponse
		err = json.Unmarshal(bodyBytes, &deleteProbeResponse)
		if err != nil {
			return DeleteProbeResponse{}, errors.New("Error in deleting probe" + err.Error())
		}
		if len(deleteProbeResponse.Errors) > 0 {
			return DeleteProbeResponse{}, errors.New(deleteProbeResponse.Errors[0].Message)
		}
		return deleteProbeResponse, nil

	} else {
		return DeleteProbeResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}

}

func GetProbeYAMLRequest(pid string, request models.GetProbeYAMLRequest, cred types.Credentials) (GetProbeYAMLResponse, error) {
	var gqlReq GetProbeYAMLGQLRequest
	gqlReq.Query = GetProbeYAMLQuery
	gqlReq.Variables.ProjectID = pid
	gqlReq.Variables.Request = request

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return GetProbeYAMLResponse{}, errors.New("Error in getting probe details" + err.Error())
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
		return GetProbeYAMLResponse{}, errors.New("Error in getting probe details" + err.Error())
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return GetProbeYAMLResponse{}, errors.New("Error in getting probe details" + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var getProbeYAMLResponse GetProbeYAMLResponse
		err = json.Unmarshal(bodyBytes, &getProbeYAMLResponse)
		if err != nil {
			return GetProbeYAMLResponse{}, errors.New("Error in getting probes details" + err.Error())
		}
		if len(getProbeYAMLResponse.Errors) > 0 {
			return GetProbeYAMLResponse{}, errors.New(getProbeYAMLResponse.Errors[0].Message)
		}
		return getProbeYAMLResponse, nil

	} else {
		return GetProbeYAMLResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
