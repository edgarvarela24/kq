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
	"k8s.io/client-go/kubernetes"
)

var podActions = []string{
	"logs",
}

var podsCmd = &cobra.Command{
	Use:   "pods",
	Short: "Interactively select and act on pods",
	Long: `List pods with fuzzy search and perform actions on them.

Actions available:
  - logs: View pod logs (with follow, timestamps options)
  - exec: Open an interactive shell (coming soon)
  - describe: Show detailed pod information (coming soon)
  - port-forward: Forward local port to pod (coming soon)`,

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

		pods, err := kube.ListPods(ctx, client, namespace)
		if err != nil {
			return fmt.Errorf("listing pods in namespace %q: %w", namespace, err)
		}

		if len(pods) == 0 {
			return fmt.Errorf("no pods found in namespace %q", namespace)
		}

		pod, err := ui.SelectOne("Select Pod", pods)
		if err != nil {
			return fmt.Errorf("selecting pod: %w", err)
		}

		action, err := ui.SelectOne("Select Action", podActions)
		if err != nil {
			return fmt.Errorf("selecting action: %w", err)
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		switch action {
		case "logs":
			return executeLogsAction(ctx, client, namespace, pod, dryRun)
		default:
			return fmt.Errorf("action %q not implemented", action)
		}
	},
}

func executeLogsAction(ctx context.Context, client kubernetes.Interface, namespace, pod string, dryRun bool) error {
	containers, err := kube.ListContainers(ctx, client, namespace, pod)
	if err != nil {
		return fmt.Errorf("listing containers in pod %q: %w", pod, err)
	}

	var container string
	if len(containers) == 1 {
		container = containers[0]
	} else {
		container, err = ui.SelectOne("Select Container", containers)
		if err != nil {
			return fmt.Errorf("selecting container: %w", err)
		}
	}

	follow, timestamps, previous, err := ui.SelectLogOptions()
	if err != nil {
		return fmt.Errorf("selecting log options: %w", err)
	}

	logOpts := kube.PodLogOptions{
		Follow:     follow,
		Timestamps: timestamps,
		Previous:   previous,
		Container:  container,
	}

	if dryRun {
		fmt.Println(buildLogsCommand(namespace, pod, logOpts))
		return nil
	}

	if err := kube.GetPodLogs(ctx, client, namespace, pod, logOpts, os.Stdout); err != nil {
		return fmt.Errorf("streaming logs from pod %q: %w", pod, err)
	}
	return nil
}

// buildLogsCommand constructs the equivalent kubectl logs command.
func buildLogsCommand(namespace, pod string, opts kube.PodLogOptions) string {
	cmd := "kubectl logs"
	if opts.Follow {
		cmd += " -f"
	}
	if opts.Timestamps {
		cmd += " --timestamps"
	}
	if opts.Previous {
		cmd += " -p"
	}
	if opts.Container != "" {
		cmd += fmt.Sprintf(" -c %s", opts.Container)
	}
	cmd += fmt.Sprintf(" -n %s %s", namespace, pod)
	return cmd
}
