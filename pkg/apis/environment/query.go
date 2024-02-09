package environment

const (
	CreateEnvironmentQuery = `mutation createEnvironment($projectID: ID!, $request: CreateEnvironmentRequest!) {
					  createEnvironment(
						projectID: $projectID
						request: $request
					  ) {
						environmentID
						name	
					  }
					}
					`
	ListEnvironmentQuery = `query listEnvironments($projectID: ID!, $request: ListEnvironmentRequest) {
	                 listEnvironments(projectID: $projectID,request: $request){
						environments {
							environmentID
							name
							createdAt
							updatedAt
							createdBy{
								username
							  }
						}
					}
	               }`

	DeleteEnvironmentQuery = `mutation deleteEnvironment($projectID: ID!, $environmentID: ID!) {
					deleteEnvironment(
					projectID: $projectID
					environmentID: $environmentID
					)
				}`

	GetEnvironmentQuery = `query getEnvironment($projectID: ID!, $environmentID : ID!) {
	                 getEnvironment(projectID: $projectID,environmentID: $environmentID){
							environmentID
							name
							createdAt
							updatedAt
							createdBy{
								username
							  }
							updatedBy{
								username
							}
							infraIDs
							type
							tags
						}
	               }`
)
