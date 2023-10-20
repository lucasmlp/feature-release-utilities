package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

	fmt.Println("Checking assets for each chart:")
	for _, key := range keys {
		for _, version := range data[key] {
			filename := fmt.Sprintf("%s-%s.tgz", key, version)
			filePath := filepath.Join(assetsDir, key, filename)
			if fileExists(filePath) {
				fmt.Printf("Found: %s\n", filename)
			} else {
				fmt.Printf("Not found: %s\n", filename)
			}
		}
	}
}
