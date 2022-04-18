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
package types

// CreateChaosWorkFlowInput describes the payload of the API request for workflow creation.
type CreateChaosWorkFlowInput struct {
	WorkflowID          string            `json:"workflow_id"`
	WorkflowManifest    string            `json:"workflow_manifest"`
	CronSyntax          string            `json:"cronSyntax"`
	WorkflowName        string            `json:"workflow_name"`
	WorkflowDescription string            `json:"workflow_description"`
	Weightages          []WeightagesInput `json:"weightages"`
	IsCustomWorkflow    bool              `json:"isCustomWorkflow"`
	ProjectID           string            `json:"project_id"`
	ClusterID           string            `json:"cluster_id"`
}

type WeightagesInput struct {
	ExperimentName string `json:"experiment_name"`
	Weightage      int    `json:"weightage"`
}

type CreateChaosWorkFlowGraphQLRequest struct {
	Query     string `json:"query"`
	Variables struct {
		CreateChaosWorkFlowInput CreateChaosWorkFlowInput `json:"ChaosWorkFlowInput"`
	} `json:"variables"`
}
