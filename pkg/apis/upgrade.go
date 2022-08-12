package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type manifestData struct {
	Data   data `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type data struct {
	GetManifest string `json:"getManifest"`
}

type ClusterData struct {
	Data   GetAgentDetails `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type GetAgentDetails struct {
	GetAgentDetails ClusterDetails `json:"getAgentDetails"`
}

type ClusterDetails struct {
	ClusterID      string  `json:"clusterID"`
	AccessKey      string  `json:"accessKey"`
	AgentNamespace *string `json:"agentNamespace"`
}

func UpgradeAgent(c context.Context, cred types.Credentials, projectID string, clusterID string, kubeconfig string) (string, error) {

	// Query to fetch agent details from server
	query := `{"query":"query {\n getAgentDetails(clusterID : \"` + clusterID + `\", \n projectID : \"` + projectID + `\"){\n agentNamespace accessKey clusterID \n}}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query), string(types.Post))
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var agent ClusterData

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(bodyBytes, &agent)
		if err != nil {
			return "", err
		}
		if len(agent.Errors) > 0 {
			return "", errors.New(agent.Errors[0].Message)
		}
	} else {
		return "", errors.New(resp.Status)
	}

	// Query to fetch upgraded manifest from the server
	query = `{"query":"query {\n getManifest(projectID : \"` + projectID + `\",\n clusterID : \"` + agent.Data.GetAgentDetails.ClusterID + `\",\n accessKey :\"` + agent.Data.GetAgentDetails.AccessKey + `\")}"}`
	resp, err = SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query), string(types.Post))
	if err != nil {
		return "", err
	}

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	// Checks if status code is OK(200)
	if resp.StatusCode == http.StatusOK {
		var manifest manifestData
		err = json.Unmarshal(bodyBytes, &manifest)
		if err != nil {
			return "", err
		}

		if len(manifest.Errors) > 0 {
			return "", errors.New(manifest.Errors[0].Message)
		}

		// To write the manifest data into a temporary file
		err = ioutil.WriteFile("chaos-delegate-manifest.yaml", []byte(manifest.Data.GetManifest), 0644)
		if err != nil {
			return "", err
		}

		// Fetching agent-config from the subscriber
		configData, err := k8s.GetConfigMap(c, "agent-config", *agent.Data.GetAgentDetails.AgentNamespace)
		if err != nil {
			return "", err
		}
		var configMapString string

		metadata := new(bytes.Buffer)
		fmt.Fprintf(metadata, "\n%s: %s\n%s: %s\n%s: \n  %s: %s\n  %s: %s\n%s:\n", "apiVersion", "v1",
			"kind", "ConfigMap", "metadata", "name", "agent-config", "namespace", *agent.Data.GetAgentDetails.AgentNamespace, "data")

		for k, v := range configData {
			b := new(bytes.Buffer)
			if k == "COMPONENTS" {
				fmt.Fprintf(b, "  %s: |\n    %s", k, v)
			} else if k == "START_TIME" || k == "IS_CLUSTER_CONFIRMED" {
				fmt.Fprintf(b, "  %s: \"%s\"\n", k, v)
			} else {
				fmt.Fprintf(b, "  %s: %s\n", k, v)
			}
			configMapString = configMapString + b.String()

		}

		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    cred.Token,
			Endpoint: cred.Endpoint,
			YamlPath: "chaos-delegate-manifest.yaml",
		}, kubeconfig, true)

		if err != nil {
			return yamlOutput, err
		}
		utils.White.Print("\n", yamlOutput)

		err = os.Remove("chaos-delegate-manifest.yaml")
		if err != nil {
			return "Error removing Chaos Delegate manifest: ", err
		}

		// Creating a backup for current agent-config in the SUBSCRIBER
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		configMapString = metadata.String() + configMapString
		err = ioutil.WriteFile(home+"/backupAgentConfig.yaml", []byte(configMapString), 0644)
		if err != nil {
			return "Error creating backup for agent config: ", err
		}

		utils.White_B.Print("\n ** A backup of agent-config configmap has been saved in your system's home directory as backupAgentConfig.yaml **\n")

		return "Manifest applied successfully", nil
	} else {
		return "GQL error: ", errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
