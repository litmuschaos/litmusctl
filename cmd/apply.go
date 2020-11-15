package cmd

import (
	"fmt"
	"os/exec"
)

func ApplyYaml(c AgentRegistrationData, cred Credentials, yamlPath string) (output string, err error) {
	path := fmt.Sprintf("%s/%s/%s.yaml", cred.Host, yamlPath, c.Data.UserAgentReg.Token)
	args := []string{"kubectl", "apply", "-f", path}
	stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return string(stdout), err
}
