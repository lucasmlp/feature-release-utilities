package git

import (
	"bytes"
	"os/exec"
	"path/filepath"
	// ... other imports ...
)

func CheckLastCommitMessage(filePath string) (string, error) {
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
