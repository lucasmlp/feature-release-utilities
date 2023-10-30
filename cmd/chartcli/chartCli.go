package main

import (
	"log"

	"github.com/lucasmlp/release-yaml-utils/pkg/release"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chartcli",
	Short: "CLI tool for charts management",
}

func main() {
	rootCmd.AddCommand(release.GenerateCmd)
	rootCmd.AddCommand(release.ToBeReleasedCmd)
	release.ToBeReleasedCmd.Flags().StringVarP(&release.ReleasedFilePath, "released", "r", "released.yaml", "Path to the released.yaml file.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %v", err)
	}
}
