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

type User struct {
	ExpiresIn string `yaml:"expires_in" json:"expires_in"`
	Token     string `yaml:"token" json:"token"`
	Username  string `yaml:"username" json:"username"`
}

type Account struct {
	Users          []User `yaml:"users" json:"users"`
	Endpoint       string `yaml:"endpoint" json:"endpoint"`
	ServerEndpoint string `yaml:"serverEndpoint" json:"serverEndpoint"`
}

type LitmuCtlConfig struct {
	Accounts       []Account `yaml:"accounts" json:"accounts"`
	APIVersion     string    `yaml:"apiVersion" json:"apiVersion"`
	CurrentAccount string    `yaml:"current-account" json:"current-account"`
	CurrentUser    string    `yaml:"current-user" json:"current-user"`
	Kind           string    `yaml:"kind" json:"kind"`
}

type Current struct {
	CurrentAccount string `yaml:"current-account" json:"current-account"`
	CurrentUser    string `yaml:"current-user" json:"current-user"`
}

type UpdateLitmusCtlConfig struct {
	CurrentAccount string  `yaml:"current-account" json:"current-account"`
	CurrentUser    string  `yaml:"current-user" json:"current-user"`
	Account        Account `yaml:"account" json:"account"`
	ServerEndpoint string  `yaml:"serverEndpoint" json:"serverEndpoint"`
}
