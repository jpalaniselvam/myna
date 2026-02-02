package collection

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

/*
Summary of Test Cases:
1. CreateEnvironment:
   - Success: Verify creation of .toml file in environments directory.
   - MissingArgs: Error when collectionPath or envName is missing.
   - NoEnvDir: Error when environments directory doesn't exist.
   - AlreadyExists: Error when environment file already exists.

2. AddEnvVar:
   - Success: Verify adding a key-value pair.
   - MissingArgs: Error when args are missing.
   - FileMissing: Error when environment file doesn't exist.
   - KeyExists: Error when key to add already exists.

3. UpdateEnvVar:
   - Success: Verify updating a key-value pair.
   - KeyMissing: Error when key to update doesn't exist.
   - FileMissing: Error when environment file doesn't exist.

4. DeleteEnvVar:
   - Success: Verify deletion of a key.
   - KeyMissing: Error when key to delete doesn't exist.
   - FileMissing: Error when environment file doesn't exist.
   - MissingArgs: Error when args are missing.
*/

func TestCreateEnvironment(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")

	t.Run("Success", func(t *testing.T) {
		err := CreateEnvironment(collPath, "dev")
		if err != nil {
			t.Fatalf("CreateEnvironment failed: %v", err)
		}

		envPath := filepath.Join(collPath, "environments", "dev.toml")
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			t.Error("expected environment file to create")
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		if err := CreateEnvironment("", "dev"); err == nil {
			t.Error("expected error for missing collection path")
		}
		if err := CreateEnvironment(collPath, ""); err == nil {
			t.Error("expected error for missing env name")
		}
	})

	t.Run("NoEnvDir", func(t *testing.T) {
		// Create a separate dir structure without environments folder
		badPath := filepath.Join(tmpDir, "bad-col")
		os.Mkdir(badPath, 0755)
		if err := CreateEnvironment(badPath, "dev"); err == nil {
			t.Error("expected error when environments dir missing")
		}
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		// "dev" created in Success test
		if err := CreateEnvironment(collPath, "dev"); err == nil {
			t.Error("expected error when environment already exists")
		}
	})
}

func TestAddEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")
	CreateEnvironment(collPath, "dev")

	t.Run("Success", func(t *testing.T) {
		err := AddEnvVar(collPath, "dev", "API_KEY", "12345")
		if err != nil {
			t.Fatalf("AddEnvVar failed: %v", err)
		}

		// Verify content
		envPath := filepath.Join(collPath, "environments", "dev.toml")
		data, _ := os.ReadFile(envPath)
		var m map[string]interface{}
		toml.Unmarshal(data, &m)
		if m["API_KEY"] != "12345" {
			t.Errorf("expected API_KEY=12345, got %v", m["API_KEY"])
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		if err := AddEnvVar("", "dev", "k", "v"); err == nil {
			t.Error("expected error missing collection")
		}
		if err := AddEnvVar(collPath, "", "k", "v"); err == nil {
			t.Error("expected error missing env name")
		}
		if err := AddEnvVar(collPath, "dev", "", "v"); err == nil {
			t.Error("expected error missing key")
		}
	})

	t.Run("FileMissing", func(t *testing.T) {
		if err := AddEnvVar(collPath, "prod", "k", "v"); err == nil {
			t.Error("expected error for missing environment file")
		}
	})

	t.Run("KeyExists", func(t *testing.T) {
		// API_KEY added in Success
		if err := AddEnvVar(collPath, "dev", "API_KEY", "newval"); err == nil {
			t.Error("expected error when key already exists")
		}
	})
}

func TestUpdateEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")
	CreateEnvironment(collPath, "dev")
	AddEnvVar(collPath, "dev", "API_KEY", "original")

	t.Run("Success", func(t *testing.T) {
		err := UpdateEnvVar(collPath, "dev", "API_KEY", "updated")
		if err != nil {
			t.Fatalf("UpdateEnvVar failed: %v", err)
		}

		envPath := filepath.Join(collPath, "environments", "dev.toml")
		data, _ := os.ReadFile(envPath)
		var m map[string]interface{}
		toml.Unmarshal(data, &m)
		if m["API_KEY"] != "updated" {
			t.Errorf("expected API_KEY=updated, got %v", m["API_KEY"])
		}
	})

	t.Run("KeyMissing", func(t *testing.T) {
		if err := UpdateEnvVar(collPath, "dev", "NON_EXISTENT", "val"); err == nil {
			t.Error("expected error when key missing")
		}
	})

	t.Run("FileMissing", func(t *testing.T) {
		if err := UpdateEnvVar(collPath, "prod", "k", "v"); err == nil {
			t.Error("expected error for missing environment file")
		}
	})
}

func TestDeleteEnvVar(t *testing.T) {
	tmpDir := t.TempDir()
	collName := "test-col"
	collPath := filepath.Join(tmpDir, collName)
	Create(tmpDir, collName, "desc")
	CreateEnvironment(collPath, "dev")
	AddEnvVar(collPath, "dev", "API_KEY", "val")

	t.Run("Success", func(t *testing.T) {
		err := DeleteEnvVar(collPath, "dev", "API_KEY")
		if err != nil {
			t.Fatalf("DeleteEnvVar failed: %v", err)
		}

		envPath := filepath.Join(collPath, "environments", "dev.toml")
		data, _ := os.ReadFile(envPath)
		var m map[string]interface{}
		toml.Unmarshal(data, &m)
		if _, ok := m["API_KEY"]; ok {
			t.Error("expected API_KEY to be deleted")
		}
	})

	t.Run("KeyMissing", func(t *testing.T) {
		if err := DeleteEnvVar(collPath, "dev", "NON_EXISTENT"); err == nil {
			t.Error("expected error when key missing")
		}
	})

	t.Run("FileMissing", func(t *testing.T) {
		if err := DeleteEnvVar(collPath, "prod", "k"); err == nil {
			t.Error("expected error for missing environment file")
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		if err := DeleteEnvVar("", "dev", "k"); err == nil {
			t.Error("expected error missing collection")
		}
		if err := DeleteEnvVar(collPath, "", "k"); err == nil {
			t.Error("expected error missing env name")
		}
		if err := DeleteEnvVar(collPath, "dev", ""); err == nil {
			t.Error("expected error missing key")
		}
	})
}
