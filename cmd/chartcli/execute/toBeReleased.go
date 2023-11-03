package execute

import (
	"fmt"
	"log"
	"time"

	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
)

func ToBeReleased(cmd *cobra.Command, args []string) {
	bkpDir := "bkp"
	timeSuffix := time.Now().Format("20060102-150405")

	// Back up the file
	bkpFile, err := backupFile(utils.ReleasedFilePath, bkpDir, timeSuffix)
	if err != nil {
		fmt.Printf("Failed to backup file %s: %v\n", utils.ReleasedFilePath, err)
		return
	}

	fmt.Printf("Backed up %s file to %s\n", utils.ReleasedFilePath, bkpFile)

	releaseData, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	releasedData, err := utils.ReadYaml(utils.ReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading released.yaml: %v", err)
	}

	toBeReleasedData := make(map[string][]string)

	for chart, versions := range releaseData {
		if releasedVersions, ok := releasedData[chart]; ok {
			for _, v := range versions {
				if !utils.Contains(releasedVersions, v) {
					toBeReleasedData[chart] = append(toBeReleasedData[chart], v)
				}
			}
		} else {
			toBeReleasedData[chart] = versions
		}
	}

	// Consolidate versions and count duplicates
	consolidatedData, duplicates := consolidateVersions(toBeReleasedData)

	// Print the number of duplicates
	fmt.Printf("File %s had %d duplicated chart/versions\n", utils.ToBeReleasedFilePath, duplicates)

	// Load release.yaml to determine the order
	releaseOrder, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		fmt.Printf("Failed to read release.yaml: %v\n", err)
		return
	}

	// Sort the data according to the release order
	sortReleaseData(&consolidatedData, &releaseOrder)

	// Write the sorted data back to the file
	err = utils.WriteYaml(utils.ToBeReleasedFilePath, consolidatedData)
	if err != nil {
		fmt.Printf("Failed to write sorted data to %s: %v\n", utils.ReleasedFilePath, err)
		return
	}

	fmt.Printf("Generated, consolidated and sorted %s\n", utils.ToBeReleasedFilePath)
}
