package payload

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	baseTypes "github.com/jpalaniselvam/myna/internal/types"
	"github.com/jpalaniselvam/myna/internal/varsub"
)

/*
Summary of Test Cases:
1. File Payload (Absolute Path): Verifies reading content from an existing file using an absolute path.
2. File Payload (Relative Path): Verifies reading content using a path relative to the base directory.
3. File Payload (Non-Existent): Verifies error handling when the file does not exist.
4. Data Payload (Inline): Verifies using inline string content when no file is specified.
5. Precedence (File over Data): Verifies that 'file' takes precedence over 'data' when both are present.
6. Empty Payload: Verifies that nil is returned when neither 'file' nor 'data' is specified.
7. Variable Substitution (File): Verifies variables are resolved within file content.
8. Variable Substitution (Data): Verifies variables are resolved within inline data.
*/

func TestResolve(t *testing.T) {
	// Create a temporary directory for file-based tests
	tmpDir, err := os.MkdirTemp("", "payload_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a dummy file for absolute path test
	absFilePath := filepath.Join(tmpDir, "abs_payload.json")
	absContent := `{"key": "value"}`
	if err := os.WriteFile(absFilePath, []byte(absContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Create a dummy file for relative path test
	relFileName := "rel_payload.json"
	relFilePath := filepath.Join(tmpDir, relFileName)
	relContent := `{"foo": "bar"}`
	if err := os.WriteFile(relFilePath, []byte(relContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Create a dummy file for variable substitution test
	varFileName := "var_payload.txt"
	varFilePath := filepath.Join(tmpDir, varFileName)
	varContent := `Hello {{ name }}!`
	if err := os.WriteFile(varFilePath, []byte(varContent), 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	// Setup variable resolver with some test data
	vars := varsub.Vars{
		"name": "World",
		"env":  "production",
	}
	resolver := varsub.New(vars)

	tests := []struct {
		name          string
		payload       baseTypes.Payload
		baseDir       string
		wantErr       bool
		expectedData  string // Substring match or exact match depending on context, or use strict equality
		errorContains string
	}{
		{
			name: "File Payload (Absolute Path)",
			payload: baseTypes.Payload{
				File: absFilePath,
			},
			baseDir:      "", // Should not matter for absolute path
			wantErr:      false,
			expectedData: absContent,
		},
		{
			name: "File Payload (Relative Path)",
			payload: baseTypes.Payload{
				File: relFileName,
			},
			baseDir:      tmpDir,
			wantErr:      false,
			expectedData: relContent,
		},
		{
			name: "File Payload (Non-Existent)",
			payload: baseTypes.Payload{
				File: "non_existent_file.json",
			},
			baseDir:       tmpDir,
			wantErr:       true,
			errorContains: "failed to read payload file",
		},
		{
			name: "Data Payload (Inline)",
			payload: baseTypes.Payload{
				Data: `{"inline": "data"}`,
			},
			baseDir:      "",
			wantErr:      false,
			expectedData: `{"inline": "data"}`,
		},
		{
			name: "Precedence (File over Data)",
			payload: baseTypes.Payload{
				File: absFilePath,
				Data: `{"inline": "ignored"}`,
			},
			baseDir:      "",
			wantErr:      false,
			expectedData: absContent,
		},
		{
			name: "Empty Payload",
			payload: baseTypes.Payload{
				File: "",
				Data: "",
			},
			baseDir:      "",
			wantErr:      false,
			expectedData: "", // Expecting nil byte slice which is empty string-ish
		},
		{
			name: "Variable Substitution (File)",
			payload: baseTypes.Payload{
				File: varFileName,
			},
			baseDir:      tmpDir,
			wantErr:      false,
			expectedData: "Hello World!",
		},
		{
			name: "Variable Substitution (Data)",
			payload: baseTypes.Payload{
				Data: "Environment: {{ env }}",
			},
			baseDir:      "",
			wantErr:      false,
			expectedData: "Environment: production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Resolve(tt.payload, tt.baseDir, resolver)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Resolve() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Resolve() error = %v, want error containing %v", err, tt.errorContains)
				}
				return
			}

			if err != nil {
				t.Errorf("Resolve() unexpected error = %v", err)
				return
			}

			if tt.expectedData == "" && len(got) == 0 {
				// Matched nil/empty return
				return
			}

			if string(got) != tt.expectedData {
				t.Errorf("Resolve() got = %s, want %s", string(got), tt.expectedData)
			}
		})
	}
}
