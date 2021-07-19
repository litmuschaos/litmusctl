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
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"math/big"
	"os"
)

var (
	Red     = color.New(color.FgRed)
	White_B = color.New(color.FgWhite, color.Bold)
	White   = color.New(color.FgWhite)
)

func Scanner() string {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		return scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	return ""
}
func PrintError(err error) {
	if err != nil {
		Red.Println(err)
		os.Exit(1)
	}
}

func GetLitmusConfigPath(cmd *cobra.Command) string {
	configFilePath, err := cmd.Flags().GetString("config")
	PrintError(err)

	if configFilePath == "" {
		home, err := homedir.Dir()
		PrintError(err)

		configFilePath = home + "/" + DefaultFileName
	}

	return configFilePath
}

func GetCredentials(cmd *cobra.Command) (types.Credentials, error) {
	configFilePath := GetLitmusConfigPath(cmd)

	obj, err := config.YamltoObject(configFilePath)
	PrintError(err)

	if obj.CurrentUser == "" || obj.CurrentAccount == "" {
		return types.Credentials{}, errors.New("Current user or current account is not set")
	}

	var token string
	for _, account := range obj.Accounts {
		if account.Endpoint == obj.CurrentAccount {
			for _, user := range account.Users {
				if user.Username == obj.CurrentUser {
					token = user.Token
				}
			}
		}
	}

	return types.Credentials{
		Username: obj.CurrentUser,
		Token:    token,
		Endpoint: obj.CurrentAccount,
	}, nil
}

func PrintInJsonFormat(inf interface{}) {
	var out bytes.Buffer
	byt, err := json.Marshal(inf)
	PrintError(err)

	err = json.Indent(&out, byt, "", "  ")
	PrintError(err)

	White.Println(out.String())

}

func PrintInYamlFormat(inf interface{}) {
	byt, err := yaml.Marshal(inf)
	PrintError(err)

	White.Println(string(byt))
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}
