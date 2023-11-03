package execute

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/lucasmlp/release-yaml-utils/pkg/git"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
)

func ExecuteGeneration() {

	data, err := utils.ReadYaml(utils.OriginalReleaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	var forwardPortedLines, notForwardPortedLines []string
	for chart, versions := range data {
		for _, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", chart, version)
			filePath := filepath.Join(utils.AssetsDir, chart, filename)
			commitMsg, err := git.CheckLastCommitMessage(filePath)
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
