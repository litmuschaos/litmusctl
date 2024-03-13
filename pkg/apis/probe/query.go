package probe

const (
	ListProbeQuery = `query ListProbes($projectID: ID!, $probeNames: [ID!], $filter: ProbeFilterInput) {
		listProbes(projectID: $projectID, probeNames: $probeNames, filter: $filter) {
		  name
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
	DeleteProbeQuery = `mutation deleteProbe($probeName: ID!, $projectID: ID!) {
		deleteProbe(probeName: $probeName, projectID: $projectID)
	  }
	`
)
