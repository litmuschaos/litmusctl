package cmd

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"

	"golang.org/x/crypto/ssh/terminal"
)

func getUsername() string {
	var username string
	fmt.Print("🤔 Username [", defaultUsername, "]: ")
	fmt.Scanln(&username)
	if username == "" {
		username = defaultUsername
	}
	return username
}

func getPortalURL() (*url.URL, error) {
	var host string
	fmt.Print("👉 Litmus Portal URL: ")
	fmt.Scanln(&host)
	for host == "" {
		fmt.Println("⛔ Portal URL can't be empty!!")
		fmt.Print("👉 Litmus Portal URL: ")
		fmt.Scanln(&host)
	}
	host = strings.TrimRight(host, "/")
	newUrl, err := url.Parse(host)
	if err != nil {
		return &url.URL{}, err
	}
	return newUrl, nil
}

func getPassword() []byte {
	var pass []byte
	fmt.Print("🙈 Password: ")
	pass, _ = terminal.ReadPassword(0)
	if pass == nil {
		fmt.Println("\n⛔ Password cannot be empty!")
		os.Exit(1)
	}
	return pass
}

func confirm() {
	var wish string
	fmt.Print("\n🤷 Do you want to continue with the above details? [Y/N]: ")
	fmt.Scanln(&wish)
	if wish == "Y" || wish == "Yes" || wish == "yes" || wish == "y" {
		fmt.Println("👍 Continuing agent registration!!")
	} else {
		fmt.Println("✋ Exiting agent registration!!")
		os.Exit(1)
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
func Summary(agent Agent) {
	fmt.Println("\n📌 Summary --------------------------")
	fmt.Println("\nAgent Name:        ", agent.AgentName)
	fmt.Println("Agent Description: ", agent.Description)
	fmt.Println("Platform Name:     ", agent.PlatformName)
	if ok, _ := NsExists(agent.Namespace); ok {
		fmt.Println("Namespace:         ", agent.Namespace)
	} else {
		fmt.Println("Namespace:         ", agent.Namespace, "(new)")
	}
	if SAExists(agent.Namespace, agent.ServiceAccount) {
		fmt.Println("Service Account:   ", agent.ServiceAccount)
	} else {
		fmt.Println("Service Account:   ", agent.ServiceAccount, "(new)")
	}
	fmt.Println("Installation Mode: ", agent.Mode)
	fmt.Println("\n-------------------------------------")
}
