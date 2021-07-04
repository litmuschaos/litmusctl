package types

type User struct {
	ExpiresIn string `yaml:"expires_in",json:"expires_in"`
	//Password   string `yaml:"password",json:"password"`
	Token      string `yaml:"token",json:"token"`
	Username   string `yaml:"username",json:"username"`
}

type Account struct {
	Users []User `yaml:"users",json:"users"`
	Endpoint string `yaml:"endpoint",json:"endpoint"`
}

type LitmuCtlConfig struct {
	Accounts 	   []Account `yaml:"accounts",json:"accounts"`
	APIVersion     string `yaml:"apiVersion",json:"apiVersion"`
	CurrentAccount string `yaml:"current-account",json:"current-account"`
	CurrentUser string `yaml:"current-user",json:"current-user"`
	Kind           string `yaml:"kind",json:"kind"`
}

type Current struct {
	CurrentAccount string `yaml:"current-account",json:"current-account"`
	CurrentUser string `yaml:"current-user",json:"current-user"`
}

type UpdateLitmusCtlConfig struct {
	CurrentAccount string `yaml:"current-account",json:"current-account"`
	CurrentUser string `yaml:"current-user",json:"current-user"`
	Account 	 Account `yaml:"account",json:"account"`
}

const DefaultFileName = "litmusconfig.yaml"

