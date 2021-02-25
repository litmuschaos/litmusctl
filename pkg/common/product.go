package common

import (
	"fmt"

	resty "github.com/go-resty/resty/v2"
)

type LaunchProductResponse struct {
	Data LaunchProductData `json:"data"`
}
type LaunchProductData struct {
	LaunchProduct string `json:"launchProduct"`
}

// GetUserDetails fetches details of the input user
func LaunchProduct(t Token, c Credentials, Product string) (LaunchProductResponse, error) {
	var new LaunchProductResponse
	client := resty.New()
	bodyData := `{"query":"query{\n  launchProduct(type: ` + fmt.Sprintf("%s", Product) + `)\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", t.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(bodyData).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&new).
		Post(
			fmt.Sprintf(
				"%s/api/graphql/query",
				c.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return LaunchProductResponse{}, err
	}

	return new, nil
}
