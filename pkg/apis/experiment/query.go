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
                          description
                          weightages {
                            faultName
                            weightage
                          }
                          isCustomExperiment
                          updatedAt
                          createdAt
                          infra {
                            projectID
                            name
                            infraID
                            infraType
                          }
                          isRemoved
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
                          projectID
                          infraID
                          infraType
                          }
                          isRemoved
                          updatedAt
                          updatedBy{
                              username
                              email
                          }
                          phase
                          resiliencyScore
                          faultsPassed
                          faultsFailed
                          faultsAwaited
                          faultsStopped
                          faultsNa
                          totalFaults
                          executionData
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
