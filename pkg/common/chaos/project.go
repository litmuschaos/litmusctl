package chaos

import (
	"fmt"

	util "github.com/litmuschaos/litmusctl/pkg/common"

	resty "github.com/go-resty/resty/v2"
)

type ProjectDetails struct {
	Data Data `json:"data"`
}
type Data struct {
	GetUser GetUser `json:"getUser"`
}
type GetUser struct {
	Projects []Project `json:"projects"`
}
type Project struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(t util.Token, c util.Credentials) (ProjectDetails, interface{}) {
	var user ProjectDetails

	client := resty.New()
	bodyData := `{"query":"query {\n  getUser(username: \"` + fmt.Sprintf("%s", c.Username) + `\"){\n projects{\n id\n name\n}\n}\n}"}`
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", fmt.Sprintf("%s", t.AccessToken)).
		SetHeader("Accept-Encoding", "gzip, deflate, br").
		SetBody(bodyData).
		// SetResult automatic unmarshalling for the request,
		// if response status code is between 200 and 299
		SetResult(&user).
		Post(
			fmt.Sprintf(
				"%s/api/query",
				c.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return ProjectDetails{}, err
	}
	return user, nil
}

// GetProject display list of projects and returns the project id based on input
func GetProject(u ProjectDetails) string {
	var pid int
	fmt.Println("\n✨ Projects List:")
	for index := range u.Data.GetUser.Projects {
		projectNo := index + 1
		fmt.Printf("%d.  %s\n", projectNo, u.Data.GetUser.Projects[index].Name)
	}
	fmt.Print("\n🔎 Select Project: ")
	fmt.Scanln(&pid)
	for pid < 1 || pid > len(u.Data.GetUser.Projects) {
		fmt.Println("❗ Invalid Project. Please select a correct one.")
		fmt.Print("\n🔎 Select Project: ")
		fmt.Scanln(&pid)
	}
	pid = pid - 1
	return u.Data.GetUser.Projects[pid].ID
}
