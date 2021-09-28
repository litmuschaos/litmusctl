package apis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/litmuschaos/litmusctl/pkg/k8s"
	"github.com/litmuschaos/litmusctl/pkg/types"
	"github.com/litmuschaos/litmusctl/pkg/utils"
	"io/ioutil"
	"net/http"
)

type manifestData struct {
	Data data `json:"data`
}

type data struct {
	GetManifest string `json:"getManifest"`
}

func GetManifest(c context.Context, cred types.Credentials) (string, error) {

	// EXTRACT clusterID FROM NAMESPACE BY READING AGENT CONFIG
	query := `{"query":"query {\n getManifest(projectID : \"02b81577-f2a0-4e23-ac2a-d446d1aec59b\",\n clusterID : \"055396c9-0a86-440b-942b-a76d8a7112dd\",\n accessKey :\"gNvZEZxd0EtVtMweDRWJgQniNQMGvUuz\")}"}`
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
		//
		//for _, x := range manifestsArray {
		//	x
		//}
		//fmt.Println("m0 ", manifestsArray[1])

		// NAMESPACE IS HARDCODED
		response, err := k8s.ClusterResource(c, manifest.Data.GetManifest, "litmus")
		if err != nil {
			fmt.Println("ERROR",err)
			return "", err
		}

		content := response.UnstructuredContent()
		agentData := content["data"].(map[string]interface{})

		fmt.Println("Agent Data ", agentData)

		//fmt.Println("-----------------------------------------------")
		//fmt.Println("MAP",response.UnstructuredContent())

		//fmt.Println("AGENT CONFIG FROM MANIFEST ",response)

		fmt.Println("-----------------------------------------------")

		configData,err := k8s.GetConfigMap1(c)
		if err != nil {
			fmt.Println("ERROR",err)
			return "", err
		}
		fmt.Println("Config Data ", configData)

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


		//err = ioutil.WriteFile("temp1.yaml", []byte(manifest.Data.GetManifest), 0644)
		//if err != nil {
		//	return "", err
		//}
		//
		//yamlOutput, err := k8s.ApplyYaml(k8s.ApplyYamlPrams{
		//	Token:    cred.Token,
		//	Endpoint: cred.Endpoint,
		//	YamlPath: "temp1.yaml",
		//}, "")
		//
		//if err != nil {
		//	fmt.Println("ERROR")
		//	return "", err
		//}
		//utils.White_B.Print("\n", yamlOutput)

		fmt.Println("SUCCESSFUL")
		return "File Written Succesfully", nil
	} else {
		return "", errors.New("Unmatched status code:" + string(bodyBytes))
	}
}
