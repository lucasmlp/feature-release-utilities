package execute

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
)

func Format(cmd *cobra.Command, args []string) {
	files := []string{"toBeReleased.yaml", "released.yaml"}
	bkpDir := "bkp"
	timeSuffix := time.Now().Format("20060102-150405")

	// Create backup directory if it doesn't exist
	if _, err := os.Stat(bkpDir); os.IsNotExist(err) {
		os.Mkdir(bkpDir, os.ModePerm)
	}

	// Load release.yaml to determine the order
	releaseOrder, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		fmt.Printf("Failed to read release.yaml: %v\n", err)
		return
	}

	for _, file := range files {
		// Back up the file
		bkpFile, err := backupFile(file, bkpDir, timeSuffix)
		if err != nil {
			fmt.Printf("Failed to backup file %s: %v\n", file, err)
			continue
		}

		// Load the YAML file to be formatted
		data, err := utils.ReadYaml(file)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file, err)
			continue
		}

		// Consolidate versions and count duplicates
		consolidatedData, duplicates := consolidateVersions(data)

		// Print the number of duplicates
		fmt.Printf("File %s had %d duplicated chart/versions\n", file, duplicates)

		// Sort the data according to the release order
		sortReleaseData(&consolidatedData, &releaseOrder)

		// Write the sorted data back to the file
		err = utils.WriteYaml(file, consolidatedData)
		if err != nil {
			fmt.Printf("Failed to write sorted data to %s: %v\n", file, err)
			continue
		}

		fmt.Printf("Formatted and backed up %s to %s\n", file, bkpFile)
	}
}

func backupFile(file, bkpDir, timeSuffix string) (string, error) {
	// Extract the file extension and filename without extension
	ext := filepath.Ext(file)
	nameWithoutExt := strings.TrimSuffix(file, ext)

	// Construct the backup file path with the time suffix before the file extension
	bkpFile := filepath.Join(bkpDir, fmt.Sprintf("%s-%s%s", nameWithoutExt, timeSuffix, ext))

	// Perform the backup by copying
	err := utils.CopyFile(file, bkpFile)
	if err != nil {
		return "", err
	}

	return bkpFile, nil
}

func consolidateVersions(data models.ReleaseData) (models.ReleaseData, int) {
	consolidated := make(models.ReleaseData)
	duplicateCount := 0

	for chart, versions := range data {
		uniqueVersions := make(map[string]bool)
		var consolidatedVersions []string
		for _, version := range versions {
			if !uniqueVersions[version] {
				uniqueVersions[version] = true
				consolidatedVersions = append(consolidatedVersions, version)
			} else {
				duplicateCount++
			}
		}
		consolidated[chart] = consolidatedVersions
	}

	return consolidated, duplicateCount
}
