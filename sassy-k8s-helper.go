package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var ns string
	flag.StringVar(&ns, "namespace", "", "namespace")
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	log.Println("Using kubeconfig file: ", kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalln("failed to get nodes:", err)
	}
	// print nodes
	for i, node := range nodes.Items {

		/* 	for _, condition := range node.Status.Conditions {
			fmt.Printf("%s:%s\n", condition.Type, condition.Status)
		} */
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			println("%v", err)
		}
		nodeNames := []string{}
		nodeNames = append(nodeNames, node.Name)
		for _, namespace := range namespaces.Items {
			for _, name := range nodeNames {
				// pods need a namespace to be listed.
				pods, err := clientset.CoreV1().Pods(namespace.Name).List(context.TODO(), metav1.ListOptions{
					TypeMeta:             metav1.TypeMeta{},
					LabelSelector:        "",
					FieldSelector:        "spec.nodeName=" + name,
					Watch:                false,
					AllowWatchBookmarks:  false,
					ResourceVersion:      "",
					ResourceVersionMatch: "",
					TimeoutSeconds:       new(int64),
					Limit:                0,
					Continue:             "",
				})
				if err != nil {
					println("%v", err)
				}
				for _, pod := range pods.Items {
					if pod.Status.Phase != "Running" {
						fmt.Printf("[%d] %s\n", i, node.GetName())
						for _, condition := range node.Status.Conditions {
							fmt.Printf("%s:%s\n", condition.Type, condition.Status)
						}
						fmt.Printf("Pods not running:\n")
						fmt.Printf("%s: %s\n", pod.Namespace, pod.Name)
						fmt.Printf("%s\n", pod.Status.Phase)
					}

				}
			}
		}
	}
}
