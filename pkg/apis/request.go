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
	"bytes"
	"net/http"
)

type SendRequestParams struct {
	Endpoint string
	Token    string
}

func SendRequest(params SendRequestParams, payload []byte, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, params.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return &http.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", params.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
