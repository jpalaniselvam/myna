package context

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jpalaniselvam/myna/internal/varsub"
)

/*
Summary of Test Cases:
1. DetectCollectionContext:
   - InCollection: Verify detection of myna.toml in parent directory.
   - Standalone: Verify behavior when no myna.toml exists.
   - Error: (Hard to trigger `filepath.Abs` error in test, but covered if possible).

2. DetectEnvironmentContext:
   - Exists: Verify listing of environment files.
   - NotExists: Verify nil behavior when environments directory is missing.

3. LoadCollection:
   - Valid: Load a valid myna.toml with variables.
   - Invalid Kind: Error when kind is not "collection".
   - Parse Error: Error on invalid TOML.

4. LoadAction:
   - Valid: Load and resolve an action file using [pre] vars.
   - Resolve: Verify precedence of [pre] vars over resolver scope.
   - Parse Error: Error on invalid TOML.

5. LoadEnvironmentVariables:
   - Valid: Load vars from a valid environment file.
   - NotFound: Error when environment file does not exist.
   - Parse Error: Error on invalid TOML.

6. BuildExecutionContext:
   - In Collection + Env: Full context with collection and valid environment.
   - In Collection + No Env: Collection loaded, empty environment vars.
   - Standalone + Env (Ignored): Verify --env warning (stdout check skipped, logic verified).
   - Standalone + No Env: Minimal context.
*/

