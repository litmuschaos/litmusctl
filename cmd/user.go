package cmd

type NewUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	CompanyName string `json:"company_name"`
	Name        string `json:"name"`
	ProjectName string `json:"project_name"`
}
