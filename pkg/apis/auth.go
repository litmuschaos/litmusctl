package apis

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

func Auth(input types.AuthInput) (types.AuthResponse, error) {
	type Payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	payloadBytes, err := json.Marshal(Payload{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		return types.AuthResponse{}, err
	}

	resp, err := SendRequest(SendRequestParams{input.Endpoint + "/auth/login", ""}, payloadBytes)
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

	return types.AuthResponse{}, nil
}
