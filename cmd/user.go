package cmd

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type NewUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
}

type GetUserData struct {
	Data Data `json:"data"`
}
type Projects struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type GetUser struct {
	Projects []Projects `json:"projects"`
	ID       string     `json:"id"`
	Username string     `json:"username"`
}
type Data struct {
	GetUser GetUser `json:"getUser"`
}

// GetUserDetails fetches details of the input user
func GetUserDetails(t Token, c Credentials) (GetUserData, interface{}) {
	var new GetUserData
	client := resty.New()
	bodyData := `{"query":" \nquery{\n  getUser(username: \"` + fmt.Sprintf("%s", c.Username) + `\"){\n    projects{\n      id\n      name\n    }\n    id\n    username\n  }\n}"}`
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
				"%s/api/query",
				c.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return GetUserData{}, resp.Error()
	}

	return new, nil
}
