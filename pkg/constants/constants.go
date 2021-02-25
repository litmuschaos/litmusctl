package constants

const (

	// CLI version
	CLIVersion = "v0.2.0"

	// Default username
	DefaultUsername = "admin"

	// Default installation mode
	DefaultMode = "namespace"

	// Platform list
	PlatformList = "1. AWS\n2. GKE\n3. Openshift\n4. Rancher\n5. Others"

	// AWS identifier
	AWSIdentifier = "aws://"

	// GKE identifier
	GKEIdentifier = "gce://"

	// Openshift identifier
	OpenshiftIdentifier = "node.openshift.io/os_id"

	// Default platform name
	DefaultPlatform = "Others"

	// Label of subscriber agent being deployed
	ChaosAgentLabel = "app=subscriber"

	// Agent type is "external" for agents connected via Litmusctl
	AgentType = "external"

	// Default namespace for agent installation
	DefaultNs = "kubera"

	// Default service account used for agent installation
	DefaultSA = "kubera"

	// Chaos agent registration yaml path
	//ChaosYamlPath = "chaos/api/graphql/file"
	ChaosYamlPath = "api/file"

	ChaosAgentPath = "chaos/agents"
)

// Propel constants
const (

	// Propel agent type
	PropelAgentType = "External"

	// Propel agent registration yaml path
	PropelYamlPath = "propel/api/graphql/agent/gen"

	// Propel agent label
	PropelAgentLabel = "propel.kubera.mayadata.io/app-name=propel-agent-subscriber"

	PropelAgentPath = "propel/clusters/SelfCluster"
)