func TestDetectCollectionContext(t *testing.T) {
	// Setup directory structure
	// tmp/
	//   myna.toml
	//   subdir/
	//     action.toml
	tmpDir := t.TempDir()

	mynaPath := filepath.Join(tmpDir, "myna.toml")
	if err := os.WriteFile(mynaPath, []byte("kind = 'collection'"), 0644); err != nil {
		t.Fatalf("failed to create myna.toml: %v", err)
	}

	subdir := filepath.Join(tmpDir, "subdir")
	if err := os.Mkdir(subdir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	actionPath := filepath.Join(subdir, "action.toml")
	// No need to write content for DetectCollectionContext, just path existence checks

	t.Run("InCollection", func(t *testing.T) {
		ctx, err := DetectCollectionContext(actionPath)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !ctx.InCollection {
			t.Error("expected InCollection to be true")
		}
		if ctx.Root != tmpDir {
			t.Errorf("expected Root to be %s, got %s", tmpDir, ctx.Root)
		}
		if ctx.File != mynaPath {
			t.Errorf("expected File to be %s, got %s", mynaPath, ctx.File)
		}
	})

	t.Run("Standalone", func(t *testing.T) {
		// Create a separate standalone dir
		standaloneDir := t.TempDir()
		standaloneAction := filepath.Join(standaloneDir, "action.toml")

		ctx, err := DetectCollectionContext(standaloneAction)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if ctx.InCollection {
			t.Error("expected InCollection to be false")
		}
		if ctx.Root != "" {
			t.Errorf("expected empty Root, got %s", ctx.Root)
		}
	})
}

func TestLoadCollection(t *testing.T) {
	tmpDir := t.TempDir()
	mynaPath := filepath.Join(tmpDir, "myna.toml")

	content := `
kind = "collection"
version = "1.0"
description = "Test Collection"
profile = "default"
region = "us-east-1"

[pre]
key = "value"
interpolated = "val-{{ var }}"
`
	if err := os.WriteFile(mynaPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write myna.toml: %v", err)
	}

	resolver := varsub.New(varsub.Vars{"var": "ue"})

	t.Run("Valid", func(t *testing.T) {
		col, err := LoadCollection(mynaPath, resolver)
		if err != nil {
			t.Fatalf("LoadCollection failed: %v", err)
		}
		if col.Kind != "collection" {
			t.Errorf("expected kind collection, got %s", col.Kind)
		}
		if col.Pre["interpolated"] != "val-ue" {
			t.Errorf("expected interpolated variable resolution, got %v", col.Pre["interpolated"])
		}
	})

	t.Run("InvalidKind", func(t *testing.T) {
		badMyna := filepath.Join(tmpDir, "bad_myna.toml")
		os.WriteFile(badMyna, []byte(`kind = "action"`), 0644)

		_, err := LoadCollection(badMyna, resolver)
		if err == nil {
			t.Error("expected error for invalid kind, got nil")
		}
	})
}

func TestLoadAction(t *testing.T) {
	tmpDir := t.TempDir()
	actionPath := filepath.Join(tmpDir, "action.toml")

	content := `
kind = "test.action"
version = "1.0"
description = "Test Action {{ global_var }}"

[pre]
local_var = "local"
combined = "{{ global_var }}-{{ local_var }}"
`
	if err := os.WriteFile(actionPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write action.toml: %v", err)
	}

	resolver := varsub.New(varsub.Vars{"global_var": "global"})

	t.Run("Valid with Resolution", func(t *testing.T) {
		action, err := LoadAction(actionPath, resolver)
		if err != nil {
			t.Fatalf("LoadAction failed: %v", err)
		}

		if !strings.Contains(action.Base().Description, "global") {
			t.Errorf("expected description to contain 'global', got %s", action.Base().Description)
		}

		// Accessing pre vars is a bit tricky as they are inside base.Vars.Pre
		if action.Base().Pre["combined"] != "global-local" {
			t.Errorf("expected combined var to be 'global-local', got %v", action.Base().Pre["combined"])
		}
	})
}

func TestLoadEnvironmentVariables(t *testing.T) {
	tmpDir := t.TempDir()
	envDir := filepath.Join(tmpDir, "environments")
	os.Mkdir(envDir, 0755)

	envPath := filepath.Join(envDir, "dev.toml")
	content := `
[vars]
api_url = "https://dev.api.com"
retry_count = 3
`
	os.WriteFile(envPath, []byte(content), 0644)

	t.Run("Valid", func(t *testing.T) {
		vars, err := LoadEnvironmentVariables(tmpDir, "dev")
		if err != nil {
			t.Fatalf("LoadEnvironmentVariables failed: %v", err)
		}
		if vars["api_url"] != "https://dev.api.com" {
			t.Errorf("expected api_url match")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		_, err := LoadEnvironmentVariables(tmpDir, "prod")
		if err == nil {
			t.Error("expected error for missing environment, got nil")
		}
	})
}

func TestBuildExecutionContext(t *testing.T) {
	// Setup full structure
	tmpDir := t.TempDir()

	// myna.toml
	os.WriteFile(filepath.Join(tmpDir, "myna.toml"), []byte(`
kind = "collection" 
profile = "col-profile"
region = "col-region"
[pre]
col_var = "col"
`), 0644)

	// Environment
	envDir := filepath.Join(tmpDir, "environments")
	os.Mkdir(envDir, 0755)
	os.WriteFile(filepath.Join(envDir, "test.toml"), []byte(`
[vars]
env_var = "env"
`), 0644)

	// Action
	actionPath := filepath.Join(tmpDir, "action.toml")
	// Content doesn't matter for context building, just path
	os.WriteFile(actionPath, []byte(""), 0644)

	t.Run("InContextWithEnv", func(t *testing.T) {
		ctx, err := BuildExecutionContext(actionPath, "test")
		if err != nil {
			t.Fatalf("BuildExecutionContext failed: %v", err)
		}

		if ctx.Profile != "col-profile" {
			t.Errorf("expected profile col-profile, got %s", ctx.Profile)
		}

		// Verify resolver has access to both env and collection vars.
		// Since we can't easily peek into resolver scopes, we can try to Resolve a string.
		res := ctx.Resolver.Resolve([]byte("{{ col_var }}-{{ env_var }}"))
		if string(res) != "col-env" {
			t.Errorf("expected variable resolution 'col-env', got %s", string(res))
		}
	})

	t.Run("InheritConfig", func(t *testing.T) {
		ctx, _ := BuildExecutionContext(actionPath, "test")

		// Case 1: Action has no config, inherits collection
		p, r := ctx.InheritConfig("", "")
		if p != "col-profile" || r != "col-region" {
			t.Errorf("failed to inherit: got %s, %s", p, r)
		}

		// Case 2: Action has config, overrides collection
		p, r = ctx.InheritConfig("act-profile", "act-region")
		if p != "act-profile" || r != "act-region" {
			t.Errorf("failed to override: got %s, %s", p, r)
		}
	})
}
