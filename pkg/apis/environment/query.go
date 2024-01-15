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
							updatedBy{
								username
							}
							infraIDs
							type
						}
					}
	               }`
)
