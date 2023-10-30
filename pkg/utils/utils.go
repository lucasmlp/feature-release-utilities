package utils

import (
	"io/ioutil"

	"github.com/lucasmlp/release-yaml-utils/pkg/models"
	"gopkg.in/yaml.v2"
)

func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func ReadYaml(filePath string) (models.ReleaseData, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var data models.ReleaseData
	if err := yaml.Unmarshal(content, &data); err != nil {
		return nil, err
	}

	return data, nil
}
