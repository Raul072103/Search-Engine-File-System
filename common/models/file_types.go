package utils

import (
	"encoding/json"
	"os"
)

type FileTypesConfig struct {
	TypesMapping      map[int32]string `json:"types_mapping"`
	ExtensionMappings map[int32]string `json:"extension_mappings"`
}

func ParseFileTypesConfig(filePath string) (FileTypesConfig, error) {
	var fileTypesConfig FileTypesConfig

	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fileTypesConfig, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&fileTypesConfig); err != nil {
		return fileTypesConfig, err
	}

	return fileTypesConfig, nil
}
