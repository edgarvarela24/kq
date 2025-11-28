package kube

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

// TestListNameSpaces tests the ListNameSpaces function.
//
// This is a "table-driven test" — a Go idiom where you define test cases
// in a slice and loop through them. Benefits:
// - Easy to add new cases
// - Clear what each case tests
// - Shared test logic
//
// TODO: Implement the test
//
// Steps:
// 1. Create fake namespace objects (corev1.Namespace)
// 2. Create a fake clientset with those namespaces
// 3. Call ListNameSpaces with the fake clientset
// 4. Assert the returned names match what you expect
//
// Useful:
// - fake.NewSimpleClientset(objects...) creates a fake clientset
// - Objects must satisfy runtime.Object (k8s resources do)
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
			// TODO: Build fake namespace objects from tt.namespaces
			//
			// You need to create []runtime.Object containing *corev1.Namespace
			// Each namespace needs at minimum: ObjectMeta.Name
			var objects []runtime.Object
			for _, nsName := range tt.namespaces {
				ns := &corev1.Namespace{
					ObjectMeta: metav1.ObjectMeta{
						Name: nsName,
					},
				}
				objects = append(objects, ns)
			}
			// TODO: Create fake clientset
			// fake.NewSimpleClientset(objects...)
			clientset := fake.NewSimpleClientset(objects...)
			// TODO: Call ListNameSpaces
			got, err := ListNameSpaces(clientset)
			// TODO: Compare result with tt.want
			// Use t.Errorf() on mismatch
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
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
