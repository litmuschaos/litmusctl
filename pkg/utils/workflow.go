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
	"net/url"
	"regexp"
	"strconv"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/api/litmuschaos/v1alpha1"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	"sigs.k8s.io/yaml"
)

// ParseWorkflowManifest reads the manifest that is passed as an argument and
// populates the payload for the CreateChaosWorkflow API request. The manifest
// can be either a local file or a remote file.
func ParseWorkflowManifest(file string, chaosWorkFlowRequest *model.ChaosWorkFlowRequest) error {

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

		// Marshal the workflow back to JSON for API payload.
		workflowStr, ok := json.Marshal(workflow)
		if ok != nil {
			return ok
		}
		chaosWorkFlowRequest.WorkflowManifest = string(workflowStr)
		chaosWorkFlowRequest.WorkflowName = workflow.ObjectMeta.Name
		chaosWorkFlowRequest.IsCustomWorkflow = true

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

		// Marshal the workflow back to JSON for API payload.
		workflowStr, _ := json.Marshal(cronWorkflow)
		chaosWorkFlowRequest.WorkflowManifest = string(workflowStr)
		chaosWorkFlowRequest.WorkflowName = cronWorkflow.ObjectMeta.Name
		chaosWorkFlowRequest.IsCustomWorkflow = true

		// Set the schedule for the workflow
		chaosWorkFlowRequest.CronSyntax = cronWorkflow.Spec.Schedule

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

// Helper fucntion to check the presence of a stirng in a slice
func sliceContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// FetchWeightages takes in the templates present in the workflow spec and
// assigns weightage to each of the experiments present in them. It can parse
// both artifacts and remote experiment specs.
func FetchWeightages(chaosWorkFlowRequest *model.ChaosWorkFlowRequest, templates []v1alpha1.Template) error {

	var chaosEngines []string

	// Fetch all present experiments and append them to the experiments array
	for _, t := range templates {
		// Only the template named "install-chaos-experiments" contains ChaosExperiment(s)
		if t.Name == "install-chaos-experiments" {
			for _, a := range t.Inputs.Artifacts {
				chaosEngines = append(chaosEngines, a.Name)
			}
		}

		// Template that contains ChaosEngine manifest
		if sliceContains(chaosEngines, t.Name) {
			var weightageInput model.WeightagesInput
			var err error

			var chaosEngine chaosTypes.ChaosEngine
			err = yaml.Unmarshal([]byte(t.Inputs.Artifacts[0].Raw.Data), &chaosEngine)
			if err != nil {
				return errors.New("Error parsing ChaosEngine: " + err.Error())
			}
			weightageInput.ExperimentName = chaosEngine.ObjectMeta.GenerateName
			w, ok := t.Metadata.Labels["weight"]

			if !ok {
				White.Println("Weightage for ChaosExperiment/" + weightageInput.ExperimentName + " not provided, defaulting to 10.")
				w = "10"
			}
			weightageInput.Weightage, err = strconv.Atoi(w)
			if err != nil {
				return errors.New("Invalid weightage for ChaosExperiment/" + weightageInput.ExperimentName + ".")
			}

			chaosWorkFlowRequest.Weightages = append(chaosWorkFlowRequest.Weightages, &weightageInput)
		}
	}
	return nil
}
