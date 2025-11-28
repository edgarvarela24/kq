// cmd/pods.go - The 'pods' subcommand
//
// This is a SCAFFOLD. The actual implementation will:
// 1. List namespaces (let user select one)
// 2. List pods in that namespace (fuzzy search)
// 3. Show actions for the selected pod (logs, exec, describe, etc.)
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
	//
	// TODO (Phase 3): Replace this with actual implementation:
	// 1. Get k8s client (from internal/kube/client.go)
	// 2. List namespaces, let user select one
	// 3. List pods in namespace, fuzzy search
	// 4. Show action menu for selected pod
	// 5. Execute chosen action
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🚧 pods command - not yet implemented")
		fmt.Println("")
		fmt.Println("This will eventually:")
		fmt.Println("  1. Let you select a namespace")
		fmt.Println("  2. Fuzzy-search pods")
		fmt.Println("  3. Choose an action (logs, exec, describe, port-forward)")
		fmt.Println("")

		// DEBUGGER EXERCISE:
		// Set a breakpoint here. Inspect the `cmd` variable.
		// Look at cmd.Flags() - what's in there?
		// Look at cmd.Parent() - that's the root command!

		// TODO: Check if --dry-run flag is set
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		fmt.Printf("Dry-run mode: %v\n", dryRun)
		// TODO: Check if --namespace flag is set (to skip namespace selection)
		namespace, _ := cmd.Flags().GetString("namespace")
		fmt.Printf("Namespace flag: %s\n", namespace)
	},
}

// init registers any flags specific to the pods command
func init() {
	// === LOCAL FLAGS FOR PODS ===
	// These only apply to `kq pods`, not other commands
	//
	// We don't have any pods-specific flags yet, but you could add:
	// podsCmd.Flags().BoolP("all-namespaces", "A", false, "List pods in all namespaces")
}
