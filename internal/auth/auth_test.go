package auth

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

/*
Summary of Test Cases:
1. Region Override: Verifies that the explicit region argument takes precedence over the AWS_REGION environment variable.
2. Region Default: Verifies that if no region is provided, it falls back to the AWS_REGION environment variable.
3. Profile Loading: Verifies that providing a profile name loads the configuration (specifically region) associated with that profile from the config file.
4. Default Profile: Verifies that if no profile is provided, it falls back to the "default" profile or environment settings.
*/

func TestLoadConfig(t *testing.T) {
	// 1. Setup temporary AWS config/credentials files
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config")
	credsFile := filepath.Join(tmpDir, "credentials")

	configContent := `
[default]
region = us-east-1

[profile custom-profile]
region = eu-central-1
`
	// credentials are required for LoadDefaultConfig to fully succeed in many contexts,
	// though sometimes it might just load config. Safer to provide them.
	credsContent := `
[default]
aws_access_key_id = default_key
aws_secret_access_key = default_secret

[custom-profile]
aws_access_key_id = custom_key
aws_secret_access_key = custom_secret
`

	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("failed to write config file: %v", err)
	}
	if err := os.WriteFile(credsFile, []byte(credsContent), 0644); err != nil {
		t.Fatalf("failed to write credentials file: %v", err)
	}

	// 2. Set environment variables to point to these files
	t.Setenv("AWS_CONFIG_FILE", configFile)
	t.Setenv("AWS_SHARED_CREDENTIALS_FILE", credsFile)
	// Clear standard env vars to ensure we rely on files/args
	t.Setenv("AWS_REGION", "")
	t.Setenv("AWS_PROFILE", "")

	t.Run("Region Override", func(t *testing.T) {
		// Explicit region should win
		cfg, err := LoadConfig("default", "us-west-2")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.Region != "us-west-2" {
			t.Errorf("expected region 'us-west-2', got '%s'", cfg.Region)
		}
	})

	t.Run("Region Default (from validation)", func(t *testing.T) {
		// No explicit region, should load from default profile
		cfg, err := LoadConfig("default", "")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.Region != "us-east-1" {
			t.Errorf("expected region 'us-east-1' (from default profile), got '%s'", cfg.Region)
		}
	})

	t.Run("Profile Loading", func(t *testing.T) {
		// Load custom profile
		cfg, err := LoadConfig("custom-profile", "")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.Region != "eu-central-1" {
			t.Errorf("expected region 'eu-central-1' (from custom-profile), got '%s'", cfg.Region)
		}

		// Verify credentials provider (indirectly)
		// Retrieving credentials to ensure the right profile key was loaded
		creds, err := cfg.Credentials.Retrieve(context.TODO())
		if err != nil {
			t.Fatalf("failed to retrieve credentials: %v", err)
		}
		if creds.AccessKeyID != "custom_key" {
			t.Errorf("expected access key 'custom_key', got '%s'", creds.AccessKeyID)
		}
	})

	t.Run("Profile Explicit Region Override", func(t *testing.T) {
		// Load custom profile BUT override region
		cfg, err := LoadConfig("custom-profile", "sa-east-1")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.Region != "sa-east-1" {
			t.Errorf("expected region 'sa-east-1', got '%s'", cfg.Region)
		}
	})

	t.Run("Environment Variable Fallback", func(t *testing.T) {
		t.Setenv("AWS_REGION", "ap-south-1")
		// Pass empty profile and region
		cfg, err := LoadConfig("", "")
		if err != nil {
			t.Fatalf("LoadConfig failed: %v", err)
		}
		if cfg.Region != "ap-south-1" {
			t.Errorf("expected region 'ap-south-1' (from env), got '%s'", cfg.Region)
		}
	})
}
