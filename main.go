package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// from https://github.com/kubernetes/client-go/blob/master/examples/in-cluster-client-configuration/main.go
func main() {
	explorer := NewApiExplorer()

	for {
		fmt.Printf("\n\n========== \nRunning discovery at %s\n========== \n", time.Now())

		explorer.PrintPodCount()
		explorer.PrintDeprecatedApis()

		time.Sleep(1 * time.Hour)
	}
}

type ApiExplorer struct {
	clientset        *kubernetes.Clientset
	dynamicClientSet dynamic.Interface
}

func NewApiExplorer() *ApiExplorer {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return &ApiExplorer{
		clientset:        clientset,
		dynamicClientSet: dynamicClientset,
	}
}

func (a *ApiExplorer) PrintDeprecatedApis() {
	gvrs := []schema.GroupVersionResource{
		{Group: "apps", Version: "v1", Resource: "deployments"},
		{Group: "storage.k8s.io", Version: "v1beta1", Resource: "csistoragecapacities"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta2", Resource: "flowschemas"},
		{Group: "flowcontrol.apiserver.k8s.io", Version: "v1beta2", Resource: "prioritylevelconfigurations"},
	}

	for _, g := range gvrs {
		ri := a.dynamicClientSet.Resource(g)
		fmt.Printf("\n=== Retrieving: %s.%s.%s ===\n", g.Resource, g.Version, g.Group)
		rs, err := ri.List(context.Background(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Failed to retrieve: %s: %s", g, err)
			continue
		}

		for _, i := range rs.Items {
			if i.GetNamespace() != "" {
				fmt.Printf("Found %s in namespace: %s\n", i.GetName(), i.GetNamespace())
			} else {
				fmt.Printf("Found %s\n", i.GetName())
			}
		}
	}

}

func (a *ApiExplorer) PrintPodCount() {
	// get pods in all the namespaces by omitting namespace
	// Or specify namespace to get pods in particular namespace
	pods, err := a.clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// Examples for error handling:
	// - Use helper functions e.g. errors.IsNotFound()
	// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
	_, err = a.clientset.CoreV1().Pods("default").Get(context.TODO(), "example-xxxxx", metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod example-xxxxx not found in default namespace\n")
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %v\n", statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found example-xxxxx pod in default namespace\n")
	}
}
