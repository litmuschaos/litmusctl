package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
)

var decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

type manifestData struct {
	Data data `json:"data`
}

type data struct {
	GetManifest string `json:"getManifest"`
}

func createKeyValuePairs(m map[string]string) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=\"%s\"\n", key, value)
	}
	return b.String()
}

func GetManifest(c context.Context, cred types.Credentials, projectID string, namespace string) (string, error) {

	configData, err := k8s.GetConfigMap1(c, namespace)
	if err != nil {
		fmt.Println("ERROR", err)
		return "", err
	}

	clusterID := configData["CLUSTER_ID"]
	accessKey := configData["ACCESS_KEY"]

	// EXTRACT clusterID FROM NAMESPACE BY READING AGENT CONFIG
	query := `{"query":"query {\n getManifest(projectID : \"` + projectID + `\",\n clusterID : \"` + clusterID + `\",\n accessKey :\"` + accessKey + `\")}"}`
	resp, err := SendRequest(SendRequestParams{Endpoint: cred.Endpoint + utils.GQLAPIPath, Token: cred.Token}, []byte(query))
	if err != nil {
		return "", err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var manifest manifestData
		err = json.Unmarshal(bodyBytes, &manifest)

		//var y unstructured.Unstructured
		//err = json.Unmarshal(manifest, &y)
		//
		//x, err := strconv.Unquote(manifest.Data.GetManifest)
		//
		//manifest1, err := yaml.Marshal(&manifest.Data.GetManifest)

		if err != nil {
			return "", err
		}

		// NAMESPACE IS HARDCODED
		response, err := k8s.ClusterResource(c, manifest.Data.GetManifest, "litmus")
		fmt.Println("response", response)
		if err != nil {
			fmt.Println("ERROR", err)
			return "", err
		}

		content := response.UnstructuredContent()
		agentData := content["data"].(map[string]interface{})
		fmt.Println("Agent dataaa", agentData)
		//fmt.Println("-----------------------------------------------")
		//fmt.Println("MAP",response.UnstructuredContent())

		//fmt.Println("AGENT CONFIG FROM MANIFEST ",response)

		fmt.Println("-----------------------------------------------")

		//res, err := k8s.GetConfigMap1()
		//if err != nil {
		//	return "", err
		//}
		//fmt.Println("config map", res)
		//
		//response, err = k8s.ClusterResource1(c, res, "litmus")
		//if err != nil {
		//	fmt.Println("ERROR",err)
		//	return "", err
		//}
		//fmt.Println("Config RESPONSE ",response)

		err = ioutil.WriteFile("temp1.yaml", []byte(manifest.Data.GetManifest), 0644)
		if err != nil {
			return "", err
		}

		fmt.Println("Agent Data ", agentData)

		// Decode YAML manifest into unstructured.Unstructured
		obj := &unstructured.Unstructured{}

		manifestsArray := strings.Split(manifest.Data.GetManifest, "---")

		var agentConfig string

		for _, x := range manifestsArray[1:] {

			_, _, err = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme).Decode([]byte(x), nil, obj)
			if err != nil {
				return "", err
			}
			if obj.GetName() == "agent-config" {
				agentConfig = x
				break
			}

		}

		var configMapString string

		metadata := new(bytes.Buffer)
		fmt.Fprintf(metadata, "\n%s: %s\n%s: %s\n%s: \n  %s: %s\n  %s: %s\n%s:\n", "apiVersion", "v1",
			"kind", "ConfigMap", "metadata", "name", "agent-config", "namespace", "litmus", "data")

		for k, v := range configData {
			b := new(bytes.Buffer)
			if k == "VERSION" {
				fmt.Println("Version", agentData["VERSION"])
				fmt.Fprintf(b, "  %s: %s\n", k, agentData["VERSION"])
			} else if k == "COMPONENTS" {
				fmt.Fprintf(b, "  %s: |\n    %s", k, v)
			} else if k == "START_TIME" || k == "IS_CLUSTER_CONFIRMED" {
				fmt.Fprintf(b, "  %s: \"%s\"\n", k, v)
			} else {
				fmt.Fprintf(b, "  %s: %s\n", k, v)
			}
			configMapString = configMapString + b.String()
		}
		configMapString = metadata.String() + configMapString
		fmt.Println("Config Data ", configMapString)

		fileContent, err := ioutil.ReadFile("temp1.yaml")
		if err != nil {
			return "", err
		}

		var newContent = string(fileContent)

		newContent = strings.Replace(newContent, agentConfig, configMapString, -1)

		err = ioutil.WriteFile("temp1.yaml", []byte(newContent), 0644)
		if err != nil {
			return "", err
		}
		yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
			Token:    cred.Token,
			Endpoint: cred.Endpoint,
			YamlPath: "temp1.yaml",
		}, "")

		if err != nil {
			fmt.Println("ERROR")
			return "", err
		}
		utils.White_B.Print("\n", yamlOutput)

		fmt.Println("SUCCESSFUL")
		return "File Written Succesfully", nil
	} else {
		return "", errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
