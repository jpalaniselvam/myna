package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myna",
	Short: "A Git-first CLI for executing and replaying AWS serverless actions",
	Long: `myna lets you invoke Lambda functions, send SQS messages, publish SNS events, 
emit EventBridge events, upload to S3, and run serverless workflows using simple, 
version-controlled TOML files.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
