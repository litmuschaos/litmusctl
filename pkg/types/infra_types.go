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

type Infra struct {
	InfraName      string `json:"infraName"`
	Mode           string
	Description    string `json:"description,omitempty"`
	PlatformName   string `json:"platformName"`
	EnvironmentID  string `json:"environmentID"`
	ProjectId      string `json:"projectID"`
	InfraType      string `json:"infraType"`
	NodeSelector   string `json:"nodeSelector"`
	Tolerations    string
	Namespace      string
	ServiceAccount string
	NsExists       bool
	SAExists       bool
	SkipSSL        bool
}

type Toleration struct {
	Key               string `json:"key"`
	Value             string `json:"value"`
	Operator          string `json:"operator"`
	Effect            string `json:"effect"`
	TolerationSeconds int    `json:"tolerationSeconds"`
}
