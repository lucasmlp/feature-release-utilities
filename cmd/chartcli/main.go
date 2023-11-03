package main

import (
	"log"

	"github.com/lucasmlp/release-yaml-utils/cmd/chartcli/execute"
	"github.com/lucasmlp/release-yaml-utils/pkg/release"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chartcli",
	Short: "CLI tool for charts management",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the filtered and updated chart files",
	Run: func(cmd *cobra.Command, args []string) {
		execute.ExecuteGeneration()
	},
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
	rootCmd.AddCommand(release.ToBeReleasedCmd)
	rootCmd.AddCommand(countCmd)
	rootCmd.AddCommand(mergeCmd)
	release.ToBeReleasedCmd.Flags().StringVarP(&release.ReleasedFilePath, "released", "r", "released.yaml", "Path to the released.yaml file.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %v", err)
	}
}
