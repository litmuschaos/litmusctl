package probe

const (
	ListProbeQuery = `query ListProbes($projectID: ID!, $infrastructureType: InfrastructureType, $probeNames: [ID!], $filter: ProbeFilterInput) {
		listProbes(projectID: $projectID, infrastructureType: $infrastructureType, probeNames: $probeNames, filter: $filter) {
		  name
		  description
		  type
		  createdAt
		  createdBy{
			username
		  }
		}
	  }
	`
)
