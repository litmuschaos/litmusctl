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

const (
	DefaultFileName = ".litmusconfig"

	// Default username
	DefaultUsername = "admin"

	// Default installation mode
	DefaultMode = "cluster"

	// Platform list
	PlatformList = "1. AWS\n2. GKE\n3. Openshift\n4. Rancher\n5. Others"

	// AWS identifier
	AWSIdentifier = "aws://"

	// GKE identifier
	GKEIdentifier = "gce://"

	// Openshift identifier
	OpenshiftIdentifier = "node.openshift.io/os_id"

	// Default platform name
	DefaultPlatform = "Others"

	// Label of subscriber infra being deployed
	ChaosInfraLabel = "app=subscriber"

	// infra type is "external" for infras connected via litmusctl
	InfraTypeExternal = "external"

	// member owner role
	MemberOwnerRole = "Owner"

	// member editor role
	MemberEditorRole = "Editor"

	// Default namespace for infra installation
	DefaultNs = "litmus"

	// Default service account used for infra installation
	DefaultSA = "litmus"

	// Chaos infra connection yaml path
	ChaosYamlPath = "/api/file"

	ChaosInfraPath = "targets"

	// Graphql server api path
	GQLAPIPath = "/query"

	// Auth server api path
	AuthAPIPath = ""
)
