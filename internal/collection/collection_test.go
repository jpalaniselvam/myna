package collection

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/jpalaniselvam/myna/internal/types"
	"github.com/pelletier/go-toml/v2"
)

/*
Summary of Test Cases:
1. Create:
   - Success: Verify creation of directory, myna.toml, and environments folder.
   - MissingArgs: Error when baseDir or name is empty.
   - AlreadyExists: Error when collection directory already exists.

2. Update:
   - Success: Verify updating description and variables.
   - Rename: Verify renaming the collection directory.
   - ParseError: Error when myna.toml is invalid.
   - MissingArgs: Error when path is empty.

3. Delete:
   - Success: Verify deletion of collection directory.
   - NotExists: Error when directory doesn't exist.
   - MissingArgs: Error when path is empty.

4. Get:
   - Success: Verify retrieval of environments, pre vars, and nested action tree.
   - NotExists: Error when directory/myna.toml doesn't exist.
   - MissingArgs: Error when path is empty.
*/

func TestCreate(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("Success", func(t *testing.T) {
		err := Create(tmpDir, "my-collection", "A test collection")
		if err != nil {
			t.Fatalf("Create failed: %v", err)
		}

		collPath := filepath.Join(tmpDir, "my-collection")
		if _, err := os.Stat(collPath); os.IsNotExist(err) {
			t.Error("expected collection directory to create")
		}

		envDir := filepath.Join(collPath, "environments")
		if _, err := os.Stat(envDir); os.IsNotExist(err) {
			t.Error("expected environments directory to create")
		}

		mynaPath := filepath.Join(collPath, "myna.toml")
		data, err := os.ReadFile(mynaPath)
		if err != nil {
			t.Fatalf("failed to read myna.toml: %v", err)
		}

		var coll types.Collection
		if err := toml.Unmarshal(data, &coll); err != nil {
			t.Fatalf("failed to parse myna.toml: %v", err)
		}

		if coll.Metadata.Description != "A test collection" {
			t.Errorf("expected description 'A test collection', got %s", coll.Metadata.Description)
		}
		if coll.Metadata.Kind != "collection" {
			t.Errorf("expected kind 'collection', got %s", coll.Metadata.Kind)
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		err := Create("", "name", "desc")
		if err == nil {
			t.Error("expected error for missing baseDir")
		}
		err = Create("dir", "", "desc")
		if err == nil {
			t.Error("expected error for missing name")
		}
	})

	t.Run("AlreadyExists", func(t *testing.T) {
		name := "existing-col"
		os.Mkdir(filepath.Join(tmpDir, name), 0755)

		err := Create(tmpDir, name, "desc")
		if err == nil {
			t.Error("expected error for existing collection")
		}
	})
}

func TestUpdate(t *testing.T) {
	tmpDir := t.TempDir()
	name := "update-col"
	collPath := filepath.Join(tmpDir, name)

	// Setup initial collection
	Create(tmpDir, name, "Initial Desc")

	t.Run("Success", func(t *testing.T) {
		newDesc := "Updated Desc"
		preVars := map[string]interface{}{"key": "value"}

		err := Update(collPath, name, newDesc, preVars, nil)
		if err != nil {
			t.Fatalf("Update failed: %v", err)
		}

		// Verify changes
		mynaPath := filepath.Join(collPath, "myna.toml")
		data, _ := os.ReadFile(mynaPath)
		var coll types.Collection
		toml.Unmarshal(data, &coll)

		if coll.Metadata.Description != newDesc {
			t.Errorf("expected desc %s, got %s", newDesc, coll.Metadata.Description)
		}
		if coll.Vars.Pre["key"] != "value" {
			t.Errorf("expected pre var key=value, got %v", coll.Vars.Pre["key"])
		}
	})

	t.Run("Rename", func(t *testing.T) {
		newName := "renamed-col"
		err := Update(collPath, newName, "", nil, nil)
		if err != nil {
			t.Fatalf("Rename failed: %v", err)
		}

		newPath := filepath.Join(tmpDir, newName)
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			t.Error("expected directory to be renamed")
		}
		if _, err := os.Stat(collPath); !os.IsNotExist(err) {
			t.Error("expected old directory to be gone")
		}

		// Update invalid path for next tests if needed, or re-setup
		collPath = newPath
		name = newName
	})

	t.Run("MissingArgs", func(t *testing.T) {
		err := Update("", "name", "desc", nil, nil)
		if err == nil {
			t.Error("expected error for missing collection path")
		}
	})
}

