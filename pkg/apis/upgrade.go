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
	GetManifest string `json:"getInfraManifest"`
}

type GetInfraResponse struct {
	Data   GetInfraData `json:"data"`
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
}

type GetInfraData struct {
	GetInfraDetails InfraDetails `json:"getInfraDetails"`
}

type InfraDetails struct {
	InfraID        string  `json:"infraID"`
	InfraNamespace *string `json:"infraNamespace"`
}

func UpgradeInfra(c context.Context, cred types.Credentials, projectID string, infraID string, kubeconfig string) (string, error) {

	// Query to fetch Infra details from server
	query := `{"query":"query {\n getInfraDetails(infraID : \"` + infraID + `\", \n projectID : \"` + projectID + `\"){\n infraNamespace infraID \n}}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query), string(types.Post))
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	var infra GetInfraResponse

	if resp.StatusCode == http.StatusOK {
		err = json.Unmarshal(bodyBytes, &infra)
		if err != nil {
			return "", err
		}
		if len(infra.Errors) > 0 {
			return "", errors.New(infra.Errors[0].Message)
		}
	} else {
		return "", errors.New(resp.Status)
	}

	// Query to fetch upgraded manifest from the server
	query = `{"query":"query {\n getInfraManifest(projectID : \"` + projectID + `\",\n infraID : \"` + infra.Data.GetInfraDetails.InfraID + `\", \n upgrade: true)}"}`
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
		err = ioutil.WriteFile("chaos-infra-manifest.yaml", []byte(manifest.Data.GetManifest), 0644)
		if err != nil {
			return "", err
		}

		// Fetching subscriber-config from the subscriber
		configData, err := k8s.GetConfigMap(c, "subscriber-config", *infra.Data.GetInfraDetails.InfraNamespace)
		if err != nil {
			return "", err
		}
		var configMapString string

		metadata := new(bytes.Buffer)
		fmt.Fprintf(metadata, "\n%s: %s\n%s: %s\n%s: \n  %s: %s\n  %s: %s\n%s:\n", "apiVersion", "v1",
			"kind", "ConfigMap", "metadata", "name", "subscriber-config", "namespace", *infra.Data.GetInfraDetails.InfraNamespace, "data")

		for k, v := range configData {
			b := new(bytes.Buffer)
			if k == "COMPONENTS" {
				fmt.Fprintf(b, "  %s: |\n    %s", k, v)
			} else if k == "START_TIME" || k == "IS_INFRA_CONFIRMED" {
				fmt.Fprintf(b, "  %s: \"%s\"\n", k, v)
			} else {
				fmt.Fprintf(b, "  %s: %s\n", k, v)
			}
			configMapString = configMapString + b.String()

		}

		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    cred.Token,
			Endpoint: cred.Endpoint,
			YamlPath: "chaos-infra-manifest.yaml",
		}, kubeconfig, true)

		if err != nil {
			return "", err
		}
		utils.White.Print("\n", yamlOutput)

		err = os.Remove("chaos-infra-manifest.yaml")
		if err != nil {
			return "Error removing Chaos Infrastructure manifest: ", err
		}

		// Creating a backup for current subscriber-config in the SUBSCRIBER
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		configMapString = metadata.String() + configMapString
		err = ioutil.WriteFile(home+"/backupSubscriberConfig.yaml", []byte(configMapString), 0644)
		if err != nil {
			return "Error creating backup for subscriber config: ", err
		}

		utils.White_B.Print("\n ** A backup of subscriber-config configmap has been saved in your system's home directory as backupSubscriberConfig.yaml **\n")

		return "Manifest applied successfully", nil
	} else {
		return "GQL error: ", errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
