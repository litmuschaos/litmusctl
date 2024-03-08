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
	GetProbeQuery = `query getProbe($projectID: ID!, $probeName: ID!) {
		getProbe(projectID: $projectID, probeName: $probeName) {
		  name
		  description
		  type
		  infrastructureType
		  createdAt
		  createdBy{
			username
		  }
		  updatedAt
		  updatedBy{
			username
		  }
		  tags
		}
	  }
	`
)
