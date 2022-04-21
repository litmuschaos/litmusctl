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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"regexp"
	"strconv"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
	chaosTypes "github.com/litmuschaos/chaos-operator/pkg/apis/litmuschaos/v1alpha1"
	"github.com/litmuschaos/litmus/litmus-portal/graphql-server/graph/model"
	yamlutil "k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/yaml"
)

// ParseWorkflowManifest reads the manifest that is passed as an argument and
// populates the payload for the CreateChaosWorkflow API request. The manifest
// can be either a local file or a remote file.
func ParseWorkflowManifest(file string, chaosWorkFlowInput *model.ChaosWorkFlowInput) error {

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
		chaosWorkFlowInput.WorkflowManifest = string(workflowStr)
		chaosWorkFlowInput.WorkflowName = workflow.ObjectMeta.Name
		chaosWorkFlowInput.IsCustomWorkflow = true

		// Fetch the weightages for experiments present in the spec.
		err = FetchWeightages(chaosWorkFlowInput, workflow.Spec.Templates)
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
		chaosWorkFlowInput.WorkflowManifest = string(workflowStr)
		chaosWorkFlowInput.WorkflowName = cronWorkflow.ObjectMeta.Name
		chaosWorkFlowInput.IsCustomWorkflow = true

		// Set the schedule for the workflow
		chaosWorkFlowInput.CronSyntax = cronWorkflow.Spec.Schedule

		// Fetch the weightages for experiments present in the spec.
		err = FetchWeightages(chaosWorkFlowInput, cronWorkflow.Spec.WorkflowSpec.Templates)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid resource kind found in manifest.")
	}

	return nil
}

// FetchWeightages takes in the templates present in the workflow spec and
// assigns weightage to each of the experiments present in them. It can parse
// both artifacts and remote experiment specs.
func FetchWeightages(chaosWorkFlowInput *model.ChaosWorkFlowInput, templates []v1alpha1.Template) error {

	var chaosExperiments []chaosTypes.ChaosExperiment

	// Fetch all present experiments and append them to the experiments array
	for _, t := range templates {
		var c chaosTypes.ChaosExperiment

		// Only the template named "install-chaos-experiments" contains ChaosExperiment(s)
		if t.Name == "install-chaos-experiments" {

			if len(t.Inputs.Artifacts) != 0 {

				// These are experiments with the spec passed as an artifact
				for _, a := range t.Inputs.Artifacts {
					err := yaml.Unmarshal([]byte(a.Raw.Data), &c)
					if err != nil {
						return errors.New("Error parsing ChaosExperiment: " + err.Error())
					}
					chaosExperiments = append(chaosExperiments, c)
				}

			} else { // These are experiments with their spec passed as remote file URLs

				// Extracting the URL from the container arguments
				re := regexp.MustCompile(`kubectl apply -f\s(?P<first>.+)\s-n.+`)
				extractURL := fmt.Sprintf("${%s}", re.SubexpNames()[1])
				fileURL := re.ReplaceAllString(t.Container.Args[0], extractURL)

				body, err := ReadRemoteFile(fileURL)
				if err != nil {
					return errors.New("Error reading ChaosExperiment: " + err.Error())
				}

				// Using a decoder since there might be multiple YAML documents
				// inside the same remote file.
				decoder := yamlutil.NewYAMLToJSONDecoder(bytes.NewReader(body))
				for {
					if err := decoder.Decode(&c); err != nil {
						break
					}
					chaosExperiments = append(chaosExperiments, c)
				}
			}
		}
	}

	// Assign weights to each chaos experiment
	for _, c := range chaosExperiments {

		// Fetch experiment weightage from annotation
		w, ok := c.ObjectMeta.Annotations["litmuschaos.io/experiment-weightage"]
		if !ok {
			White.Println("Weightage for ChaosExperiment/" + c.ObjectMeta.Name + " not provided, defaulting to 10.")
			w = "10"
		}

		var win model.WeightagesInput
		var err error
		win.ExperimentName = c.ObjectMeta.Name
		win.Weightage, err = strconv.Atoi(w)
		if err != nil {
			return errors.New("Invalid weightage for ChaosExperiment/" + c.ObjectMeta.Name + ".")
		}
		chaosWorkFlowInput.Weightages = append(chaosWorkFlowInput.Weightages, &win)
	}

	return nil
}
