// Package cmd implements the CLI commands for kq.
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time via ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:     "kq",
	Short:   "Kubernetes Quick Actions - interactive kubectl companion",
	Version: Version,
	Long: `kq is an interactive CLI for Kubernetes that reduces friction
between "I want to do something" and actually doing it.

Instead of memorizing kubectl flags and resource names, you fuzzy-find
your way to the right resource and action.`,
	Example: `  kq pods     # Interactively select and act on pods
  kq logs     # Quick access to pod logs`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Welcome to kq! Use --help to see available commands.")
		fmt.Println("Try: kq pods")
	},
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.SilenceUsage = true

	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Kubernetes namespace")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Print kubectl command instead of executing")

	rootCmd.AddCommand(podsCmd)
	rootCmd.AddCommand(logsCmd)
}
