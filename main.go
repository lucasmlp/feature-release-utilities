package main

import (
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"
)

type ReleaseData map[string][]string

func readYaml(filePath string) (ReleaseData, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data ReleaseData
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func generatePrefixedYaml(data ReleaseData, qaData ReleaseData) string {
	var lines []string
	for chart, versions := range data {
		prefix := ":todo_added:\t"
		if _, exists := qaData[chart]; exists {
			prefix = ":todo_done:\t"
		}
		lines = append(lines, prefix+chart+":")
		for _, version := range versions {
			lines = append(lines, prefix+"- "+version)
		}
	}
	return strings.Join(lines, "\n")
}

func main() {
	releaseFilePath := "/Users/machado/development/suse/charts/release.yaml"
	qaSignoffFilePath := "qa-signoff.yaml"
	updatedFilePath := "updated.yaml"

	data, err := readYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	qaData, err := readYaml(qaSignoffFilePath)
	if err != nil {
		log.Fatalf("Error reading qa-signoff.yaml: %v", err)
	}

	updatedContent := generatePrefixedYaml(data, qaData)
	err = ioutil.WriteFile(updatedFilePath, []byte(updatedContent), 0644)
	if err != nil {
		log.Fatalf("Error writing to updated.yaml: %v", err)
	}
}
