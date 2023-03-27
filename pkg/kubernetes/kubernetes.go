package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/den-is/ktempl/pkg/logging"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Node struct {
	Name        string
	InternalIP  string
	Annotations map[string]string
	Labels      map[string]string
}

// Connects to kubernetes cluster
func Connect(kubeconfig *string) (*kubernetes.Clientset, error) {

	// var kubeconfig *string
	// TODO: add logic to normalize path especially if using ~
	var thisKubeconfig string
	if *kubeconfig == "" {
		thisKubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
	} else {
		thisKubeconfig = *kubeconfig
	}

	config, err := clientcmd.BuildConfigFromFlags("", thisKubeconfig)
	if err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "error", "Was not able to build kubernetes connection config: ", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "error", "Was not able to connect to kubernetes: ", err)
	}

	// TODO: which err?, function body has couple of errors. combine or ...
	return clientset, err

}

// Returns Final list of nodes corresponding user query
func GetHostList(conn *kubernetes.Clientset, namespace *string, selector *string, usePods *bool) *[]Node {

	result := []Node{}

	if *usePods {

		pods := QueryPods(conn, namespace, selector)
		nodesFromPods := GetPodsNodes(conn, pods)

		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "info", fmt.Sprintf("Got %d pods, running on %d nodes", len(pods.Items), len(*nodesFromPods)))

		for _, podNode := range *nodesFromPods {
			n := Node{}
			n.Name = podNode.Name
			n.InternalIP = GetNodeInternalIP(podNode)
			n.Labels = podNode.Labels
			n.Annotations = podNode.Annotations
			result = append(result, n)
		}

	} else {

		nodes := QueryNodes(conn, selector)

		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "info", fmt.Sprintf("Got %d nodes", len(nodes.Items)))

		for _, node := range nodes.Items {
			if IfHostReady(&node) {
				n := Node{}
				n.Name = node.Name
				n.InternalIP = GetNodeInternalIP(&node)
				n.Labels = node.Labels
				n.Annotations = node.Annotations
				result = append(result, n)
			}
		}

	}

	return &result

}

// Query Pods in given Namespace using provided LabelSelector if any
func QueryPods(clientset *kubernetes.Clientset, namespace *string, labels *string) *v1.PodList {

	listFilter := metav1.ListOptions{}

	if *labels != "" {
		listFilter.LabelSelector = *labels
	}

	pods, err := clientset.CoreV1().Pods(*namespace).List(context.TODO(), listFilter)
	if err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "error", "Was not able to get pods list. ", err)

	}

	return pods

}

// Returns slice of unique nodes on which Pods are running
func GetPodsNodes(clientset *kubernetes.Clientset, pods *v1.PodList) *[]*v1.Node {

	existingNodes := make(map[string]bool)
	var nodeItems []*v1.Node

	for _, pod := range pods.Items {
		nodeName := pod.Spec.NodeName
		node, err := clientset.CoreV1().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
		if err != nil {
			logging.LogWithFields(
				logging.Fields{
					"component": "kubernetes",
				}, "error", fmt.Sprintf("Was not able to get node for the pod %q. ", pod.Name), err)
		}
		if _, exists := existingNodes[node.Name]; !exists {
			existingNodes[node.Name] = true
			nodeItems = append(nodeItems, node)
		} else {
			continue
		}
	}

	return &nodeItems
}

// Query for Nodes using given LabelSelector if any
func QueryNodes(clientset *kubernetes.Clientset, labels *string) *v1.NodeList {

	listFilter := metav1.ListOptions{}

	if *labels != "" {
		listFilter.LabelSelector = *labels
	}

	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), listFilter)
	if err != nil {
		logging.LogWithFields(
			logging.Fields{
				"component": "kubernetes",
			}, "error", "Was not able to get nodes list. ", err)
	}

	return nodes

}

// Returns node InternalIP. returns as soon as first instance of InternalIP is found
func GetNodeInternalIP(node *v1.Node) string {
	for _, v := range node.Status.Addresses {
		if v.Type == "InternalIP" {
			return v.Address
		}
	}
	return ""
}

// Checks whether node is Ready, a.k.a healthy, or not
func IfHostReady(node *v1.Node) bool {

	lastMsg := node.Status.Conditions[len(node.Status.Conditions)-1]
	if lastMsg.Status == "True" && lastMsg.Type == "Ready" {
		return true
	}
	return false

}
