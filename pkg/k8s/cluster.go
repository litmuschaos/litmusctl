package k8s

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"log"
	"os"
	"strings"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/restmapper"

	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	memory "k8s.io/client-go/discovery/cached"
)

var (
	decUnstructured = yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	dr              dynamic.ResourceInterface
	AgentNamespace  = os.Getenv("AGENT_NAMESPACE")
)

// This function handles cluster operations
func ClusterResource(c context.Context, manifest string, namespace string) (*unstructured.Unstructured, error) {

	// Getting dynamic and discovery client
	discoveryClient, dynamicClient, err := GetDynamicAndDiscoveryClient()
	if err != nil {
		fmt.Println("ERROR1 ", err)
		return nil, err
	}

	// Create a mapper using dynamic client
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	// Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}

	manifestsArray := strings.Split(manifest, "---")

	var gvk *schema.GroupVersionKind

	for _, x := range manifestsArray[1:] {

		_, gvk, err = decUnstructured.Decode([]byte(x), nil, obj)
		if err != nil {
			return nil, err
		}
		if obj.GetName() == "agent-config" {
			break
		}
	}
		// Find GVR
		mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
		if err != nil {
			return nil, err
		}

		//deploymentRes := schema.GroupVersionResource{Group:"apps", Version: "v1", Resource: "deployments"}

		// Obtain REST interface for the GVR
		if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
			// namespaced resources should specify the namespace
			fmt.Println("inside if")
			dr = dynamicClient.Resource(mapping.Resource).Namespace(namespace)
		} else {
			// for cluster-wide resources
			dr = dynamicClient.Resource(mapping.Resource)
		}

		response, err := dr.Get(c, obj.GetName(), metaV1.GetOptions{})
		if k8serrors.IsAlreadyExists(err) {
			// This doesn't ever happen even if it does already exist
			log.Print("Already exists")
			return nil, nil
		}

		if err != nil {
			return nil, err
		}
	log.Println("Resource successfully created")

	return response, nil
}


// This function handles cluster operations
func ClusterResource1(c context.Context, manifest string, namespace string) (*unstructured.Unstructured, error) {

	// Getting dynamic and discovery client
	discoveryClient, dynamicClient, err := GetDynamicAndDiscoveryClient()
	if err != nil {
		fmt.Println("ERROR1 ", err)
		return nil, err
	}

	// Create a mapper using dynamic client
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(discoveryClient))
	// Decode YAML manifest into unstructured.Unstructured
	obj := &unstructured.Unstructured{}


	var gvk *schema.GroupVersionKind
		_, gvk, err = decUnstructured.Decode([]byte(manifest), nil, obj)
		if err != nil {
			return nil, err
		}
		fmt.Println("namee ",obj.GetName())

	// Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return nil, err
	}
	fmt.Println("mapping ", mapping)

	//deploymentRes := schema.GroupVersionResource{Group:"apps", Version: "v1", Resource: "deployments"}

	// Obtain REST interface for the GVR
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		fmt.Println("inside if")
		dr = dynamicClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		// for cluster-wide resources
		dr = dynamicClient.Resource(mapping.Resource)
	}

	fmt.Println("dr ", dr)
	fmt.Println("namesas ", obj.GetName())
	response, err := dr.Get(c, obj.GetName(), metaV1.GetOptions{})
	if k8serrors.IsAlreadyExists(err) {
		// This doesn't ever happen even if it does already exist
		log.Print("Already exists")
		return nil, nil
	}

	if err != nil {
		return nil, err
	}
	log.Println("Resource successfully created")

	return response, nil
}
