package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sqve/kamaji/internal/version"
)

func main() {
	if err := rootCmd().Execute(); err != nil {
		if !errors.Is(err, errSprintFailed) && !errors.Is(err, errConfigInvalid) {
			fmt.Fprintln(os.Stderr, err)
		}
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "kamaji",
		Short:         "Kamaji orchestrates autonomous coding sprints",
		Version:       version.Full(),
		SilenceErrors: true,
	}

	cmd.AddCommand(startCmd())
	cmd.AddCommand(validateCmd())

	return cmd
}
