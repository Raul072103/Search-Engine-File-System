package models

import (
	"encoding/json"
	"os"
)

type FileTypesConfig struct {
	TypesMapping      map[int32]string `json:"types_mapping"`
	ExtensionMappings map[string]int32 `json:"extension_mappings"`
}

func (cfg *FileTypesConfig) GetTypeByExtension(extension string) string {
	typeValue, exists := cfg.ExtensionMappings[extension]
	if exists {
		return cfg.TypesMapping[typeValue]
	} else {
		return cfg.TypesMapping[-1]
	}
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
