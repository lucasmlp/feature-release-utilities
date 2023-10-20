package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
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

func writeYaml(filePath string, data ReleaseData) error {
	encodedData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, encodedData, 0644)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func checkLastCommitMessage(filePath string) (string, error) {
	cmd := exec.Command("git", "log", "-n", "1", "--pretty=format:%s", filePath)
	cmd.Dir = filepath.Dir(filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func main() {
	releaseFilePath := "/Users/machado/development/suse/charts/release.yaml"
	assetsDir := "/Users/machado/development/suse/charts/assets"

	data, err := readYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	affectedCharts := make(map[string][]string)

	for _, key := range keys {
		for _, version := range data[key] {
			filename := fmt.Sprintf("%s-%s.tgz", key, version)
			filePath := filepath.Join(assetsDir, key, filename)
			if fileExists(filePath) {
				commitMsg, err := checkLastCommitMessage(filePath)
				if err != nil {
					fmt.Printf("Error checking commit for %s: %v\n", filename, err)
					continue
				}
				if strings.Contains(commitMsg, "forward") || strings.Contains(commitMsg, "port") {
					affectedCharts[key] = append(affectedCharts[key], version)
				}
			} else {
				fmt.Printf("Not found: %s\n", filename)
			}
		}
	}

	// Create a filtered release data by excluding affected charts
	filteredReleaseData := make(ReleaseData)
	for key, versions := range data {
		if _, exists := affectedCharts[key]; !exists {
			filteredReleaseData[key] = versions
		}
	}

	// Write filtered data to release-filtered.yaml in the current working directory
	err = writeYaml("release-filtered.yaml", filteredReleaseData)
	if err != nil {
		log.Fatalf("Error writing to release-filtered.yaml: %v", err)
	}

	fmt.Println("release-filtered.yaml has been generated successfully.")
}
