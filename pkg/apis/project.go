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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/utils"
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

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query))
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
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query))
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
	resp, err := SendRequest(SendRequestParams{Endpoint: c.Endpoint + utils.GQLAPIPath, Token: c.Token}, []byte(query))
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
