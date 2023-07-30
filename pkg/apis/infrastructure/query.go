package infrastructure

const (
	DisconnectInfraQuery = `mutation deleteInfra($projectID: ID!, $infraID: String!) {
                      deleteInfra(
                        projectID: $projectID
                        infraID: $infraID
                      )
                    }`

	RegisterInfraQuery = `mutation registerInfra($projectID: ID!, $request: RegisterInfraRequest!) {
					  registerInfra(
						projectID: $projectID
						request: $request
					  ) {
						infraID
						name
						token
					  }
					}
					`
	ListInfraQuery = `query listInfras($projectID: ID!, $request: ListInfraRequest!){
					listInfras(projectID: $projectID, request: $request){
						totalNoOfInfras
						infras {
							infraID
							name
							isActive
							environmentID
						}
					}
					}`
)
