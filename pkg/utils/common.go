package utils

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/config"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
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

type ApplyYamlPrams struct {
	Token    string
	Endpoint string
	YamlPath string
}

func ApplyYaml(params ApplyYamlPrams, kubeconfig string) (output string, err error) {
	path := fmt.Sprintf("%s/%s/%s.yaml", params.Endpoint, params.YamlPath, params.Token)

	var args []string
	if kubeconfig != "" {
		args = []string{"kubectl", "apply", "-f", path, "--kubeconfig", kubeconfig}
	} else {
		args = []string{"kubectl", "apply", "-f", path}
	}

	stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdout), err
}

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func GetCredentials(filename string) (types.Credentials, error) {
	obj, err := config.YamltoObject(filename)
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

	fmt.Println(out.String())

}

func PrintInYamlFormat(inf interface{}) {
	byt, err := yaml.Marshal(inf)
	PrintError(err)

	fmt.Println(string(byt))
}
