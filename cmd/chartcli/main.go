package main

import (
	"log"

	"github.com/lucasmlp/release-yaml-utils/cmd/chartcli/execute"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chartcli",
	Short: "CLI tool for charts management",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the filtered and updated chart files",
	Run:   execute.Generate,
}

var toBeReleasedCmd = &cobra.Command{
	Use:   "tobereleased",
	Short: "Generates a toBeReleased.yaml file with charts that have not been released.",
	Run:   execute.ToBeReleased,
}

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count the number of versions for charts in release.yaml, released.yaml and toBeReleased.yaml",
	Run:   execute.ExecuteCount,
}

var mergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge released.yaml and toBeReleased.yaml into release.yaml in the cwd",
	Run:   execute.ExecuteMerge,
}

func main() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(toBeReleasedCmd)
	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(mergeCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %v", err)
	}
}
