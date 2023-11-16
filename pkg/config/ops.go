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
package config

import (
	"errors"
	"io/ioutil"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/types"
	"gopkg.in/yaml.v2"
)

var (
	SkipSSLVerify bool   = false
	CACert        string = ""
)

func CreateNewLitmusCtlConfig(filename string, config types.LitmuCtlConfig) error {

	configByte, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	_, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, configByte, 0644)
	if err != nil {
		return err
	}

	return nil
}

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func GetFileLength(filename string) (int, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return -1, err
	}

	return len(string(data)), nil
}

func YamltoObject(filename string) (types.LitmuCtlConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return types.LitmuCtlConfig{}, errors.New("File reading error " + err.Error())
	}

	obj := &types.LitmuCtlConfig{}
	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return types.LitmuCtlConfig{}, errors.New("File format not correct " + err.Error())
	}

	return *obj, nil
}

func ConfigSyntaxCheck(filename string) error {

	obj, err := YamltoObject(filename)
	if err != nil {
		return err
	}

	if obj.APIVersion != "v1" || obj.Kind != "Config" {
		return errors.New("File format not correct")
	}

	return nil
}

func UpdateLitmusCtlConfig(litmusconfig types.UpdateLitmusCtlConfig, filename string) error {
	obj, err := YamltoObject(filename)
	if err != nil {
		return err
	}

	var outerflag = false
	for i, act := range obj.Accounts {
		if act.Endpoint == litmusconfig.Account.Endpoint {
			var innerflag = false
			obj.Accounts[i].ServerEndpoint = litmusconfig.ServerEndpoint
			for j, user := range act.Users {
				if user.Username == litmusconfig.Account.Users[0].Username {
					obj.Accounts[i].Users[j].Username = litmusconfig.Account.Users[0].Username
					obj.Accounts[i].Users[j].Token = litmusconfig.Account.Users[0].Token
					obj.Accounts[i].Users[j].ExpiresIn = litmusconfig.Account.Users[0].ExpiresIn
					innerflag, outerflag = true, true
				}
			}

			if !innerflag {
				obj.Accounts[i].Users = append(obj.Accounts[i].Users, litmusconfig.Account.Users[0])
				outerflag = true
			}
		}
	}

	if !outerflag {
		obj.Accounts = append(obj.Accounts, litmusconfig.Account)
	}

	obj.CurrentAccount = litmusconfig.CurrentAccount
	obj.CurrentUser = litmusconfig.CurrentUser

	err = writeObjToFile(obj, filename)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCurrent(current types.Current, filename string) error {
	obj, err := YamltoObject(filename)
	if err != nil {
		return err
	}

	obj.CurrentUser = current.CurrentUser
	obj.CurrentAccount = current.CurrentAccount

	err = writeObjToFile(obj, filename)
	if err != nil {
		return err
	}

	return nil
}

func writeObjToFile(obj types.LitmuCtlConfig, filename string) error {
	_, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	byteObj, err := yaml.Marshal(obj)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, byteObj, 0644)
	if err != nil {
		return err
	}

	return nil
}

func IsAccountExists(obj types.LitmuCtlConfig, username string, endpoint string) bool {
	for _, account := range obj.Accounts {
		if account.Endpoint == endpoint {
			for _, user := range account.Users {
				if username == user.Username {
					return true
				}
			}
		}
	}

	return false
}
