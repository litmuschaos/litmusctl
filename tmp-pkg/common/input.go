package common

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/litmuschaos/litmusctl/tmp-pkg/common/k8s"

	"github.com/litmuschaos/litmusctl/tmp-pkg/constants"
	"golang.org/x/crypto/ssh/terminal"
)

func GetUsername() string {
	var username string
	fmt.Print("🤔 Username [", constants.DefaultUsername, "]: ")
	fmt.Scanln(&username)
	if username == "" {
		username = constants.DefaultUsername
	}
	return username
}

func GetPortalURL() (*url.URL, error) {
	var host string
	fmt.Print("👉 Host URL where litmus is installed: ")
	fmt.Scanln(&host)
	for host == "" {
		fmt.Println("⛔ Host URL can't be empty!!")
		fmt.Print("👉 Host URL: ")
		fmt.Scanln(&host)
	}
	host = strings.TrimRight(host, "/")
	newUrl, err := url.Parse(host)
	if err != nil {
		return &url.URL{}, err
	}
	return newUrl, nil
}

func GetPassword() []byte {
	var pass []byte
	fmt.Print("🙈 Password: ")
	pass, _ = terminal.ReadPassword(0)
	if pass == nil {
		fmt.Println("\n⛔ Password cannot be empty!")
		os.Exit(1)
	}
	return pass
}

// GetMode gets mode of agent installation as input
func GetMode() string {
	var mode int
	fmt.Println("\n🔌 Installation Modes:\n1. Cluster\n2. Namespace")
	fmt.Print("\n👉 Select Mode [", constants.DefaultMode, "]: ")
	fmt.Scanln(&mode)
	if mode == 0 {
		return "namespace"
	}
	for mode < 1 || mode > 2 {
		fmt.Println("🚫 Invalid mode. Please enter the correct mode")
		fmt.Print("👉 Select Mode [", constants.DefaultMode, "]: ")
		fmt.Scanln(&mode)
	}
	if mode == 1 {
		return "cluster"
	}
	return constants.DefaultMode
}

func Confirm() {
	var wish string
	fmt.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&wish)
	if wish == "Y" || wish == "Yes" || wish == "yes" || wish == "y" {
		fmt.Println("👍 Continuing agent connection!!")
	} else {
		fmt.Println("✋ Exiting agent connection!!")
		os.Exit(1)
	}
}

// getPlatformName displays a list of platforms, takes the
// platform name as input and returns the selected platform
//
// TODO --
// - Entering any character other than numbers returns 0. Input validation need to be done.
// - If input is given as "123abc", "abc" will be used for next user input. Buffer need to be read completely.
// - String literals like "AWS" are used at multiple places. Need to be changed to constants.
func GetPlatformName(kubeconfig *string) string {
	var platform int
	discoveredPlatform := DiscoverPlatform(kubeconfig)
	fmt.Println("📦 Platform List")
	fmt.Println(constants.PlatformList)
	fmt.Print("🔎 Select Platform [", discoveredPlatform, "]: ")
	fmt.Scanln(&platform)
	switch platform {
	case 0:
		return discoveredPlatform
	case 1:
		return "AWS"
	case 2:
		return "GKE"
	case 3:
		return "Openshift"
	case 4:
		return "Rancher"
	default:
		return constants.DefaultPlatform
	}
}

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

// Summary display the agent details based on input
func Summary(agent Agent, product string, kubeconfig *string) {
	fmt.Println("\n📌 Summary --------------------------")
	fmt.Println("\nAgent Name:        ", agent.AgentName)
	fmt.Println("Agent Description: ", agent.Description)
	fmt.Println("Platform Name:     ", agent.PlatformName)
	if ok, _ := k8s.NsExists(agent.Namespace, kubeconfig); ok {
		fmt.Println("Namespace:         ", agent.Namespace)
	} else {
		fmt.Println("Namespace:         ", agent.Namespace, "(new)")
	}
	if product == "chaos" {
		if k8s.SAExists(agent.Namespace, agent.ServiceAccount, kubeconfig) {
			fmt.Println("Service Account:   ", agent.ServiceAccount)
		} else {
			fmt.Println("Service Account:   ", agent.ServiceAccount, "(new)")
		}
		fmt.Println("Installation Mode: ", agent.Mode)
	}
	fmt.Println("\n-------------------------------------")
}
