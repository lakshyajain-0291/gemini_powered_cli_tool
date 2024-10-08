package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"gemini_cli_tool/fileinfo"
	"os"
	"path/filepath"
)

type ConfigData struct {
	Directories    []string `json:"directories"`
	SkipType       []string `json:"skip_types"`
	SkipFile       []string `json:"skip_files"`
	RelevanceIndex float32  `json:"relevance_index"`
	APIKeys        []string `json:"api_keys"`
}

func setConfig(addDirectories, deleteDirectories, addSkipTypes, deleteSkipTypes, addSkipFiles, deleteSkipFiles, addAPIKeys, deleteAPIKeys []string, relevanceIndex float32) error {

	config, err := LoadConfig()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if len(addDirectories) == 0 {
				cwd, err := os.Getwd()
				if err != nil {
					return fmt.Errorf("failed to get working directory: %w", err)
				}
				addDirectories = append(addDirectories, cwd)
			}
			config = &ConfigData{}
		} else {
			return err
		}
	}

	for _, dir := range addDirectories {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", dir)
		}
		if !contains(config.Directories, dir) {
			config.Directories = append(config.Directories, dir)
		}
	}

	for _, dir := range deleteDirectories {
		config.Directories = removeElements(config.Directories, dir)
	}

	if len(addSkipTypes) > 0 {
		config.SkipType = append(config.SkipType, addSkipTypes...)
	}

	for _, skipType := range deleteSkipTypes {
		config.SkipType = removeElements(config.SkipType, skipType)
	}

	if len(addSkipFiles) > 0 {
		config.SkipFile = append(config.SkipFile, addSkipFiles...)
	}

	for _, fileName := range deleteSkipFiles {
		config.SkipFile = removeElements(config.SkipFile, fileName)
	}

	for _, apiKey := range addAPIKeys {
		if !contains(config.APIKeys, apiKey) {
			config.APIKeys = append(config.APIKeys, apiKey)
		}
	}

	for _, api := range deleteAPIKeys {
		config.APIKeys = removeElements(config.APIKeys, api)
	}

	if relevanceIndex > 0 {
		config.RelevanceIndex = relevanceIndex
	}

	if err := SaveConfig(config); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println(fileinfo.Green("Configurations Saved Successfully"))
	return nil
}

// Helper function to check if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func removeElements(slice []string, element string) []string {
	for i, v := range slice {
		if v == element {
			slice = append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func LoadConfig() (*ConfigData, error) {

	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, ".gencli-config.json")
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config ConfigData
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}

func SaveConfig(config *ConfigData) error {

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, ".gencli-config.json")
	file, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(config)
}

func showConfigFormatted(config *ConfigData) {
	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println(fileinfo.Red(fmt.Sprintf("Failed to format config: %v\n", err)))
		return
	}
	fmt.Println(fileinfo.Cyan("----------------------------------------------------------------------------------------------------------------------------------\n"))
	fmt.Printf("%s\n%s\n", fileinfo.Green("\nCurrent Configuration :\n"), fileinfo.Gray(string(configBytes)))
	fmt.Println(fileinfo.Cyan("\n----------------------------------------------------------------------------------------------------------------------------------\n"))

}
