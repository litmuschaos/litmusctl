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
package main

import (
	"fmt"
	"os"

	rootCmd "github.com/litmuschaos/litmusctl/pkg/cmd/root"
)

var CLIVersion string

func main() {
	err := os.Setenv("CLIVersion", CLIVersion)
	if err != nil {
		fmt.Println("Failed to fetched CLIVersion")
	}

	rootCmd.Execute()
}
