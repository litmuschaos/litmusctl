package apis

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"io/ioutil"
	"log"
	"net/http"
)

func Auth(input types.AuthInput)(types.AuthResponse, error){

	type Payload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data := Payload{
		Username: input.Username,
		Password: input.Password,
	}
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return types.AuthResponse{}, err
	}

	body := bytes.NewReader(payloadBytes)
	ep := input.Endpoint + "/auth/login"
	fmt.Println(ep)

	req, err := http.NewRequest("POST", ep, body)
	if err != nil {
		return types.AuthResponse{}, err
	}

	req.Header.Set("Content-Type", "application/json")


	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return types.AuthResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode == http.StatusOK {

		var authResponse types.AuthResponse
		err = json.Unmarshal(bodyBytes, &authResponse)
		if err != nil {
			return types.AuthResponse{}, err
		}

		return authResponse, nil
	} else {
		return types.AuthResponse{}, errors.New("Unmatached status code:" + string(bodyBytes))
	}

	return types.AuthResponse{}, nil
}