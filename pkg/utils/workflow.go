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
	"net/url"

	"github.com/argoproj/argo-workflows/v3/pkg/apis/workflow/v1alpha1"
)

// ReadWorkflowManifest reads the manifest that is passed as an argument.
// It can be either a local file or a remote file.
func ReadWorkflowManifest(file string, workflow *v1alpha1.Workflow) error {
	parsedURL, err := url.ParseRequestURI(file)
	if err != nil || !(parsedURL.Scheme == "http" || parsedURL.Scheme == "https") {
		err = UnmarshalLocalFile(file, &workflow)
	} else {
		err = UnmarshalRemoteFile(file, &workflow)
	}
	return err
}
