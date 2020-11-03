package cmd

import (
	"fmt"
	"os/exec"
)

func ApplyYaml(c AgentRegistrationData, cred Credentials) (output string, err error) {
	path := fmt.Sprintf("%s/api/file/%s.yaml", cred.Host, c.Data.UserAgentReg.Token)
	args := []string{"kubectl", "apply", "-f", path}
	stdout, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("Error: %v", err)
	}
	return string(stdout), err
}
