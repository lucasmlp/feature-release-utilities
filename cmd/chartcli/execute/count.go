package execute

import (
	"log"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
)

func ExecuteCount(cmd *cobra.Command, args []string) {
	releaseData, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	releasedData, err := utils.ReadYaml(utils.ReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading released.yaml: %v", err)
	}

	toBeReleasedData, err := utils.ReadYaml(utils.ToBeReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading toBeReleased.yaml: %v", err)
	}

	mergedReleaseData, err := utils.ReadYaml(utils.MergedReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading merged release.yaml: %v", err)
	}
	releaseCount := countVersions(releaseData)
	releasedCount := countVersions(releasedData)
	toBeReleasedCount := countVersions(toBeReleasedData)
	mergedReleaseCount := countVersions(mergedReleaseData)

	log.Printf("Number of versions in release.yaml: %d\n", releaseCount)
	log.Printf("Number of versions in released.yaml: %d\n", releasedCount)
	log.Printf("Number of versions in toBeReleased.yaml: %d\n", toBeReleasedCount)
	log.Printf("Number of versions in merged release.yaml: %d\n", mergedReleaseCount)
}

func countVersions(data models.ReleaseData) int {
	count := 0
	for _, versions := range data {
		count += len(versions)
	}
	return count
}
