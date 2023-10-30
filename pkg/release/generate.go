package release

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/lucasmlp/release-yaml-utils/pkg/git"
	"github.com/lucasmlp/release-yaml-utils/pkg/utils"
)

// ... other imports ...

var (
	releaseFilePath = "/Users/machado/development/suse/charts/release.yaml"
	assetsDir       = "/Users/machado/development/suse/charts/assets"
)

var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the filtered and updated chart files",
	Run: func(cmd *cobra.Command, args []string) {
		executeGeneration()
	},
}

func executeGeneration() {

	data, err := utils.ReadYaml(releaseFilePath)
	if err != nil {
		log.Fatalf("Error reading release.yaml: %v", err)
	}

	var forwardPortedLines, notForwardPortedLines []string
	for chart, versions := range data {
		for _, version := range versions {
			filename := fmt.Sprintf("%s-%s.tgz", chart, version)
			filePath := filepath.Join(assetsDir, chart, filename)
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
