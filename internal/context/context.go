package context

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jpalaniselvam/myna/internal/parser"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
	"github.com/jpalaniselvam/myna/internal/varsub"
)

// CollectionContext holds information about the collection context
type CollectionContext struct {
	Root         string // Absolute path to collection root
	InCollection bool   // Whether action is in a collection
	File         string // path to myna file
}

// DetectCollectionContext walks up the directory tree from the action file
// looking for myna.toml. Returns collection context information.
func DetectCollectionContext(actionPath string) (*CollectionContext, error) {
	ctx := &CollectionContext{
		InCollection: false,
	}

	// Get absolute path
	absPath, err := filepath.Abs(actionPath)
	if err != nil {
		return ctx, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Start from the directory containing the action file
	dir := filepath.Dir(absPath)

	// Walk up the directory tree
	for {
		// Check if myna.toml exists in current directory
		mynaPath := filepath.Join(dir, "myna.toml")
		if _, err := os.Stat(mynaPath); err == nil {
			// Found myna.toml - this is the collection root
			ctx.Root = dir
			ctx.InCollection = true
			ctx.File = mynaPath

			return ctx, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root of filesystem, no collection found
			break
		}
		dir = parent
	}

	// No collection found - standalone context
	return ctx, nil
}

func DetectEnvironmentContext(collectionRoot string) ([]string, error) {
	envDir := filepath.Join(collectionRoot, "environments")
	if _, err := os.Stat(envDir); os.IsNotExist(err) {
		return nil, nil
	}

	files, err := filepath.Glob(filepath.Join(envDir, "*.toml"))
	if err != nil {
		return nil, fmt.Errorf("failed to list environment files: %w", err)
	}

	return files, nil
}

// LoadCollection reads and parses the myna.toml file
func LoadCollection(filePath string, resolver *varsub.Resolver) (*baseTypes.Collection, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read myna.toml: %w", err)
	}

	content = resolver.Resolve(content)

	var metadata baseTypes.Collection
	if err := parser.Unmarshal(content, &metadata); err != nil {
		return nil, fmt.Errorf("failed to parse myna.toml: %w", err)
	}

	// Validate it's actually a collection
	if metadata.Kind != "collection" {
		return nil, fmt.Errorf("myna.toml is not a collection (kind: %s)", metadata.Kind)
	}

	return &metadata, nil
}

type ResolvedAction struct {
	base    baseTypes.BaseAction
	content []byte
}

func (r *ResolvedAction) Base() baseTypes.BaseAction {
	return r.base
}

func (r *ResolvedAction) Content() []byte {
	return r.content
}

func LoadAction(filePath string, resolver *varsub.Resolver) (*ResolvedAction, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read action file: %w", err)
	}

	// 1. Initial parse into generic map to get [pre] block (unresolved)
	// We use map[string]interface{} to avoid type errors with unresolved {{var}} in non-string fields
	var genericMap map[string]interface{}
	if err := parser.Unmarshal(content, &genericMap); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", filePath, err)
	}

	// 2. Add [pre] variables to resolver scope if they exist
	if pre, ok := genericMap["pre"].(map[string]interface{}); ok {
		resolver.AddScope(pre)
	}

	// 3. Resolve variables in the file content using the updated resolver
	content = resolver.Resolve(content)

	// 4. Final parse of the resolved content into the typed struct
	var baseAction baseTypes.BaseAction
	if err := parser.Unmarshal(content, &baseAction); err != nil {
		return nil, fmt.Errorf("failed to parse resolved %s: %w", filePath, err)
	}

	action := &ResolvedAction{
		base:    baseAction,
		content: content,
	}

	return action, nil
}

// LoadEnvironmentVariables loads variables from an environment TOML file
func LoadEnvironmentVariables(collectionRoot, envName string) (map[string]interface{}, error) {
	envPath := filepath.Join(collectionRoot, "environments", envName+".toml")

	// Check if file exists
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("environment '%s' not found at %s", envName, envPath)
	}

	content, err := os.ReadFile(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment file: %w", err)
	}

	var envData struct {
		Vars map[string]interface{} `toml:"vars"`
	}

	if err := parser.Unmarshal(content, &envData); err != nil {
		return nil, fmt.Errorf("failed to parse environment file: %w", err)
	}

	if envData.Vars == nil {
		return make(map[string]interface{}), nil
	}

	return envData.Vars, nil
}

// ExecutionContext holds the resolved execution context for an action
type ExecutionContext struct {
	Resolver    *varsub.Resolver
	Collection  *baseTypes.Collection
	Credentials *baseTypes.Credentials
}

// BuildExecutionContext creates an execution context by loading environment, collection, and resolver
func BuildExecutionContext(actionPath, envName string) (*ExecutionContext, error) {
	ctx, err := DetectCollectionContext(actionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to detect collection context: %w", err)
	}

	var envVars map[string]interface{}
	if envName != "" {
		if ctx.InCollection {
			envVars, err = LoadEnvironmentVariables(ctx.Root, envName)
			if err != nil {
				return nil, fmt.Errorf("failed to load environment variables: %w", err)
			}
		} else {
			fmt.Fprintf(os.Stderr, "WARN: --env flag ignored (action not in collection context)\n")
			envVars = make(map[string]interface{})
		}
	} else {
		envVars = make(map[string]interface{})
	}

	resolver := varsub.New(envVars)

	var collection *baseTypes.Collection
	if ctx.InCollection {
		collection, err = LoadCollection(ctx.File, resolver)
		if err != nil {
			return nil, fmt.Errorf("failed to load collection: %w", err)
		}
		resolver.AddScope(collection.Pre)
	}

	var creds *baseTypes.Credentials
	if collection != nil {
		creds = &collection.Credentials
	}

	return &ExecutionContext{
		Resolver:    resolver,
		Collection:  collection,
		Credentials: creds,
	}, nil
}

// InheritConfig returns profile and region, inheriting from collection if action values are empty
func (ctx *ExecutionContext) InheritConfig(actionProfile, actionRegion string) (profile, region string) {
	profile = actionProfile
	if profile == "" && ctx.Credentials != nil {
		profile = ctx.Credentials.Profile
	}

	region = actionRegion
	if region == "" && ctx.Credentials != nil {
		region = ctx.Credentials.Region
	}

	return profile, region
}
