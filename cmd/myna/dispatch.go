package main

import (
	"fmt"

	"github.com/jpalaniselvam/myna/internal/action"
	"github.com/spf13/cobra"
)

var envFlag string

var dispatchCmd = &cobra.Command{
	Use:   "dispatch [file]",
	Short: "Execute a specific action or workflow",
	Long:  `Dispatch reads the provided TOML file (Action or Workflow) and executes it against the configured environment.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		file := args[0]
		if err := action.RunAction(file, envFlag); err != nil {
			fmt.Printf("Error executing action: %v\n", err)
		}
	},
}

func init() {
	dispatchCmd.Flags().StringVar(&envFlag, "env", "", "Environment to use (e.g., dev, prod)")
	rootCmd.AddCommand(dispatchCmd)
}
