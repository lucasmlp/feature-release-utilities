package execute

import (
	"io/ioutil"
	"log"

	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func ToBeReleased(cmd *cobra.Command, args []string) {
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

	toBeReleasedBytes, err := yaml.Marshal(toBeReleasedData)
	if err != nil {
		log.Fatalf("Error marshaling toBeReleasedData to YAML: %v", err)
	}

	err = ioutil.WriteFile("toBeReleased.yaml", toBeReleasedBytes, 0644)
	if err != nil {
		log.Fatalf("Error writing toBeReleased.yaml: %v", err)
	}
}
