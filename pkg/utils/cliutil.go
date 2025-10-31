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
	"fmt"
)

// FormatError formats error messages consistently with appropriate emoji and styling
func FormatError(message string, err error) string {
	return fmt.Sprintf("❌ %s: %v", message, err)
}

// PrintFormattedError prints a formatted error message to stderr in red color
func PrintFormattedError(message string, err error) {
	Red.Println(FormatError(message, err))
}
