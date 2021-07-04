package apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

type createProjectResponse struct {
	Data struct {
		CreateProject struct {
			Name string `json:"name"`
		} `json:"createProject"`
	} `json:"data"`
}

func CreateProjectRequest(projectName string, cred types.Credentials) error {
	query := `{"query":"mutation{createProject(projectName: \"` + projectName + `\"){name}}"}`

	resp, err := SendRequest(cred.Endpoint+"/api/query", cred.Token, []byte(query))
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var project createProjectResponse
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return err
		}

		fmt.Println("project/" + project.Data.CreateProject.Name + " created")
		return nil
	} else {
		return errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

type listProjectResponse struct {
	Data struct {
		ListProjects []struct {
			ID        string `json:"id"`
			Name      string `json:"name"`
			CreatedAt string `json:"created_at"`
		} `json:"listProjects"`
	} `json:"data"`
}

func ListProject(cred types.Credentials) (listProjectResponse, error) {
	query := `{"query":"query{listProjects{id name created_at}}"}`
	resp, err := SendRequest(cred.Endpoint+"/api/query", cred.Token, []byte(query))
	if err != nil {
		return listProjectResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return listProjectResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data listProjectResponse
		err = json.Unmarshal(bodyBytes, &data)
		if err != nil {
			return listProjectResponse{}, err
		}

		return data, nil
	} else {
		return listProjectResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

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
	ID        string `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(c types.Credentials) (ProjectDetails, error) {
	query := `{"query":"query {\n  getUser(username: \"` + c.Username + `\"){\n projects{\n id\n name\n}\n}\n}"}`
	resp, err := SendRequest(c.Endpoint+"/api/query", c.Token, []byte(query))
	if err != nil {
		return ProjectDetails{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return ProjectDetails{}, err
	}

	if resp.StatusCode == http.StatusOK {
		var project ProjectDetails
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return ProjectDetails{}, err
		}

		return project, nil
	} else {
		return ProjectDetails{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}

	return ProjectDetails{}, nil
}
