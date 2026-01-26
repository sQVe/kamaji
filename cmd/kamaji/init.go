package main

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sqve/kamaji/internal/output"
)

// configTemplate is the default kamaji.yaml content with explanatory comments.
const configTemplate = `# Sprint name (required)
# A short identifier for this sprint
name: my-sprint

# Base branch to create ticket branches from
# Typically "main" or "develop"
base_branch: main

# Rules for the AI agent to follow during this sprint
# These guidelines help maintain code quality and consistency
rules:
  - Follow existing code patterns
  - Write tests for new functionality
  - Keep commits atomic and well-documented

# Tickets define the work items in this sprint
tickets:
  # Ticket name (required) - used as identifier
  - name: example-ticket

    # Git branch name for this ticket
    # Created from base_branch when the ticket starts
    branch: feature/example-ticket

    # What this ticket accomplishes
    description: Example ticket demonstrating the configuration format

    # Tasks break down the ticket into smaller units of work
    tasks:
      # Task description (required) - what this task does
      - description: Implement the example feature

        # Optional step-by-step guidance for the agent
        steps:
          - Create the necessary files
          - Write the implementation
          - Add tests

        # How to verify the task is complete
        verify: Run tests and confirm feature works as expected
`

func initCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create a new kamaji.yaml configuration file",
		RunE: func(_ *cobra.Command, _ []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return err
			}

			configPath := filepath.Join(workDir, configFile)

			_, err = os.Stat(configPath)
			if err == nil {
				output.PrintError(configFile + " already exists")
				return errFileExists
			}
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}

			if err := os.WriteFile(configPath, []byte(configTemplate), 0o600); err != nil {
				output.PrintError(err.Error())
				return errWriteFailed
			}

			output.PrintSuccess("Created " + configFile)
			return nil
		},
	}

	cmd.SilenceUsage = true

	return cmd
}
