package main

import (
	"log"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"github.com/lucasmlp/release-yaml-utils/pkg/release"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "chartcli",
	Short: "CLI tool for charts management",
}

func main() {
	rootCmd.AddCommand(release.GenerateCmd)
	rootCmd.AddCommand(release.ToBeReleasedCmd)
	rootCmd.AddCommand(countCmd) // Add this line
	release.ToBeReleasedCmd.Flags().StringVarP(&release.ReleasedFilePath, "released", "r", "released.yaml", "Path to the released.yaml file.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %v", err)
	}
}

var releaseFilePath = "/Users/machado/development/suse/charts/release.yaml"

var countCmd = &cobra.Command{
	Use:   "count",
	Short: "Count the number of versions for charts in release.yaml and toBeReleased.yaml",
	Run:   executeCount,
}

func executeCount(cmd *cobra.Command, args []string) {
	releaseData, err := utils.ReadYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	toBeReleasedData, err := utils.ReadYaml("toBeReleased.yaml")
	if err != nil {
		log.Fatalf("Error reading toBeReleased.yaml: %v", err)
	}

	releaseCount := countVersions(releaseData)
	toBeReleasedCount := countVersions(toBeReleasedData)

	log.Printf("Number of versions in release.yaml: %d\n", releaseCount)
	log.Printf("Number of versions in toBeReleased.yaml: %d\n", toBeReleasedCount)
}

func countVersions(data models.ReleaseData) int {
	count := 0
	for _, versions := range data {
		count += len(versions)
	}
	return count
}
