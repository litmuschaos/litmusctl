package chaos

import (
	"fmt"

	util "github.com/litmuschaos/litmusctl/pkg/common"

	resty "github.com/go-resty/resty/v2"
)

type ProjectDetails struct {
	Data Data `json:"data"`
}
type Members struct {
	UserUID string `json:"user_uid"`
	Role    string `json:"role"`
}
type GetProjects struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Members []Members `json:"members"`
}
type Data struct {
	GetProjects []GetProjects `json:"getProjects"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(t util.Token, c util.Credentials, product string) (ProjectDetails, interface{}) {
	var new ProjectDetails
	client := resty.New()
	bodyData := `{"query":"\nquery{\n  getProjects{\n    id\n    name\n    members{\n      user_uid\n      role\n    }\n  }\n}"}`
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
				"%s/%s/api/graphql/query",
				c.Host,
				product,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return ProjectDetails{}, resp.Error()
	}

	return new, nil
}

// GetProject display list of projects and returns the project id based on input
func GetProject(u ProjectDetails) string {
	var pid int
	fmt.Println("\n✨ Projects List:")
	for index := range u.Data.GetProjects {
		projectNo := index + 1
		fmt.Printf("%d.  %s\n", projectNo, u.Data.GetProjects[index].Name)
	}
	fmt.Print("\n🔎 Select Project: ")
	fmt.Scanln(&pid)
	for pid < 1 || pid > len(u.Data.GetProjects) {
		fmt.Println("❗ Invalid Project. Please select a correct one.")
		fmt.Print("\n🔎 Select Project: ")
		fmt.Scanln(&pid)
	}
	pid = pid - 1
	return u.Data.GetProjects[pid].ID
}
