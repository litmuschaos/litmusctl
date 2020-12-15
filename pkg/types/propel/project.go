package propel

type PropelProjects struct {
	Errors []Errors `json:"errors"`
	Data   Data     `json:"data"`
}
type Members struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
type Projects struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Members []Members `json:"members"`
}
type ListProjects struct {
	Projects []Projects `json:"projects"`
}
type Data struct {
	ListProjects ListProjects `json:"listProjects"`
}
