package propel

type Errors struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}
