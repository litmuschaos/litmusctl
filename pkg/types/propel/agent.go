package propel

type PropelAgentList struct {
	Data PropelAgentData `json:"data"`
}
type Agents struct {
	AgentID      string `json:"cluster_id"`
	AgentName    string `json:"cluster_name"`
	PlatformName string `json:"platform_name"`
	AgentType    string `json:"cluster_type"`
	Description  string `json:"description"`
	Token        string `json:"token"`
}
type ListAgents struct {
	Agents []Agents `json:"clusters"`
}
type PropelAgentData struct {
	ListAgents ListAgents `json:"listClusters"`
}

type AgentDetails struct {
	Errors    []Errors  `json:"errors"`
	AgentData AgentData `json:"data"`
}
type NewCluster struct {
	ClusterID          string `json:"cluster_id"`
	ProjectID          string `json:"project_id"`
	ClusterName        string `json:"cluster_name"`
	Description        string `json:"description"`
	PlatformName       string `json:"platform_name"`
	AccessKey          string `json:"access_key"`
	IsRegistered       bool   `json:"is_registered"`
	IsClusterConfirmed bool   `json:"is_cluster_confirmed"`
	IsActive           bool   `json:"is_active"`
	UpdatedAt          string `json:"updated_at"`
	CreatedAt          string `json:"created_at"`
	ClusterType        string `json:"cluster_type"`
	Token              string `json:"token"`
	ClusterYamlRoute   string `json:"cluster_yaml_route"`
	ClusterURL         string `json:"cluster_url"`
	IsSelfCluster      bool   `json:"is_self_cluster"`
	Namespace          string `json:"namespace"`
}
type AddCluster struct {
	ClusterID          string     `json:"cluster_id"`
	ClusterToken       string     `json:"cluster_token"`
	ClusterAccessKey   string     `json:"cluster_access_key"`
	IsClusterConfirmed bool       `json:"isClusterConfirmed"`
	YamlRoute          string     `json:"yaml_route"`
	NewCluster         NewCluster `json:"new_cluster"`
}
type AgentData struct {
	AddCluster AddCluster `json:"addCluster"`
}
