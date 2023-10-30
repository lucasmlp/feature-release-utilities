package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type ReleaseData map[string][]string

var releasedFilePath = "released.yaml"

var releaseFilePath = "/Users/machado/development/suse/charts/release.yaml"

var assetsDir = "/Users/machado/development/suse/charts/assets"

var rootCmd = &cobra.Command{
	Use:   "chartcli",
	Short: "CLI tool for charts management",
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the filtered and updated chart files",
	Run: func(cmd *cobra.Command, args []string) {
		executeGeneration()
	},
}

var toBeReleasedCmd = &cobra.Command{
	Use:   "tobereleased",
	Short: "Generates a toBeReleased.yaml file with charts that have not been released.",
	Run:   executeToBeReleased,
}

func executeToBeReleased(cmd *cobra.Command, args []string) {
	releaseData, err := readYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	releasedData, err := readYaml(releasedFilePath)
	if err != nil {
		log.Fatalf("Error reading released.yaml: %v", err)
	}

	toBeReleasedData := make(map[string][]string)

	for chart, versions := range releaseData {
		if releasedVersions, ok := releasedData[chart]; ok {
			for _, v := range versions {
				if !contains(releasedVersions, v) {
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

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

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

func executeGeneration() {

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

func main() {
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(toBeReleasedCmd)
	toBeReleasedCmd.Flags().StringVarP(&releasedFilePath, "released", "r", "released.yaml", "Path to the released.yaml file.")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing CLI command: %v", err)
	}
}
