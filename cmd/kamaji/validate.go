package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sqve/kamaji/internal/config"
	"github.com/sqve/kamaji/internal/output"
)

func validateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate kamaji.yaml configuration",
		RunE: func(cmd *cobra.Command, _ []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return err
			}

			sprintPath := filepath.Join(workDir, configFile)

			sprint, err := config.LoadSprint(sprintPath)
			if err != nil {
				output.PrintError(err.Error())
				return errConfigInvalid
			}

			validationErrors := config.ValidateSprint(sprint)
			if len(validationErrors) > 0 {
				output.PrintError("Configuration validation failed")
				for _, ve := range validationErrors {
					fmt.Fprintf(os.Stderr, "  %s: %s\n", ve.Field, ve.Message)
				}
				return errConfigInvalid
			}

			output.PrintSuccess("Configuration is valid")
			return nil
		},
	}

	cmd.SilenceUsage = true

	return cmd
}
