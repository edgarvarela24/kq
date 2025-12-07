// Package kube provides Kubernetes client utilities for connecting to clusters
// and interacting with resources.
package kube

import (
	"context"
	"errors"
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
// It checks $KUBECONFIG first, then falls back to ~/.kube/config.
func NewClient() (kubernetes.Interface, error) {
	kubeconfigPath := os.Getenv("KUBECONFIG")
	if kubeconfigPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("getting user home directory: %w", err)
		}
		kubeconfigPath = filepath.Join(homeDir, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("building config from kubeconfig: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}
	return clientset, nil
}

// ListNamespaces returns the names of all namespaces in the cluster.
func ListNamespaces(ctx context.Context, clientset kubernetes.Interface) ([]string, error) {
	nsList, err := clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing namespaces: %w", err)
	}

	names := make([]string, len(nsList.Items))
	for i, ns := range nsList.Items {
		names[i] = ns.Name
	}
	return names, nil
}

// ListPods returns the names of all pods in a given namespace.
func ListPods(ctx context.Context, clientset kubernetes.Interface, namespace string) ([]string, error) {
	podList, err := clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("listing pods: %w", err)
	}

	pods := make([]string, len(podList.Items))
	for i, pod := range podList.Items {
		pods[i] = pod.Name
	}
	return pods, nil
}

// PodLogOptions configures how logs are retrieved.
type PodLogOptions struct {
	Follow     bool
	Timestamps bool
	Previous   bool
	Container  string
}

// GetPodLogs streams logs from a pod to the provided writer.
func GetPodLogs(ctx context.Context, clientset kubernetes.Interface, namespace, podName string, opts PodLogOptions, writer io.Writer) error {
	logOpts := &corev1.PodLogOptions{
		Follow:     opts.Follow,
		Timestamps: opts.Timestamps,
		Previous:   opts.Previous,
		Container:  opts.Container,
	}

	req := clientset.CoreV1().Pods(namespace).GetLogs(podName, logOpts)
	reader, err := req.Stream(ctx)
	if err != nil {
		return fmt.Errorf("opening log stream: %w", err)
	}
	defer reader.Close()

	if _, err := io.Copy(writer, reader); err != nil {
		// Context cancellation is expected when user presses Ctrl+C during follow
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return fmt.Errorf("reading log stream: %w", err)
	}
	return nil
}

// ListContainers returns the names of all containers in a pod.
func ListContainers(ctx context.Context, clientset kubernetes.Interface, namespace, podName string) ([]string, error) {
	pod, err := clientset.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("getting pod: %w", err)
	}

	names := make([]string, len(pod.Spec.Containers))
	for i, c := range pod.Spec.Containers {
		names[i] = c.Name
	}
	return names, nil
}
