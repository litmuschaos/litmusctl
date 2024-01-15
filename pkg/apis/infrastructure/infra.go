/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a1 copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/apis"
	"github.com/litmuschaos/litmusctl/pkg/types"

	models "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

// GetInfraList lists the Chaos Infrastructure connected to the specified project
func GetInfraList(c types.Credentials, pid string, request models.ListInfraRequest) (InfraData, error) {
	var gplReq ListInfraGraphQLRequest
	gplReq.Query = ListInfraQuery
	gplReq.Variables.ProjectID = pid
	gplReq.Variables.ListInfraRequest = request

	query, err := json.Marshal(gplReq)
	if err != nil {
		return InfraData{}, err
	}
	resp, err := apis.SendRequest(apis.SendRequestParams{Endpoint: c.ServerEndpoint + utils.GQLAPIPath, Token: c.Token}, query, string(types.Post))
	if err != nil {
		return InfraData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {

		return InfraData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var Infra InfraData
		err = json.Unmarshal(bodyBytes, &Infra)
		if err != nil {
			return InfraData{}, err
		}

		if len(Infra.Errors) > 0 {
			return InfraData{}, errors.New(Infra.Errors[0].Message)
		}

		return Infra, nil
	} else {
		return InfraData{}, fmt.Errorf("error getting detais from server")
	}
}

// ConnectInfra connects the  Infra with the given details
func ConnectInfra(infra types.Infra, cred types.Credentials) (InfraConnectionData, error) {
	var gqlReq RegisterInfraGqlRequest
	gqlReq.Query = RegisterInfraQuery
	gqlReq.Variables.ProjectId = infra.ProjectId
	gqlReq.Variables.RegisterInfraRequest = CreateRegisterInfraRequest(infra)

	if infra.NodeSelector != "" {
		gqlReq.Variables.RegisterInfraRequest.NodeSelector = &infra.NodeSelector
	}

	if infra.Tolerations != "" {
		var toleration []*models.Toleration
		err := json.Unmarshal([]byte(infra.Tolerations), &toleration)
		utils.PrintError(err)
		gqlReq.Variables.RegisterInfraRequest.Tolerations = toleration
	}

	if infra.NodeSelector != "" && infra.Tolerations != "" {
		gqlReq.Variables.RegisterInfraRequest.NodeSelector = &infra.NodeSelector

		var toleration []*models.Toleration
		err := json.Unmarshal([]byte(infra.Tolerations), &toleration)
		utils.PrintError(err)
		gqlReq.Variables.RegisterInfraRequest.Tolerations = toleration
	}

	query, err := json.Marshal(gqlReq)
	resp, err := apis.SendRequest(apis.SendRequestParams{Endpoint: cred.ServerEndpoint + utils.GQLAPIPath, Token: cred.Token}, query, string(types.Post))
	if err != nil {
		return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
	}

	if resp.StatusCode == http.StatusOK {
		var connectInfra InfraConnectionData
		err = json.Unmarshal(bodyBytes, &connectInfra)
		if err != nil {
			return InfraConnectionData{}, errors.New("Error in registering Chaos Infrastructure: " + err.Error())
		}

		if len(connectInfra.Errors) > 0 {
			return InfraConnectionData{}, errors.New(connectInfra.Errors[0].Message)
		}
		return connectInfra, nil
	} else {
		return InfraConnectionData{}, err
	}
}

func CreateRegisterInfraRequest(infra types.Infra) (request models.RegisterInfraRequest) {
	return models.RegisterInfraRequest{
		Name:               infra.InfraName,
		InfraScope:         infra.Mode,
		Description:        &infra.Description,
		PlatformName:       infra.PlatformName,
		EnvironmentID:      infra.EnvironmentID,
		InfrastructureType: models.InfrastructureTypeKubernetes,
		InfraNamespace:     &infra.Namespace,
		ServiceAccount:     &infra.ServiceAccount,
		InfraNsExists:      &infra.NsExists,
		InfraSaExists:      &infra.SAExists,
		SkipSsl:            &infra.SkipSSL,
	}
}

// DisconnectInfra sends GraphQL API request for disconnecting Chaos Infra(s).
func DisconnectInfra(projectID string, infraID string, cred types.Credentials) (DisconnectInfraData, error) {

	var gqlReq DisconnectInfraGraphQLRequest
	var err error

	gqlReq.Query = DisconnectInfraQuery
	gqlReq.Variables.ProjectID = projectID
	gqlReq.Variables.InfraID = infraID

	query, err := json.Marshal(gqlReq)
	if err != nil {
		return DisconnectInfraData{}, err
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
		return DisconnectInfraData{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return DisconnectInfraData{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var disconnectInfraData DisconnectInfraData
		err = json.Unmarshal(bodyBytes, &disconnectInfraData)
		if err != nil {
			return DisconnectInfraData{}, err
		}

		if len(disconnectInfraData.Errors) > 0 {
			return DisconnectInfraData{}, errors.New(disconnectInfraData.Errors[0].Message)
		}

		return disconnectInfraData, nil
	} else {
		return DisconnectInfraData{}, err
	}
}
