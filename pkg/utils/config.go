package utils

import (
	"errors"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func CreateNewLitmusCtlConfig(filename string, config types.LitmuCtlConfig) error {

	configByte, err := yaml.Marshal(config)
	if err != nil{
		return err
	}

	_, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		return err
	}

	err = ioutil.WriteFile(filename, configByte, 0644)
	if err != nil{
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
		return types.LitmuCtlConfig{}, errors.New("File reading error "+ err.Error())
	}

	obj := &types.LitmuCtlConfig{}
	err = yaml.Unmarshal(data, obj)
	if err != nil {
		return types.LitmuCtlConfig{}, errors.New("File format not correct " + err.Error())
	}

	return *obj, nil
}

func ConfigSyntaxCheck(filename string) error {

	obj, err:= YamltoObject(filename)
	if err != nil{
		return err
	}

	if obj.APIVersion != "v1" || obj.Kind != "Config" {
		return errors.New("File format not correct")
	}

	return nil
}

func UpdateLitmusCtlConfig(account types.Account, filename string) error {
	obj, err:= YamltoObject(filename)
	if err != nil {
		return err
	}

	var outerflag = false
	for i, act := range obj.Accounts {
		if act.Endpoint == account.Endpoint {
			var innerflag = false
			for _, user := range act.Users{
				if user.Username == account.Users[0].Username {
					user.Password = account.Users[0].Username
					user.Token = account.Users[0].Token
					innerflag, outerflag = true, true
				}
			}

			if !innerflag {
				obj.Accounts[i].Users = append(obj.Accounts[i].Users, account.Users[0])
				outerflag = true
			}
		}
	}

	if !outerflag {
		obj.Accounts = append(obj.Accounts, account)
	}

	_, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
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