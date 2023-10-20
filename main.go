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
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filePath, content, 0644)
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
	startFilePath := "start.yaml"
	filteredFilePath := "filtered.yaml"
	forwardPortFilePath := "forward-port.yaml"

	data, err := readYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	forwardPortData := make(ReleaseData)

	fmt.Println("Charts with 'forward' or 'port' in the commit message:")
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
					fmt.Printf("Chart: %s Version: %s - Commit Message: %s\n", key, version, commitMsg)
					forwardPortData[key] = append(forwardPortData[key], version)
					delete(data, key)
				}
			}
		}
	}

	err = writeYaml(forwardPortFilePath, forwardPortData)
	if err != nil {
		log.Fatalf("Error writing to forward-port.yaml: %v", err)
	}

	startData := make(ReleaseData)
	for key, versions := range data {
		for _, version := range versions {
			startData[key] = append(startData[key], ":todo_added:  "+version)
		}
	}

	err = writeYaml(startFilePath, startData)
	if err != nil {
		log.Fatalf("Error writing to start.yaml: %v", err)
	}

	err = writeYaml(filteredFilePath, data)
	if err != nil {
		log.Fatalf("Error writing to filtered.yaml: %v", err)
	}
}
