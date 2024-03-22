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
		  kubernetesHTTPProperties{
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
		  }
		  kubernetesCMDProperties{
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
		  }
		  k8sProperties {
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
		  }
		  promProperties {
			probeTimeout
			interval
			retry
			attempt
			probePollingInterval
			initialDelay
			evaluationTimeout
			stopOnFailure
		  }
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
