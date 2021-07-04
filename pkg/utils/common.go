package utils

import (
	"bufio"
	"fmt"
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

func ApplyYaml(token string, endpoint string, yamlPath string, kubeconfig string) (output string, err error) {
	path := fmt.Sprintf("%s/%s/%s.yaml", endpoint, yamlPath, token)

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

func PrintError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
