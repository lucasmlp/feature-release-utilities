// execute/merge.go
package execute

import (
	"log"
	"sort"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
)

// ExecuteMerge is called by the cobra command to merge released and toBeReleased YAML files.
func ExecuteMerge(cmd *cobra.Command, args []string) {
	// Read the original release data
	originalReleaseData, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading original release.yaml: %v", err)
	}

	// Read the released data
	releasedData, err := utils.ReadYaml(utils.ReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading released.yaml: %v", err)
	}

	// Read the toBeReleased data
	toBeReleasedData, err := utils.ReadYaml(utils.ToBeReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading toBeReleased.yaml: %v", err)
	}

	// Merge releasedData and toBeReleasedData
	mergedData := mergeReleaseData(releasedData, toBeReleasedData)

	// Print the differences
	printDifferences(originalReleaseData, mergedData)

	// Sort and write the merged data back to the release.yaml file
	// Assuming releaseData and sortOrderData are the maps you are working with.
	sortReleaseData(&mergedData, &originalReleaseData)

	utils.WriteYaml(utils.MergedReleaseFilePath, mergedData)
}

// mergeReleaseData merges the new data into the existing release data.
func mergeReleaseData(releasedData, toBeReleasedData models.ReleaseData) models.ReleaseData {
	mergedData := models.ReleaseData{}

	// Add all releasedData to mergedData
	for chart, versions := range releasedData {
		if _, exists := mergedData[chart]; !exists {
			mergedData[chart] = make([]string, 0)
		}
		mergedData[chart] = append(mergedData[chart], versions...)
	}

	// Add non-duplicate toBeReleasedData to mergedData
	for chart, versions := range toBeReleasedData {
		if _, exists := mergedData[chart]; !exists {
			mergedData[chart] = make([]string, 0)
		}
		for _, version := range versions {
			if !utils.Contains(mergedData[chart], version) {
				mergedData[chart] = append(mergedData[chart], version)
			}
		}
	}

	// Potentially you might want to sort the mergedData for consistency
	// This depends on your specific use case and how you want to handle version ordering.

	return mergedData
}

// PrintDifferences will print the differences between two ReleaseData maps.
func printDifferences(originalData, newData models.ReleaseData) {
	log.Println("Differences (present in merged release.yaml but not in original one):")
	for chart, versions := range newData {
		originalVersions, ok := originalData[chart]
		if !ok {
			// If the chart does not exist at all in the original data, print all versions
			log.Printf("New chart added: %s with versions %v\n", chart, versions)
			continue
		}

		for _, version := range versions {
			if !utils.Contains(originalVersions, version) {
				log.Printf("New version for chart %s: %s\n", chart, version)
			}
		}
	}

	log.Println("Differences (present in original release.yaml but not in merged one):")
	for chart, versions := range originalData {
		newVersions, ok := newData[chart]
		if !ok {
			// If the chart does not exist in the new data, print all versions from the original
			log.Printf("Chart removed: %s with versions %v\n", chart, versions)
			continue
		}

		for _, version := range versions {
			if !utils.Contains(newVersions, version) {
				log.Printf("Version removed for chart %s: %s\n", chart, version)
			}
		}
	}
}

// SortReleaseData sorts the data in releaseData to match the order of charts and versions in sortOrderData.
func sortReleaseData(releaseData, sortOrderData *models.ReleaseData) {
	for chart, versions := range *releaseData {
		if sortOrder, ok := (*sortOrderData)[chart]; ok {
			sort.SliceStable(versions, func(i, j int) bool {
				// Find the index of the versions i and j in the sortOrder slice.
				// If a version is not found, it is treated as if it has a higher sort order (placed last).
				indexI := utils.IndexOf(sortOrder, versions[i])
				indexJ := utils.IndexOf(sortOrder, versions[j])
				return indexI < indexJ
			})
		} else {
			// If the chart is not found in sortOrderData, just sort it alphabetically
			sort.Strings(versions)
		}
	}
}
