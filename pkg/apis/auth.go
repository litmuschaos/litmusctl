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
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

type Payload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Auth(input types.AuthInput) (types.AuthResponse, error) {
	payloadBytes, err := json.Marshal(Payload{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		return types.AuthResponse{}, err
	}

	// Sending token as empty because auth server doesn't need Authorization token to validate.
	resp, err := SendRequest(SendRequestParams{input.Endpoint + utils.AuthAPIPath + "/login", ""}, payloadBytes, string(types.Post))
	if err != nil {
		return types.AuthResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return types.AuthResponse{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var authResponse types.AuthResponse
		err = json.Unmarshal(bodyBytes, &authResponse)
		if err != nil {
			return types.AuthResponse{}, err
		}

		return authResponse, nil
	} else {
		return types.AuthResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
