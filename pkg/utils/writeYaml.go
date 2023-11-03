package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

// WriteYaml takes data of any type, marshals it into YAML, and writes it to a file specified by filePath.
func WriteYaml(filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := yaml.NewEncoder(file)
	encoder.SetIndent(2) // Setting the indent to 2 spaces for YAML formatting.

	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return encoder.Close()
}
