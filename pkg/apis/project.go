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
	"io/ioutil"
	"net/http"

	"github.com/golang-jwt/jwt"

	"github.com/litmuschaos/litmusctl/pkg/utils"

	"github.com/litmuschaos/litmusctl/pkg/types"
)

type createProjectResponse struct {
	Data struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type createProjectPayload struct {
	ProjectName string `json:"projectName"`
}

func CreateProjectRequest(projectName string, cred types.Credentials) (createProjectResponse, error) {
	payloadBytes, err := json.Marshal(createProjectPayload{
		ProjectName: projectName,
	})

	if err != nil {
		return createProjectResponse{}, err
	}
	resp, err := SendRequest(SendRequestParams{cred.Endpoint + utils.AuthAPIPath + "/create_project", "Bearer " + cred.Token}, payloadBytes, string(types.Post))
	if err != nil {
		return createProjectResponse{}, err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return createProjectResponse{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var project createProjectResponse
		err = json.Unmarshal(bodyBytes, &project)
		if err != nil {
			return createProjectResponse{}, err
		}

		if len(project.Errors) > 0 {
			return createProjectResponse{}, errors.New(project.Errors[0].Message)
		}

		utils.White_B.Println("project/" + project.Data.Name + " created")
		return project, nil
	} else {
		return createProjectResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

type listProjectResponse struct {
	Data []struct {
		ID        string `json:"ProjectID"`
		Name      string `json:"Name"`
		CreatedAt int64  `json:"CreatedAt"`
	} `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

func ListProject(cred types.Credentials) (listProjectResponse, error) {

	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.AuthAPIPath + "/list_projects", Token: "Bearer " + cred.Token}, []byte{}, string(types.Get))
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

		if len(data.Errors) > 0 {
			return listProjectResponse{}, errors.New(data.Errors[0].Message)
		}

		return data, nil
	} else {
		return listProjectResponse{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}

type ProjectDetails struct {
	Data   Data `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type Data struct {
	ID       string    `json:"ID"`
	Projects []Project `json:"Projects"`
}

type Member struct {
	Role     string `json:"Role"`
	UserID   string `json:"UserID"`
	UserName string `json:"UserName"`
}

type Project struct {
	ID        string   `json:"ProjectID"`
	Name      string   `json:"Name"`
	CreatedAt int64    `json:"CreatedAt"`
	Members   []Member `json:"Members"`
}

// GetProjectDetails fetches details of the input user
func GetProjectDetails(c types.Credentials) (ProjectDetails, error) {
	token, _ := jwt.Parse(c.Token, nil)
	if token == nil {
		return ProjectDetails{}, nil
	}
	Username, _ := token.Claims.(jwt.MapClaims)["username"].(string)
	resp, err := SendRequest(SendRequestParams{Endpoint: "http://localhost:3000" + utils.AuthAPIPath + "/get_user_with_project/" + Username, Token: "Bearer " + c.Token}, []byte{}, string(types.Get))
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
		if len(project.Errors) > 0 {
			return ProjectDetails{}, errors.New(project.Errors[0].Message)
		}

		return project, nil
	} else {
		return ProjectDetails{}, errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
