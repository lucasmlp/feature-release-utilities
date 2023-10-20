package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
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
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	var forwardPortedLines, notForwardPortedLines []string
	for chart, versions := range data {
		for _, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", chart, version)
			filePath := filepath.Join(assetsDir, chart, filename)
			commitMsg, err := checkLastCommitMessage(filePath)
			if err != nil {
				log.Printf("Error checking commit for %s: %v\n", filename, err)
				continue
			}
			commitWords := []string{"forward-port", "port", "forward", "port-forward"}
			isForwardPorted := false
			for _, word := range commitWords {
				if strings.Contains(commitMsg, word) {
					isForwardPorted = true
					break
				}
			}
			if isForwardPorted {
				forwardPortedLines = append(forwardPortedLines, fmt.Sprintf("%s %s: %s", chart, version, commitMsg))
			} else {
				notForwardPortedLines = append(notForwardPortedLines, fmt.Sprintf("%s %s: %s", chart, version, commitMsg))
			}
		}
	}

	ioutil.WriteFile("forward-port.log", []byte(strings.Join(forwardPortedLines, "\n")), 0644)
	ioutil.WriteFile("filtered.log", []byte(strings.Join(notForwardPortedLines, "\n")), 0644)
}
