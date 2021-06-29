package common

import (
	"fmt"
	"os/exec"
)

func ApplyYaml(token string, cred Credentials, yamlPath string, kubeconfig string) (output string, err error) {
	path := fmt.Sprintf("%s/%s/%s.yaml", cred.Host, yamlPath, token)

	var args []string
	if kubeconfig != "" {
		args = []string{"kubectl", "apply", "-f", path, "--kubeconfig", kubeconfig}
	} else {
		args = []string{"kubectl", "apply", "-f", path}
	}

	stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return string(stdout), err
}
