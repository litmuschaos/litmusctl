package cmd

import "fmt"

// GetProject display list of projects and returns the project id based on input
func GetProject(u GetUserData) string {
	var pid int
	fmt.Println("\n✨ Projects List:")
	for index, _ := range u.Data.GetUser.Projects {
		projectNo := index + 1
		fmt.Printf("%d.  %s\n", projectNo, u.Data.GetUser.Projects[index].Name)
	}
	fmt.Print("\n🔎 Select Project: ")
	fmt.Scanln(&pid)
	for pid < 1 || pid > len(u.Data.GetUser.Projects) {
		fmt.Println("❗ Invalid Project. Please select a correct one.")
		fmt.Print("\n🔎 Select Project: ")
		fmt.Scanln(&pid)
	}
	pid = pid - 1
	return u.Data.GetUser.Projects[pid].ID
}
