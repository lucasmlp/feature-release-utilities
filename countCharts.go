package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"strings"
)

type ReleaseData map[string][]string

func main() {
	localFilePath := "/Users/machado/development/suse/charts/release.yaml"
	assetsDir := "/Users/machado/development/suse/charts/assets"

	localData, err := readYaml(localFilePath)
	if err != nil {
		fmt.Printf("Error reading %s: %s\n", localFilePath, err)
		return
	}

	fmt.Printf("Number of items in %s: %d\n", localFilePath, len(localData))

	nonCRDCount := countNonCRDCharts(localData)
	fmt.Printf("Number of charts without -crd suffix: %d\n", nonCRDCount)

	checkCommitMessages(assetsDir, localData)
}

func readYaml(filePath string) (ReleaseData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make(ReleaseData)
	var key string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, ":") {
			key = strings.TrimSuffix(line, ":")
			data[key] = []string{}
		} else if len(line) > 2 {
			version := strings.TrimSpace(line[2:])
			data[key] = append(data[key], version)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return data, nil
}

func countNonCRDCharts(data ReleaseData) int {
	count := 0
	for key := range data {
		if !strings.HasSuffix(key, "-crd") {
			count++
		}
	}
	return count
}

func isPresentInReleaseData(fileName string, releaseData ReleaseData) bool {
	for chartName, versions := range releaseData {
		for _, version := range versions {
			expectedFileName := fmt.Sprintf("%s-%s.tgz", chartName, version)
			if strings.HasSuffix(fileName, expectedFileName) {
				return true
			}
		}
	}
	return false
}

func checkCommitMessages(directory string, releaseData ReleaseData) {
	fmt.Println("\nChecking commit messages...")

	err := fs.WalkDir(os.DirFS(directory), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		if isPresentInReleaseData(d.Name(), releaseData) {
			cmd := exec.Command("git", "log", "-1", "--oneline", "--", path)
			cmd.Dir = directory
			out, err := cmd.Output()
			if err != nil {
				fmt.Printf("Error while checking the git log for %s: %s\n", path, err)
				return nil
			}

			commitMessage := strings.TrimSpace(string(out))
			if strings.Contains(commitMessage, "forward") || strings.Contains(commitMessage, "port") {
				fmt.Printf("Chart affected: %s\nCommit Message: %s\n\n", d.Name(), commitMessage)
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error traversing the directory: %s\n", err)
	}
}
