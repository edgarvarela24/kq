package kube

import (
	"bytes"
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListNamespaces(t *testing.T) {
	tests := []struct {
		name       string
		namespaces []string
		want       []string
	}{
		{
			name:       "returns empty slice when no namespaces",
			namespaces: []string{},
			want:       []string{},
		},
		{
			name:       "returns single namespace",
			namespaces: []string{"default"},
			want:       []string{"default"},
		},
		{
			name:       "returns multiple namespaces",
			namespaces: []string{"default", "kube-system", "staging"},
			want:       []string{"default", "kube-system", "staging"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objects []runtime.Object
			for _, nsName := range tt.namespaces {
				objects = append(objects, &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{Name: nsName},
				})
			}

			clientset := fake.NewSimpleClientset(objects...)
			got, err := ListNamespaces(context.Background(), clientset)
			if err != nil {
				t.Fatalf("ListNamespaces() error: %v", err)
			}
			if !equalStringSlices(got, tt.want) {
				t.Errorf("ListNamespaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aMap := make(map[string]int)
	for _, s := range a {
		aMap[s]++
	}
	for _, s := range b {
		if aMap[s] == 0 {
			return false
		}
		aMap[s]--
	}
	return true
}

func TestListPods(t *testing.T) {
	tests := []struct {
		name      string
		namespace string
		pods      []string
		want      []string
	}{
		{
			name:      "returns empty slice when no pods",
			namespace: "default",
			pods:      []string{},
			want:      []string{},
		},
		{
			name:      "returns single pod",
			namespace: "default",
			pods:      []string{"nginx-abc123"},
			want:      []string{"nginx-abc123"},
		},
		{
			name:      "returns multiple pods",
			namespace: "default",
			pods:      []string{"nginx-abc123", "redis-def456", "busybox-ghi789"},
			want:      []string{"nginx-abc123", "redis-def456", "busybox-ghi789"},
		},
		{
			name:      "returns only pods from specified namespace",
			namespace: "staging",
			pods:      []string{"web-xyz789"},
			want:      []string{"web-xyz789"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var objects []runtime.Object
			for _, podName := range tt.pods {
				objects = append(objects, &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: tt.namespace,
					},
				})
			}

			clientset := fake.NewSimpleClientset(objects...)
			got, err := ListPods(context.Background(), clientset, tt.namespace)
			if err != nil {
				t.Fatalf("ListPods() error: %v", err)
			}
			if !equalStringSlices(got, tt.want) {
				t.Errorf("ListPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListContainers(t *testing.T) {
	tests := []struct {
		name       string
		namespace  string
		podName    string
		containers []string
		want       []string
	}{
		{
			name:       "returns single container",
			namespace:  "default",
			podName:    "nginx-abc123",
			containers: []string{"nginx"},
			want:       []string{"nginx"},
		},
		{
			name:       "returns multiple containers",
			namespace:  "default",
			podName:    "multi-container",
			containers: []string{"app", "sidecar", "init"},
			want:       []string{"app", "sidecar", "init"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			containers := make([]corev1.Container, len(tt.containers))
			for i, name := range tt.containers {
				containers[i] = corev1.Container{Name: name}
			}

			pod := &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tt.podName,
					Namespace: tt.namespace,
				},
				Spec: corev1.PodSpec{
					Containers: containers,
				},
			}

			clientset := fake.NewSimpleClientset(pod)
			got, err := ListContainers(context.Background(), clientset, tt.namespace, tt.podName)
			if err != nil {
				t.Fatalf("ListContainers() error: %v", err)
			}
			if !equalStringSlices(got, tt.want) {
				t.Errorf("ListContainers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListContainers_PodNotFound(t *testing.T) {
	clientset := fake.NewSimpleClientset()
	_, err := ListContainers(context.Background(), clientset, "default", "nonexistent")
	if err == nil {
		t.Error("ListContainers() expected error for nonexistent pod, got nil")
	}
}

// TestGetPodLogs_Integration tests log streaming against a real cluster.
// Run with: go test ./internal/kube/... -v -run TestGetPodLogs_Integration
// NOTE: Skipped if no cluster is available.
func TestGetPodLogs_Integration(t *testing.T) {
	client, err := NewClient()
	if err != nil {
		t.Skip("Skipping integration test: no cluster available")
	}

	namespace := "kube-system"
	ctx := context.Background()

	pods, err := ListPods(ctx, client, namespace)
	if err != nil || len(pods) == 0 {
		t.Skip("Skipping integration test: no pods in kube-system")
	}

	var buf bytes.Buffer
	opts := PodLogOptions{
		Follow:     false,
		Timestamps: false,
		Previous:   false,
	}

	err = GetPodLogs(ctx, client, namespace, pods[0], opts, &buf)
	if err != nil {
		t.Fatalf("GetPodLogs() error: %v", err)
	}

	if buf.Len() == 0 {
		t.Error("GetPodLogs() returned empty logs")
	}

	t.Logf("Got %d bytes of logs from %s/%s", buf.Len(), namespace, pods[0])
}
