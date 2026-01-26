package main

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/sqve/kamaji/internal/orchestrator"
)

func startCmd() *cobra.Command {
	var spawnerCmd string

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Run sprint until complete or stuck",
		Long:  "Execute tasks sequentially from kamaji.yaml until the sprint completes or a task fails 3 consecutive times.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			workDir, err := os.Getwd()
			if err != nil {
				return err
			}

			result, err := orchestrator.Run(cmd.Context(), orchestrator.RunConfig{
				WorkDir:    workDir,
				SprintPath: filepath.Join(workDir, configFile),
				SpawnerCmd: spawnerCmd,
			})
			if err != nil {
				return err
			}

			if !result.Success {
				return errSprintFailed
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&spawnerCmd, "spawner-cmd", "", "Override spawner command (for testing)")
	_ = cmd.Flags().MarkHidden("spawner-cmd")

	cmd.SilenceUsage = true

	return cmd
}
