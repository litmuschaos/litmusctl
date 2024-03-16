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
			url
			method{
				get{
					criteria
					responseCode
				}
				post{
					contentType
					body
					bodyPath
					criteria
					responseCode
				}
			}
			insecureSkipVerify
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
			command
			comparator {
				type
				value
				criteria
			}
			source
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
			group
			version
			resource
			namespace
			resourceNames
			fieldSelector
			labelSelector
			operation
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
			endpoint
			query
			queryPath
			comparator{
				type
				value
				criteria
			}
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
