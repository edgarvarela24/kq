// internal/kube/client.go - Kubernetes client wrapper
//
// This package handles connecting to a Kubernetes cluster.
//
// Key concepts:
// - kubeconfig: A file (~/.kube/config) that stores cluster connection info
// - clientset: A collection of typed API clients for different k8s resources
//
// Docs:
// - client-go package: https://pkg.go.dev/k8s.io/client-go
// - clientcmd (config loading): https://pkg.go.dev/k8s.io/client-go/tools/clientcmd
// - kubernetes.Clientset: https://pkg.go.dev/k8s.io/client-go/kubernetes#Clientset
package kube

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient creates a Kubernetes clientset from the user's kubeconfig.
//
// The kubeconfig file contains:
// - Cluster endpoints (API server URLs)
// - Authentication credentials
// - Context mappings (which user talks to which cluster)
//
// Loading order (same as kubectl):
// 1. $KUBECONFIG environment variable (if set)
// 2. ~/.kube/config (default location)
func NewClient() (*kubernetes.Clientset, error) {
	// Step 1: Find kubeconfig path
	// Check $KUBECONFIG first, then fall back to ~/.kube/config
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	// Step 2: Build REST config from kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}
	// Step 3: Create clientset from config
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}
	return clientSet, nil
}

func ListNameSpaces(clientset kubernetes.Interface) ([]string, error) {
	// Step 1: Get the list of namespaces from the cluster
	ctx := context.Background()
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Step 2: Extract and return the namespace names as a slice of strings
	names := make([]string, len(nsList.Items))
	for i, ns := range nsList.Items {
		names[i] = ns.Name
	}
	return names, nil
}

// ListPods returns the names of all pods in a given namespace.
func ListPods(clientset kubernetes.Interface, namespace string) ([]string, error) {
	ctx := context.Background()
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	pods := make([]string, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = pod.Name
	}
	return pods, nil
}

// PodLogOptions configures how logs are retrieved
type PodLogOptions struct {
	Follow     bool
	Timestamps bool
	Previous   bool
	Container  string // Empty string means default/only container
}

// GetPodLogs streams logs from a pod to the provided writer.
//
// TODO: Implement this function
//
// This is different from List functions — instead of returning data,
// it streams data to an io.Writer (usually os.Stdout).
//
// Key concepts:
// - GetLogs() returns a *rest.Request, not the logs directly
// - Call .Stream(ctx) on the request to get an io.ReadCloser
// - Use io.Copy() to pipe from the reader to your writer
// - Don't forget to close the stream!
//
// Docs:
// - PodInterface.GetLogs: https://pkg.go.dev/k8s.io/client-go/kubernetes/typed/core/v1#PodInterface
// - corev1.PodLogOptions: https://pkg.go.dev/k8s.io/api/core/v1#PodLogOptions
// - io.Copy: https://pkg.go.dev/io#Copy
func GetPodLogs(clientset kubernetes.Interface, namespace, podName string, opts PodLogOptions, writer io.Writer) error {
	// 1. Create a corev1.PodLogOptions struct with the options
	logOpts := &corev1.PodLogOptions{
		Follow:     opts.Follow,
		Timestamps: opts.Timestamps,
		Previous:   opts.Previous,
		Container:  opts.Container,
	}
	// 2. Call clientset.CoreV1().Pods(namespace).GetLogs(podName, &opts)
	ctx := context.Background()
	logs := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOpts)
	// 3. Call .Stream(ctx) on the request to get a reader
	reader, err := logs.Stream(ctx)
	if err != nil {
		return fmt.Errorf("failed to stream logs: %w", err)
	}
	defer reader.Close()
	// 4. Use io.Copy(writer, reader) to stream the logs
	_, err = io.Copy(writer, reader)
	if err != nil {
		return fmt.Errorf("failed to copy logs: %w", err)
	}
	return nil
}

func ListContainers(clientset kubernetes.Interface, namespace, podName string) ([]string, error) {
	ctx := context.Background()
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod: %w", err)
	}

	containerNames := make([]string, len(pod.Spec.Containers))
	for i, container := range pod.Spec.Containers {
		containerNames[i] = container.Name
	}
	return containerNames, nil
}
