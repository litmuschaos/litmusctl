package util

import (
	"fmt"

	resty "github.com/go-resty/resty/v2"
	util "github.com/mayadata-io/cli-utils/pkg/common"
)

type PropelProjects struct {
	Errors []util.Errors `json:"errors"`
	Data   Data          `json:"data"`
}
type Projects struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type ListProjects struct {
	Projects []Projects `json:"projects"`
}
type Data struct {
	ListProjects ListProjects `json:"listProjects"`
}

// ListPropelProjects fetches the list of projects using listProjects query
func ListPropelProjects(t util.Token, c util.Credentials) (PropelProjects, interface{}) {
	var new PropelProjects
	client := resty.New()
	bodyData := `{"query":"query{\n  listProjects{\n    projects{\n      id\n      name\n }\n  }\n}"}`
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
				"%s/propel/api/graphql/query",
				c.Host,
			),
		)
	if err != nil || !resp.IsSuccess() {
		return PropelProjects{}, resp.Error()
	}

	return new, nil
}

// SelectPropelProject display list of projects and returns the project id based on selected project
func SelectPropelProject(u PropelProjects) string {
	var pid int
	fmt.Println("\n✨ Projects List:")
	for index, _ := range u.Data.ListProjects.Projects {
		projectNo := index + 1
		fmt.Printf("%d.  %s\n", projectNo, u.Data.ListProjects.Projects[index].Name)
	}
	fmt.Print("\n🔎 Select Project: ")
	fmt.Scanln(&pid)
	for pid < 1 || pid > len(u.Data.ListProjects.Projects) {
		fmt.Println("❗ Invalid Project. Please select a correct one.")
		fmt.Print("\n🔎 Select Project: ")
		fmt.Scanln(&pid)
	}
	pid = pid - 1
	return u.Data.ListProjects.Projects[pid].ID
}
