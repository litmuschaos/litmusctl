package experiment

const (
	SaveExperimentQuery = `mutation saveChaosExperiment($projectID: ID!, $request: SaveChaosExperimentRequest!) {
                      saveChaosExperiment(projectID: $projectID, request: $request)
                     }`

	ListExperimentQuery = `query listExperiment($projectID: ID!, $request: ListExperimentRequest!) {
                      listExperiment(projectID: $projectID, request: $request) {
                        totalNoOfExperiments
                        experiments {
                          experimentID
                          experimentManifest
                          cronSyntax
                          name
                          infra {
                            name
                            infraID
                          }
                          updatedBy{
                              username
                              email
                        }
                      }
                    }
	}`

	ListExperimentRunsQuery = `query listExperimentRuns($projectID: ID!, $request: ListExperimentRunRequest!) {
                      listExperimentRun(projectID: $projectID, request: $request) {
                        totalNoOfExperimentRuns
                        experimentRuns {
                          experimentRunID
                          experimentID
                          experimentName
                          infra {
                          name
                          }
                          updatedAt
                          updatedBy{
                              username
                          }
                          phase
                          resiliencyScore
                        }
                      }
                    }`
	DeleteExperimentQuery = `mutation deleteChaosExperiment($projectID: ID!, $experimentID: String!, $experimentRunID: String) {
                      deleteChaosExperiment(
                        projectID: $projectID
                        experimentID: $experimentID
                        experimentRunID: $experimentRunID
                      )
                    }`
)
