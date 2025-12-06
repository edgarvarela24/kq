// cmd/pods.go - The 'pods' subcommand
//
// This command:
// 1. Lists namespaces (lets user select one)
// 2. Lists pods in that namespace (fuzzy search)
// 3. Shows actions for the selected pod (coming in Phase 4)
package cmd

import (
	"fmt"
	"os"

	"github.com/evarela/kq/internal/kube"
	"github.com/evarela/kq/internal/ui"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var podActions = []string{
	"logs",
}

// podsCmd represents the pods command
//
// When you run `kq pods`, this command executes.
var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Interactively select and act on pods",
	Long: `List pods with fuzzy search and perform actions on them.

Actions available:
  - logs: View pod logs (with follow, timestamps options)
  - exec: Open an interactive shell
  - describe: Show detailed pod information
  - port-forward: Forward local port to pod`,

	// Run is called when the command executes
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement the pods command flow
		//
		// Step 1: Create the Kubernetes client
		//   - Use kube.NewClient()
		//   - Handle errors (print and exit)
		client, err := kube.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}
		// Step 2: Get the namespace
		//   - Check if --namespace flag was provided (skip prompt if so)
		//   - Otherwise, list namespaces with kube.ListNameSpaces()
		//   - Let user select one with ui.SelectOne()
		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespaces, err := kube.ListNameSpaces(client)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing namespaces: %v\n", err)
				os.Exit(1)
			}
			namespace, err = ui.SelectOne("Select Namespace", namespaces)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error selecting namespace: %v\n", err)
				os.Exit(1)
			}
		}

		// Step 3: List and select a pod
		//   - Use kube.ListPods(clientset, namespace)
		//   - Let user select one with ui.SelectOne()
		pods, err := kube.ListPods(client, namespace)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing pods in namespace %q: %v\n", namespace, err)
			os.Exit(1)
		}

		if len(pods) == 0 {
			fmt.Fprintf(os.Stderr, "No pods found in namespace %q\n", namespace)
			os.Exit(1)
		}

		pod, err := ui.SelectOne("Select Pod", pods)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting pod: %v\n", err)
			os.Exit(1)
		}
		// Step 4: List pod actions
		//   - Use the podActions slice defined above
		//   - Let user select one with ui.SelectOne()
		action, err := ui.SelectOne("Select Action", podActions)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting action: %v\n", err)
			os.Exit(1)
		}
		// Step 5: Perform the selected action
		switch action {
		case "logs":
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			executeLogsAction(client, namespace, pod, dryRun)
		default:
			fmt.Fprintf(os.Stderr, "Action %q not implemented yet\n", action)
			os.Exit(1)
		}
	},
}

func executeLogsAction(client kubernetes.Interface, namespace, pod string, dryRun bool) {
	if dryRun {
		// Print the kubectl command that would be run
		fmt.Printf("kubectl logs -n %s %s\n", namespace, pod)
		return
	}
	containers, err := kube.ListContainers(client, namespace, pod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing containers in pod %q: %v\n", pod, err)
		os.Exit(1)
	}
	var container string
	if len(containers) == 1 {
		container = containers[0]
	} else {
		// Prompt user to select a container if multiple exist
		container, err = ui.SelectOne("Select Container", containers)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting container: %v\n", err)
			os.Exit(1)
		}
	}

	// Prompt for log options
	follow, timestamps, previous, err := ui.SelectLogOptions()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error selecting log options: %v\n", err)
		os.Exit(1)
	}
	// Create PodLogOptions
	logOpts := kube.PodLogOptions{
		Follow:     follow,
		Timestamps: timestamps,
		Previous:   previous,
		Container:  container,
	}
	// Stream logs
	err = kube.GetPodLogs(client, namespace, pod, logOpts, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error streaming logs from pod %q: %v\n", pod, err)
		os.Exit(1)
	}
}

// init registers any flags specific to the pods command
func init() {
	// === LOCAL FLAGS FOR PODS ===
	// These only apply to `kq pods`, not other commands
	//
	// We don't have any pods-specific flags yet, but you could add:
	// podsCmd.Flags().BoolP("all-namespaces", "A", false, "List pods in all namespaces")
}
