// Package cmd implements the CLI commands for kq.
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/edgarvarela24/kq/internal/kube"
	"github.com/edgarvarela24/kq/internal/ui"
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs [pod-name]",
	Short: "View logs from a pod",
	Long: `Stream logs from a Kubernetes pod.

If pod name is provided as argument, streams logs directly.
If not, presents an interactive selection.`,
	Example: `  kq logs nginx-7d8f9-xyz              # Logs from pod (interactive namespace)
  kq logs -n default nginx-7d8f9-xyz   # Logs from pod in specific namespace
  kq logs                              # Interactive: select namespace, pod
  kq logs -f nginx-xyz                 # Follow logs`,
	Args: cobra.MaximumNArgs(1),

	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := kube.NewClient()
		if err != nil {
			return fmt.Errorf("creating kubernetes client: %w", err)
		}

		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer cancel()

		namespace := cmd.Flag("namespace").Value.String()
		if namespace == "" {
			namespaces, err := kube.ListNamespaces(ctx, client)
			if err != nil {
				return fmt.Errorf("listing namespaces: %w", err)
			}
			namespace, err = ui.SelectOne("Select Namespace", namespaces)
			if err != nil {
				return fmt.Errorf("selecting namespace: %w", err)
			}
		}

		var podName string
		if len(args) > 0 {
			podName = args[0]
		}
		if podName == "" {
			pods, err := kube.ListPods(ctx, client, namespace)
			if err != nil {
				return fmt.Errorf("listing pods in namespace %q: %w", namespace, err)
			}
			podName, err = ui.SelectOne("Select Pod", pods)
			if err != nil {
				return fmt.Errorf("selecting pod: %w", err)
			}
		}

		// Check if any log option flags were explicitly set
		followChanged := cmd.Flags().Changed("follow")
		timestampsChanged := cmd.Flags().Changed("timestamps")
		previousChanged := cmd.Flags().Changed("previous")
		containerChanged := cmd.Flags().Changed("container")

		var follow, timestamps, previous bool
		var container string

		// If no flags were set, prompt interactively for options
		if !followChanged && !timestampsChanged && !previousChanged && !containerChanged {
			// Handle container selection for multi-container pods
			containers, err := kube.ListContainers(ctx, client, namespace, podName)
			if err != nil {
				return fmt.Errorf("listing containers in pod %q: %w", podName, err)
			}
			if len(containers) == 1 {
				container = containers[0]
			} else {
				container, err = ui.SelectOne("Select Container", containers)
				if err != nil {
					return fmt.Errorf("selecting container: %w", err)
				}
			}

			follow, timestamps, previous, err = ui.SelectLogOptions()
			if err != nil {
				return fmt.Errorf("selecting log options: %w", err)
			}
		} else {
			follow, _ = cmd.Flags().GetBool("follow")
			timestamps, _ = cmd.Flags().GetBool("timestamps")
			previous, _ = cmd.Flags().GetBool("previous")
			container, _ = cmd.Flags().GetString("container")
		}

		logOpts := kube.PodLogOptions{
			Follow:     follow,
			Timestamps: timestamps,
			Previous:   previous,
			Container:  container,
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")
		if dryRun {
			fmt.Println(buildLogsCommand(namespace, podName, logOpts))
			return nil
		}

		if err := kube.GetPodLogs(ctx, client, namespace, podName, logOpts, os.Stdout); err != nil {
			return fmt.Errorf("streaming logs from pod %q: %w", podName, err)
		}
		return nil
	},
}

func init() {
	logsCmd.Flags().BoolP("follow", "f", false, "Follow log output")
	logsCmd.Flags().Bool("timestamps", false, "Show timestamps")
	logsCmd.Flags().BoolP("previous", "p", false, "Show logs from previous container instance")
	logsCmd.Flags().StringP("container", "c", "", "Container name (for multi-container pods)")
}
