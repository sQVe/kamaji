package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/sqve/kamaji/internal/version"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "kamaji",
		Short:   "Kamaji orchestrates autonomous coding sprints",
		Version: version.Full(),
	}

	return cmd
}
