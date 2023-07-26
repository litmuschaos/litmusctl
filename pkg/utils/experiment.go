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
package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	"github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"
	"sigs.k8s.io/yaml"
)

// ParseExperimentManifest reads the manifest that is passed as an argument and
// populates the payload for the Message API request. The manifest
// can be either a local file or a remote file.
func ParseExperimentManifest(file string, chaosWorkFlowRequest *model.SaveChaosExperimentRequest) error {

	var body []byte
	var err error

	// Read the manifest file.
	parsedURL, ok := url.ParseRequestURI(file)
	if ok != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		body, err = ioutil.ReadFile(file)
	} else {
		body, err = ReadRemoteFile(file)
	}
	if err != nil {
		return err
	}

	// Extract the kind of Argo Workflow from the given manifest
	re := regexp.MustCompile(`\bkind:\s*(?P<kind>Workflow|CronWorkflow)\b`)
	extractKind := fmt.Sprintf("${%s}", re.SubexpNames()[1])
	workflowKind := re.ReplaceAllString(re.FindString(string(body)), extractKind)

	if workflowKind == "Workflow" {

		var workflow v1alpha1.Workflow
		err = UnmarshalObject(body, &workflow)
		if err != nil {
			return err
		}

		if len(workflow.ObjectMeta.Name) > 0 {
			chaosWorkFlowRequest.Name = workflow.ObjectMeta.Name
		} else if len(workflow.ObjectMeta.GenerateName) > 0 {
			workflow.ObjectMeta.Name = workflow.ObjectMeta.GenerateName + generateRandomString()
			workflow.ObjectMeta.GenerateName = "TOBEDELETED"
			chaosWorkFlowRequest.Name = workflow.ObjectMeta.Name
		} else {
			return errors.New("No name or generateName provided for the Chaos experiment.")
		}

		// Marshal the workflow back to JSON for API payload.
		workflowStr, ok := json.Marshal(workflow)
		if ok != nil {
			return ok
		}

		chaosWorkFlowRequest.Manifest = strings.Replace(string(workflowStr), "\"generateName\":\"TOBEDELETED\",", "", 1)
		//chaosWorkFlowRequest.IsCustomWorkflow = true

		// Fetch the weightages for experiments present in the spec.
		err = FetchWeightages(chaosWorkFlowRequest, workflow.Spec.Templates)
		if err != nil {
			return err
		}
	} else if workflowKind == "CronWorkflow" {

		var cronWorkflow v1alpha1.CronWorkflow
		err = UnmarshalObject(body, &cronWorkflow)
		if err != nil {
			return err
		}

		chaosWorkFlowRequest.Name = cronWorkflow.ObjectMeta.Name

		if len(cronWorkflow.ObjectMeta.Name) > 0 {
			chaosWorkFlowRequest.Name = cronWorkflow.ObjectMeta.Name
		} else if len(cronWorkflow.ObjectMeta.GenerateName) > 0 {
			cronWorkflow.ObjectMeta.Name = cronWorkflow.ObjectMeta.GenerateName + generateRandomString()
			cronWorkflow.ObjectMeta.GenerateName = "TOBEDELETED"
			chaosWorkFlowRequest.Name = cronWorkflow.ObjectMeta.Name
		} else {
			return errors.New("No name or generateName provided for the Chaos experiment.")
		}

		// Marshal the workflow back to JSON for API payload.
		workflowStr, ok := json.Marshal(cronWorkflow)
		if ok != nil {
			return ok
		}

		chaosWorkFlowRequest.Manifest = strings.Replace(string(workflowStr), "\"generateName\":\"TOBEDELETED\",", "", 1)
		//chaosWorkFlowRequest.IsCustomWorkflow = true

		// Set the schedule for the workflow
		//chaosWorkFlowRequest.CronSyntax = cronWorkflow.Spec.Schedule

		// Fetch the weightages for experiments present in the spec.
		err = FetchWeightages(chaosWorkFlowRequest, cronWorkflow.Spec.WorkflowSpec.Templates)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid resource kind found in manifest.")
	}

	return nil
}

// Helper function to check the presence of a string in a slice
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Helper function to generate a random 8 char string - used for workflow name postfix
func generateRandomString() string {
	rand.Seed(time.Now().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvxyz0123456789")
	b := make([]rune, 5)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// FetchWeightages takes in the templates present in the workflow spec and
// assigns weightage to each of the experiments present in them. It can parse
// both artifacts and remote experiment specs.
func FetchWeightages(chaosWorkFlowRequest *model.SaveChaosExperimentRequest, templates []v1alpha1.Template) error {

	for _, t := range templates {

		var err error

		if t.Inputs.Artifacts != nil && len(t.Inputs.Artifacts) > 0 {
			if t.Inputs.Artifacts[0].Raw == nil {
				continue
			}

			var data = t.Inputs.Artifacts[0].Raw.Data
			if len(data) > 0 {
				// This replacement is required because chaos engine yaml have a syntax template. example:{{ workflow.parameters.adminModeNamespace }}
				// And it is not able the unmarshal the yamlstring to chaos engine struct
				data = strings.ReplaceAll(data, "{{", "")
				data = strings.ReplaceAll(data, "}}", "")

				var chaosEngine chaosTypes.ChaosEngine

				err = yaml.Unmarshal([]byte(data), &chaosEngine)

				if err != nil {
					return errors.New("failed to unmarshal chaosengine")
				}

				if strings.ToLower(chaosEngine.Kind) == "chaosengine" {
					var weightageInput model.WeightagesInput

					weightageInput.FaultName = chaosEngine.ObjectMeta.GenerateName

					if len(weightageInput.FaultName) == 0 {
						return errors.New("empty chaos experiment name")
					}

					if len(chaosEngine.Spec.Experiments) == 0 {
						return errors.New("no experiments specified in chaosengine - " + weightageInput.FaultName)
					}

					w, ok := t.Metadata.Labels["weight"]

					if !ok {
						White.Println("Weightage for ChaosFault/" + weightageInput.FaultName + " not provided, defaulting to 10.")
						w = "10"
					}
					weightageInput.Weightage, err = strconv.Atoi(w)

					if err != nil {
						return errors.New("Invalid weightage for ChaosExperiment/" + weightageInput.FaultName + ".")
					}

					//chaosWorkFlowRequest. = append(chaosWorkFlowRequest.Weightages, &weightageInput)
				}
			}
		}
	}

	// If no experiments are present in the workflow, adds a 0 to the Weightages array so it doesn't fail (same behavior as the UI)
	//if len(chaosWorkFlowRequest.Weightages) == 0 {
	//	White.Println("No experiments found in the chaos scenario, defaulting experiments weightage to 0.")
	//	var weightageInput model.WeightagesInput
	//	weightageInput.ExperimentName = ""
	//	weightageInput.Weightage = 0
	//	chaosWorkFlowRequest.Weightages = append(chaosWorkFlowRequest.Weightages, &weightageInput)
	//}
	return nil
}
