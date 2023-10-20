package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

func writePrefixedYaml(filePath string, data ReleaseData) error {
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for i := range lines {
		if lines[i] != "" {
			lines[i] = ":todo_added:\t" + lines[i]
		}
	}
	prefixedContent := strings.Join(lines, "\n")

	return ioutil.WriteFile(filePath, []byte(prefixedContent), 0644)
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

func hasCommitWords(commitMsg string) bool {
	words := []string{"\\bforward-port\\b", "\\bport\\b", "\\bforward\\b", "\\bport-forward\\b"}
	for _, word := range words {
		if regexp.MustCompile(word).MatchString(commitMsg) {
			return true
		}
	}
	return false
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

	for _, key := range keys {
		for _, version := range data[key] {
			filename := fmt.Sprintf("%s-%s.tgz", key, version)
			filePath := filepath.Join(assetsDir, key, filename)
			if fileExists(filePath) {
				commitMsg, err := checkLastCommitMessage(filePath)
				if err != nil {
					continue
				}
				if hasCommitWords(commitMsg) {
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

	err = writePrefixedYaml(startFilePath, data)
	if err != nil {
		log.Fatalf("Error writing to start.yaml: %v", err)
	}

	err = writeYaml(filteredFilePath, data)
	if err != nil {
		log.Fatalf("Error writing to filtered.yaml: %v", err)
	}
}
