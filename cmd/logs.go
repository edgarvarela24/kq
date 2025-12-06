// cmd/logs.go - The 'logs' subcommand
//
// Two ways to use:
// 1. Power user:  kq logs -n <namespace> <pod-name>
// 2. Interactive: kq pods → select pod → choose "logs" action
//
// Both paths use the same underlying log streaming logic.
package cmd

import (
	"fmt"
	"os"

	"github.com/evarela/kq/internal/kube"
	"github.com/evarela/kq/internal/ui"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [pod-name]",
	Short: "View logs from a pod",
	Long: `Stream logs from a Kubernetes pod.

If pod name is provided as argument, streams logs directly.
If not, presents an interactive selection.

Examples:
  kq logs nginx-7d8f9-xyz              # Logs from pod (interactive namespace)
  kq logs -n default nginx-7d8f9-xyz   # Logs from pod in specific namespace
  kq logs                              # Interactive: select namespace, pod`,

	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement logs command
		//
		// Step 1: Create Kubernetes client
		client, err := kube.NewClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating Kubernetes client: %v\n", err)
			os.Exit(1)
		}
		// Step 2: Determine namespace
		//   - Check --namespace flag first
		//   - Otherwise prompt for selection
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
		// Step 3: Determine pod
		//   - Check if pod name was passed as argument (args[0])
		//   - Otherwise prompt for selection
		var podName string
		if len(args) > 0 {
			podName = args[0]
		}
		if podName == "" {
			pods, err := kube.ListPods(client, namespace)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error listing pods in namespace %q: %v\n", namespace, err)
				os.Exit(1)
			}
			podName, err = ui.SelectOne("Select Pod", pods)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error selecting pod: %v\n", err)
				os.Exit(1)
			}
		}
		// Step 4: Get log options from flags
		//   - follow (-f)
		//   - timestamps
		//   - previous
		//   - container (if multi-container pod)
		follow, _ := cmd.Flags().GetBool("follow")
		timestamps, _ := cmd.Flags().GetBool("timestamps")
		previous, _ := cmd.Flags().GetBool("previous")
		container, _ := cmd.Flags().GetString("container")

		logOpts := kube.PodLogOptions{
			Follow:     follow,
			Timestamps: timestamps,
			Previous:   previous,
			Container:  container,
		}
		// Step 5: Check --dry-run flag
		//   - If set, print kubectl equivalent and exit
		//   - Example: kubectl logs -f -n default nginx-abc123
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if dryRun {
			kubectlCmd := "kubectl logs"
			if follow {
				kubectlCmd += " -f"
			}
			if timestamps {
				kubectlCmd += " --timestamps"
			}
			if previous {
				kubectlCmd += " -p"
			}
			if container != "" {
				kubectlCmd += fmt.Sprintf(" -c %s", container)
			}
			kubectlCmd += fmt.Sprintf(" -n %s %s", namespace, podName)

			fmt.Println("Dry run mode. Equivalent command:")
			fmt.Println(kubectlCmd)
			return
		}
		// Step 6: Stream logs
		//   - Use kube.GetPodLogs()
		//   - Pipe output to os.Stdout
		err = kube.GetPodLogs(client, namespace, podName, logOpts, os.Stdout)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error streaming logs from pod %q: %v\n", podName, err)
			os.Exit(1)
		}
	},
}

func init() {
	// Local flags for logs command
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	logsCmd.Flags().Bool("timestamps", false, "Show timestamps")
	logsCmd.Flags().BoolP("previous", "p", false, "Show logs from previous container instance")
	logsCmd.Flags().StringP("container", "c", "", "Container name (for multi-container pods)")

	// Register with root
	rootCmd.AddCommand(logsCmd)
}
