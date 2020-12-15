package cmd

const (

	// CLI version
	cliVersion = "v0.2.0"

	// Default username
	defaultUsername = "admin"

	// Default installation mode
	defaultMode = "namespace"

	// Default platform name
	defaultPlatform = "others"

	// Label of subscriber agent being deployed
	chaosAgentLabel = "app=subscriber"

	// Agent type is "external" for agents connected via kuberactl
	agentType = "external"

	// Default namespace for agent installation
	defaultNs = "kubera"

	// Default service account used for agent installation
	defaultSA = "kubera"

	// Chaos agent registration yaml path
	chaosYamlPath = "chaos/api/graphql/file"

	chaosAgentPath = "chaos/agents"
)

// Propel constants
const (

	// Propel agent type
	propelAgentType = "External"

	// Propel agent registration yaml path
	propelYamlPath = "propel/api/graphql/agent/gen"

	// Propel agent label
	propelAgentLabel = "propel.kubera.mayadata.io/app-name=propel-agent-subscriber"

	propelAgentPath = "propel/clusters/SelfCluster"
)
