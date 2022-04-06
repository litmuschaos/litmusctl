/*
Copyright © 2021 The LitmusChaos Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package apis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	types "github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
)

// CreateWorkflow sends GraphQL API request for creating a workflow
func CreateWorkflow(in types.ChaosWorkFlowInput, cred types.Credentials) {

	var gqlReq types.ChaosWorkFlowGraphQLRequest

	gqlReq.Query = `mutation createChaosWorkFlow($ChaosWorkFlowInput: ChaosWorkFlowInput!) {
                      createChaosWorkFlow(input: $ChaosWorkFlowInput) {
                        workflow_id
                        cronSyntax
                        workflow_name
                        workflow_description
                        isCustomWorkflow
                      }
                    }`
	gqlReq.Variables.ChaosWorkFlowInput = in

	query, _ := json.Marshal(gqlReq)

	resp, _ := SendRequest(
		SendRequestParams{
			Endpoint: cred.Endpoint + utils.GQLAPIPath,
			Token:    cred.Token,
		},
		query,
		string(types.Post),
	)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(bodyBytes))
}
