package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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

func writeYaml(filePath string, data ReleaseData, prefix string) error {
	var content string
	for chart, versions := range data {
		content += prefix + chart + ":\n"
		for _, version := range versions {
			content += prefix + "  - " + version + "\n"
		}
	}

	return ioutil.WriteFile(filePath, []byte(content), 0644)
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

	// Read release.yaml
	data, err := readYaml(releaseFilePath)
	if err != nil {
		panic(fmt.Sprintf("Error reading YAML file: %v", err))
	}

	// Filter data based on commit messages
	filteredData := make(ReleaseData)
	for chart, versions := range data {
		for _, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", chart, version)
			filePath := filepath.Join(assetsDir, chart, filename)
			if fileExists(filePath) {
				commitMsg, err := checkLastCommitMessage(filePath)
				if err != nil {
					continue
				}
				if !strings.Contains(commitMsg, "forward") && !strings.Contains(commitMsg, "port") {
					if _, ok := filteredData[chart]; !ok {
						filteredData[chart] = []string{}
					}
					filteredData[chart] = append(filteredData[chart], version)
				}
			}
		}
	}

	// Write filtered.yaml without prefix
	err = writeYaml("filtered.yaml", filteredData, "")
	if err != nil {
		panic(fmt.Sprintf("Error writing filtered.yaml: %v", err))
	}

	// Write start.yaml with :todo_added: prefix and two spaces
	err = writeYaml("start.yaml", filteredData, ":todo_added:  ")
	if err != nil {
		panic(fmt.Sprintf("Error writing start.yaml: %v", err))
	}
}
