package collection

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

// CreateEnvironment creates a new environment configuration file in the collection
func CreateEnvironment(collectionPath, envName string) error {
	if collectionPath == "" || envName == "" {
		return fmt.Errorf("collection path and environment name are required")
	}

	envDir := filepath.Join(collectionPath, "environments")
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return fmt.Errorf("environments directory does not exist in %s", collectionPath)
	}

	envFilePath := filepath.Join(envDir, envName+".toml")
	if _, err := os.Stat(envFilePath); !os.IsNotExist(err) {
		return fmt.Errorf("environment %s already exists", envName)
	}

	// Create empty TOML file
	f, err := os.Create(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to create environment file: %w", err)
	}
	f.Close()
	return nil
}

// GetEnvironment retrieves the variables of an environment
func GetEnvironment(collectionPath, envName string) (map[string]interface{}, error) {
	if collectionPath == "" || envName == "" {
		return nil, fmt.Errorf("collection path and environment name are required")
	}

	envFilePath := filepath.Join(collectionPath, "environments", envName+".toml")
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %w", err)
	}

	var envMap map[string]interface{}
	if len(data) > 0 {
		if err := toml.Unmarshal(data, &envMap); err != nil {
			return nil, fmt.Errorf("failed to parse environment file: %w", err)
		}
	} else {
		envMap = make(map[string]interface{})
	}

	return envMap, nil
}

// AddEnvVar adds a key-value pair to an environment. Errors if key already exists.
func AddEnvVar(collectionPath, envName, key, value string) error {
	return modifyEnvVar(collectionPath, envName, key, value, false)
}

// UpdateEnvVar updates a key-value pair in an environment. Errors if key does not exist.
func UpdateEnvVar(collectionPath, envName, key, value string) error {
	return modifyEnvVar(collectionPath, envName, key, value, true)
}

// DeleteEnvVar deletes a key-value pair from an environment
func DeleteEnvVar(collectionPath, envName, key string) error {
	if collectionPath == "" || envName == "" || key == "" {
		return fmt.Errorf("collection path, environment name, and key are required")
	}

	envFilePath := filepath.Join(collectionPath, "environments", envName+".toml")
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		return fmt.Errorf("failed to read environment file: %w", err)
	}

	var envMap map[string]interface{}
	if err := toml.Unmarshal(data, &envMap); err != nil {
		return fmt.Errorf("failed to parse environment file: %w", err)
	}

	if _, exists := envMap[key]; !exists {
		return fmt.Errorf("key %s does not exist in environment %s", key, envName)
	}

	delete(envMap, key)

	newData, err := toml.Marshal(envMap)
	if err != nil {
		return fmt.Errorf("failed to marshal environment updates: %w", err)
	}

	if err := os.WriteFile(envFilePath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	return nil
}

func modifyEnvVar(collectionPath, envName, key, value string, isUpdate bool) error {
	if collectionPath == "" || envName == "" || key == "" {
		return fmt.Errorf("collection path, environment name, and key are required")
	}

	envFilePath := filepath.Join(collectionPath, "environments", envName+".toml")
	data, err := os.ReadFile(envFilePath)
	if err != nil {
		// If file doesn't exist, we can't add/update var
		return fmt.Errorf("failed to read environment file: %w", err)
	}

	var envMap map[string]interface{}
	if len(data) > 0 {
		if err := toml.Unmarshal(data, &envMap); err != nil {
			return fmt.Errorf("failed to parse environment file: %w", err)
		}
	} else {
		envMap = make(map[string]interface{})
	}

	_, exists := envMap[key]
	if isUpdate {
		if !exists {
			return fmt.Errorf("key %s does not exist in environment %s (use add to create)", key, envName)
		}
	} else {
		if exists {
			return fmt.Errorf("key %s already exists in environment %s (use update to modify)", key, envName)
		}
	}

	envMap[key] = value

	newData, err := toml.Marshal(envMap)
	if err != nil {
		return fmt.Errorf("failed to marshal environment updates: %w", err)
	}

	if err := os.WriteFile(envFilePath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write environment file: %w", err)
	}

	return nil
}
