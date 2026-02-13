package collection

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jpalaniselvam/myna/internal/types"
	"github.com/pelletier/go-toml/v2"
)

// CollectionResponse is the structure returned by Get
type CollectionResponse struct {
	Description  string                 `json:"description"`
	Environments []string               `json:"environments"`
	Pre          map[string]interface{} `json:"pre"`
	Actions      map[string]interface{} `json:"actions"`
	Credentials  types.Credentials      `json:"credentials"`
}

// Create creates a new collection
func Create(baseDir, name, desc string) error {
	if baseDir == "" || name == "" {
		return fmt.Errorf("base-dir and name are required")
	}

	collectionPath := filepath.Join(baseDir, name)
	if _, err := os.Stat(collectionPath); !os.IsNotExist(err) {
		return fmt.Errorf("collection directory already exists: %s", collectionPath)
	}

	if err := os.MkdirAll(filepath.Join(collectionPath, "environments"), 0755); err != nil {
		return fmt.Errorf("failed to create directory structure: %w", err)
	}

	coll := types.Collection{
		Metadata: types.Metadata{
			Version:     "1.0",
			Kind:        types.Kind("collection"),
			Description: desc,
		},
		Vars: types.Vars{
			Pre: make(map[string]interface{}),
		},
	}

	data, err := toml.Marshal(coll)
	if err != nil {
		return fmt.Errorf("failed to marshal collection metadata: %w", err)
	}

	mynaTomlPath := filepath.Join(collectionPath, "myna.toml")
	if err := os.WriteFile(mynaTomlPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write myna.toml: %w", err)
	}

	fmt.Printf("Collection created at %s\n", collectionPath)
	return nil
}

// Update updates an existing collection
func Update(collectionPath, name, desc string, pre map[string]interface{}, creds *types.Credentials) error {
	if collectionPath == "" {
		return fmt.Errorf("base-dir (collection path) is required")
	}

	mynaTomlPath := filepath.Join(collectionPath, "myna.toml")

	// Read existing
	data, err := os.ReadFile(mynaTomlPath)
	if err != nil {
		return fmt.Errorf("failed to read myna.toml: %w", err)
	}

	var coll types.Collection
	if err := toml.Unmarshal(data, &coll); err != nil {
		return fmt.Errorf("failed to parse myna.toml: %w", err)
	}

	// Update fields
	if desc != "" {
		coll.Metadata.Description = desc
	}

	if pre != nil {
		coll.Vars.Pre = pre
	}
	if creds != nil {
		coll.Credentials = *creds
	}

	// Save updates
	newData, err := toml.Marshal(coll)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}
	if err := os.WriteFile(mynaTomlPath, newData, 0644); err != nil {
		return fmt.Errorf("failed to write myna.toml: %w", err)
	}

	// Handle Rename if name is provided and different from current dir name
	if name != "" {
		currentName := filepath.Base(collectionPath)
		if name != currentName {
			parentDir := filepath.Dir(collectionPath)
			newPath := filepath.Join(parentDir, name)
			if err := os.Rename(collectionPath, newPath); err != nil {
				return fmt.Errorf("failed to rename collection: %w", err)
			}
			fmt.Printf("Collection renamed to %s\n", newPath)
			return nil
		}
	}

	fmt.Printf("Collection updated at %s\n", collectionPath)
	return nil
}

// Delete deletes a collection
func Delete(collectionPath string) error {
	if collectionPath == "" {
		return fmt.Errorf("base-dir (collection path) is required")
	}

	// Confirm path exists
	if _, err := os.Stat(collectionPath); os.IsNotExist(err) {
		return fmt.Errorf("collection directory not found: %s", collectionPath)
	}

	// Delete
	if err := os.RemoveAll(collectionPath); err != nil {
		return fmt.Errorf("failed to delete collection: %w", err)
	}

	fmt.Printf("Collection deleted: %s\n", collectionPath)
	return nil
}

// Get retrieves collection details
func Get(collectionPath string) (*CollectionResponse, error) {
	if collectionPath == "" {
		return nil, fmt.Errorf("base-dir (collection path) is required")
	}

	mynaTomlPath := filepath.Join(collectionPath, "myna.toml")

	// Read Metadata/Vars
	data, err := os.ReadFile(mynaTomlPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read myna.toml: %w", err)
	}
	var coll types.Collection
	if err := toml.Unmarshal(data, &coll); err != nil {
		return nil, fmt.Errorf("failed to parse myna.toml: %w", err)
	}

	response := &CollectionResponse{
		Description:  coll.Metadata.Description,
		Pre:          coll.Vars.Pre,
		Environments: []string{},
		Actions:      make(map[string]interface{}),
		Credentials:  coll.Credentials,
	}

	// List Environments
	envDir := filepath.Join(collectionPath, "environments")
	entries, err := os.ReadDir(envDir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() && strings.HasSuffix(e.Name(), ".toml") {
				response.Environments = append(response.Environments, e.Name())
			}
		}
	}

	// List Actions (recursively find .toml files excluding myna.toml and environments dir)
	err = filepath.Walk(collectionPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// skip environments directory
			if filepath.Base(path) == "environments" && filepath.Dir(path) == collectionPath {
				return filepath.SkipDir
			}
			return nil
		}

		// Check for .toml files
		if strings.HasSuffix(info.Name(), ".toml") {
			// exclude myna.toml in root
			if info.Name() == "myna.toml" && filepath.Dir(path) == collectionPath {
				return nil
			}

			// Get relative path
			relPath, err := filepath.Rel(collectionPath, path)
			if err != nil {
				return nil
			}
			relPath = filepath.ToSlash(relPath)
			parts := strings.Split(relPath, "/")
			addToTree(response.Actions, parts, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to walk collection directory: %w", err)
	}

	return response, nil
}

func addToTree(tree map[string]interface{}, parts []string, fullPath string) {
	if len(parts) == 0 {
		return
	}
	name := parts[0]
	if len(parts) == 1 {
		// Leaf node (file)
		// We can store just the file name as key, and maybe the full relative path as value?
		// User requirement: "maintaining the directory structure, just the name"
		// If we use 'name' as key, we need a value.
		// Let's use the full relative path as the value so the UI can execute it.
		tree[name] = fullPath
		return
	}

	// Directory node
	subTreeRaw, exists := tree[name]
	var subTree map[string]interface{}
	if !exists {
		subTree = make(map[string]interface{})
		tree[name] = subTree
	} else {
		subTree, _ = subTreeRaw.(map[string]interface{})
	}

	addToTree(subTree, parts[1:], fullPath)
}
