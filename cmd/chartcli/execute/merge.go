package execute

import (
	"io/ioutil"
	"log"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ExecuteMerge(cmd *cobra.Command, args []string) {
	releasedData, err := utils.ReadYaml(utils.ReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading released.yaml: %v", err)
	}

	toBeReleasedData, err := utils.ReadYaml(utils.ToBeReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading toBeReleased.yaml: %v", err)
	}

	originalReleaseData, err := utils.ReadYaml(utils.ReleasedFilePath)
	if err != nil {
		log.Fatalf("Error reading the original release.yaml: %v", err)
	}

	mergedData := make(models.ReleaseData)

	for chart, versions := range releasedData {
		mergedData[chart] = versions
	}

	for chart, versions := range toBeReleasedData {
		for _, version := range versions {
			if !contains(mergedData[chart], version) {
				mergedData[chart] = append(mergedData[chart], version)
			}
		}
	}

	// Print differences
	for chart, versions := range mergedData {
		if _, exists := originalReleaseData[chart]; !exists {
			log.Printf("New chart added: %s with versions %v", chart, versions)
			continue
		}

		for _, version := range versions {
			if !contains(originalReleaseData[chart], version) {
				log.Printf("New version for chart %s: %s", chart, version)
			}
		}
	}

	mergedBytes, err := yaml.Marshal(mergedData)
	if err != nil {
		log.Fatalf("Error marshaling mergedData to YAML: %v", err)
	}

	err = ioutil.WriteFile("release.yaml", mergedBytes, 0644)
	if err != nil {
		log.Fatalf("Error writing release.yaml: %v", err)
	}
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
