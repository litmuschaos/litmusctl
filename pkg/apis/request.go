package apis

import (
	"bytes"
	"net/http"
)

type SendRequestParams struct {
	Endpoint string
	Token    string
}

func SendRequest(params SendRequestParams, payload []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", params.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", params.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
