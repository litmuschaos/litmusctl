package constants

const (

	// CLI version
	CLIVersion = "v0.1.0"

	// Default username
	DefaultUsername = "admin"

	// Default installation mode
	DefaultMode = "cluster"

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

	// Agent type is "external" for agents connected via litmusctl
	AgentType = "external"

	// Default namespace for agent installation
	DefaultNs = "litmus"

	// Default service account used for agent installation
	DefaultSA = "litmus"

	// Chaos agent registration yaml path
	//ChaosYamlPath = "chaos/api/graphql/file"
	ChaosYamlPath = "api/file"

	ChaosAgentPath = "targets"
)