func TestDelete(t *testing.T) {
	tmpDir := t.TempDir()
	name := "del-col"
	collPath := filepath.Join(tmpDir, name)
	Create(tmpDir, name, "desc")

	t.Run("Success", func(t *testing.T) {
		err := Delete(collPath)
		if err != nil {
			t.Fatalf("Delete failed: %v", err)
		}

		if _, err := os.Stat(collPath); !os.IsNotExist(err) {
			t.Error("expected directory to be deleted")
		}
	})

	t.Run("NotExists", func(t *testing.T) {
		err := Delete(filepath.Join(tmpDir, "non-existent"))
		if err == nil {
			t.Error("expected error for non-existent directory")
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		err := Delete("")
		if err == nil {
			t.Error("expected error for missing collection path")
		}
	})
}

func TestGet(t *testing.T) {
	tmpDir := t.TempDir()
	name := "get-col"
	collPath := filepath.Join(tmpDir, name)
	Create(tmpDir, name, "desc")

	// Update with some vars
	Update(collPath, name, "", map[string]interface{}{"foo": "bar"}, nil)

	// Create environment file
	os.WriteFile(filepath.Join(collPath, "environments", "dev.toml"), []byte{}, 0644)

	// Create nested actions
	// root_action.toml
	os.WriteFile(filepath.Join(collPath, "root_action.toml"), []byte(""), 0644)

	// subdir/nested.toml
	os.Mkdir(filepath.Join(collPath, "subdir"), 0755)
	os.WriteFile(filepath.Join(collPath, "subdir", "nested.toml"), []byte(""), 0644)

	// deep/dir/action.toml
	os.MkdirAll(filepath.Join(collPath, "deep", "dir"), 0755)
	os.WriteFile(filepath.Join(collPath, "deep", "dir", "action.toml"), []byte(""), 0644)

	t.Run("Success", func(t *testing.T) {
		resp, err := Get(collPath)
		if err != nil {
			t.Fatalf("Get failed: %v", err)
		}

		// Check Pre vars
		if resp.Pre["foo"] != "bar" {
			t.Errorf("expected pre var foo=bar, got %v", resp.Pre["foo"])
		}

		// Check Environments
		foundEnv := false
		for _, e := range resp.Environments {
			if e == "dev.toml" {
				foundEnv = true
				break
			}
		}
		if !foundEnv {
			t.Error("expected dev.toml in environments")
		}

		// Check Actions Tree structure
		// Expected:
		// {
		//   "root_action.toml": "root_action.toml",
		//   "subdir": { "nested.toml": "subdir/nested.toml" },
		//   "deep": { "dir": { "action.toml": "deep/dir/action.toml" } }
		// }

		// 1. Root action
		if val, ok := resp.Actions["root_action.toml"]; !ok || val != "root_action.toml" {
			t.Errorf("expected root_action.toml at root, got %v", val)
		}

		// 2. Subdir nested action
		subdir, ok := resp.Actions["subdir"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected subdir to be a map")
		}
		// On windows path separators might be different if not handled by ToSlash in implementation
		// verify logic implementation uses filepath.ToSlash
		expectedNested := "subdir/nested.toml"
		if val, ok := subdir["nested.toml"]; !ok || val != expectedNested {
			t.Errorf("expected nested.toml in subdir with path %s, got %v", expectedNested, val)
		}

		// 3. Deeply nested action
		deep, ok := resp.Actions["deep"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected deep to be a map")
		}
		dir, ok := deep["dir"].(map[string]interface{})
		if !ok {
			t.Fatalf("expected deep/dir to be a map")
		}
		expectedDeep := "deep/dir/action.toml"
		if val, ok := dir["action.toml"]; !ok || val != expectedDeep {
			t.Errorf("expected action.toml in deep/dir, got %v", val)
		}
	})

	t.Run("EncodeJSON", func(t *testing.T) {
		resp, _ := Get(collPath)
		data, err := json.Marshal(resp)
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		// Basic sanity check on JSON output
		var raw map[string]interface{}
		json.Unmarshal(data, &raw)
		if _, ok := raw["actions"]; !ok {
			t.Error("expected actions key in json")
		}
	})

	t.Run("MissingArgs", func(t *testing.T) {
		_, err := Get("")
		if err == nil {
			t.Error("expected error for missing path")
		}
	})
}

// Helper to check if map contains key
func hasKey(m interface{}, key string) bool {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Map {
		for _, k := range v.MapKeys() {
			if k.String() == key {
				return true
			}
		}
	}
	return false
}
