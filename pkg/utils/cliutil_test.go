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
	"errors"
	"testing"
)

func TestFormatError(t *testing.T) {
	err := errors.New("connection failed")
	result := FormatError("Error creating project", err)
	expected := "❌ Error creating project: connection failed"

	if result != expected {
		t.Errorf("FormatError() = %v, want %v", result, expected)
	}
}
