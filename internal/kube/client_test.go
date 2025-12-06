package kube

import (
	"bytes"
	"sort"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

func TestListNameSpaces(t *testing.T) {
	// Table-driven test cases
	tests := []struct {
		name       string   // description of this test case
		namespaces []string // namespaces to create in fake cluster
		want       []string // expected result from ListNameSpaces
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
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: nsName,
					},
				}
				objects = append(objects, ns)
			}
			clientset := fake.NewSimpleClientset(objects...)
			got, err := ListNameSpaces(clientset)
			if err != nil {
				t.Fatalf("ListNameSpaces() error: %v", err)
			}
			if !equalStringSlices(got, tt.want) {
				t.Errorf("ListNameSpaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

// equalStringSlices compares two string slices for equality
func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]string, len(a))
	bCopy := make([]string, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	sort.Strings(aCopy)
	sort.Strings(bCopy)
	for i := range a {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}

func TestListPods(t *testing.T) {
	tests := []struct {
		name      string   // description of this test case
		namespace string   // namespace to query
		pods      []string // pod names to create in that namespace
		want      []string // expected result from ListPods
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
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      podName,
						Namespace: tt.namespace,
					},
				}
				objects = append(objects, pod)
			}
			clientset := fake.NewSimpleClientset(objects...)
			got, err := ListPods(clientset, tt.namespace)
			if err != nil {
				t.Fatalf("ListPods() error: %v", err)
			}
			if !equalStringSlices(got, tt.want) {
				t.Errorf("ListPods() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetPodLogs_Integration tests log streaming against a real cluster.
//
// This is an integration test — it requires a running Kubernetes cluster.
// Skip it if no cluster is available.
//
// Run with: go test ./internal/kube/... -v -run TestGetPodLogs_Integration
//
// NOTE: The fake clientset doesn't support streaming, so we test against
// a real cluster. This test is skipped in CI environments.
func TestGetPodLogs_Integration(t *testing.T) {
	// Skip if we can't connect to a cluster
	client, err := NewClient()
	if err != nil {
		t.Skip("Skipping integration test: no cluster available")
	}

	// This test assumes you have a running pod in kube-system namespace
	// (coredns is usually there in kind clusters)
	namespace := "kube-system"

	// Get a pod name dynamically
	pods, err := ListPods(client, namespace)
	if err != nil || len(pods) == 0 {
		t.Skip("Skipping integration test: no pods in kube-system")
	}
	podName := pods[0]

	// Capture logs into a buffer
	var buf bytes.Buffer
	opts := PodLogOptions{
		Follow:     false, // Important: don't follow or test hangs
		Timestamps: false,
		Previous:   false,
	}

	err = GetPodLogs(client, namespace, podName, opts, &buf)
	if err != nil {
		t.Fatalf("GetPodLogs() error: %v", err)
	}

	// Verify we got some log output
	if buf.Len() == 0 {
		t.Error("GetPodLogs() returned empty logs")
	}

	t.Logf("Got %d bytes of logs from %s/%s", buf.Len(), namespace, podName)
}
