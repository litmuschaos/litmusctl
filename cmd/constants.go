package cmd

const (

	// CLI version
	cliVersion = "v0.1.0"

	// Default username
	defaultUsername = "admin"

	// Default installation mode
	defaultMode = "namespace"

	// Default platform name
	defaultPlatform = "others"

	// Label of subscriber agent being deployed
	agentLabel = "app=subscriber"

	// Agent type is "external" for agents connected via kuberactl
	agentType = "external"

	// Default namespace for agent installation
	defaultNs = "kubera"

	// Default service account used for agent installation
	defaultSA = "kubera"

	// Chaos agent registration yaml path
	chaosYamlPath = "chaos/api/graphql/file"
)
