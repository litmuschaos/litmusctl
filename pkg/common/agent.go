package common

type Agent struct {
	AgentName      string `json:"cluster_name"`
	Mode           string
	Description    string `json:"description,omitempty"`
	PlatformName   string `json:"platform_name"`
	ProjectId      string `json:"project_id"`
	ClusterType    string `json:"cluster_type"`
	Namespace      string
	ServiceAccount string
	NsExists       bool
	SAExists       bool
}
