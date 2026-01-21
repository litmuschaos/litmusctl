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
                          projectID
                          experimentRunID
                          experimentType
                          experimentID
                          weightages {
                              faultName
                              weightage
                          }
                          updatedAt
                          createdAt
                          infra {
                              projectID
                              infraID                              
                              name
                              description
                              tags
                              environmentID
                              platformName
                              isActive
                              isInfraConfirmed
                              isRemoved
                              updatedAt
                              createdAt
                              noOfExperiments
                              noOfExperimentRuns
                              token
                              infraNamespace
                              serviceAccount
                              infraScope
                              infraNsExists
                              infraSaExists
                              lastExperimentTimestamp
                              startTime
                              version
                              infraType
                              updateStatus
                          }
                          experimentName
                          phase
                          resiliencyScore
                          faultsPassed
                          faultsFailed
                          faultsAwaited
                          faultsStopped
                          faultsNa
                          totalFaults
                          executionData
                          isRemoved
                          updatedBy{
                              userID
                              username
                              email
                          }
                          createdBy{
                              userID
                              username
                              email
                          }
                          notifyID
                          runSequence
                        }
                      }
                    }`

	ExperimentRunsQuery = `query getExperimentRun($projectID: ID!, $experimentRunID: ID, $notifyID: ID) {
                        getExperimentRun(
                            projectID: $projectID
                            experimentRunID: $experimentRunID
                            notifyID: $notifyID
                          ) 
                          {
                          experimentRunID
                          experimentID
                          experimentName
                          infra {
                          name
                          infraID
                          }
                          updatedAt
                          updatedBy{
                              username
                          }
                          phase
                          resiliencyScore
                          executionData
                        }
                      }`
	DeleteExperimentQuery = `mutation deleteChaosExperiment($projectID: ID!, $experimentID: String!, $experimentRunID: String) {
                      deleteChaosExperiment(
                        projectID: $projectID
                        experimentID: $experimentID
                        experimentRunID: $experimentRunID
                      )
                    }`
	GetPodLogsQuery = `subscription podLog($request: PodLogRequest!) {
                      getPodLog(request: $request) {
                        log
                        __typename
                      }
                    }`
)
