// cmd/root.go - Root command definition using Cobra
//
// Cobra concepts you'll learn here:
// - Commands: The actions your CLI can perform (kq, kq pods, kq logs, etc.)
// - Flags: Options that modify behavior (--namespace, --follow)
// - PersistentFlags: Flags inherited by all subcommands
//
// Cobra docs: https://cobra.dev/
// Cobra GitHub: https://github.com/spf13/cobra
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
//
// When you run just `kq`, this is what executes.
// Subcommands (like `kq pods`) are added to this via AddCommand().
var rootCmd = &cobra.Command{
	// Use is the one-line usage message
	// The first word is the command name
	Use: "kq",

	// Short is a brief description shown in help
	Short: "Kubernetes Quick Actions - interactive kubectl companion",

	// Long is a longer description shown in 'kq --help'
	Long: `kq is an interactive CLI for Kubernetes that reduces friction
between "I want to do something" and actually doing it.

Instead of memorizing kubectl flags and resource names, you fuzzy-find
your way to the right resource and action.

Example:
  kq pods     # Interactively select and act on pods
  kq logs     # Quick access to pod logs
  kq exec     # Interactive shell into a pod`,

	// Run is executed when the command is called directly (no subcommand)
	//
	// TODO: For now, we just print a message. Later, this could launch
	// an interactive menu to choose what resource type to work with.
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to kq! Use --help to see available commands.")
		fmt.Println("Try: kq pods")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// init() is a special Go function that runs automatically when the package is imported.
// This is the standard Cobra pattern for setting up flags and adding subcommands.
//
// Learn more about Go's init(): https://go.dev/doc/effective_go#init
func init() {
	// === PERSISTENT FLAGS ===
	// These flags will be available to this command AND all subcommands.
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Kubernetes namespace")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Print kubectl command instead of executing")

	// === LOCAL FLAGS ===
	// These flags only apply to this command, not subcommands.

	// === ADD SUBCOMMANDS ===
	// Each subcommand is added here. As you build more, add them.
	rootCmd.AddCommand(podsCmd)
}
