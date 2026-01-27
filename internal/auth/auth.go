package auth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// LoadConfig loads AWS configuration based on profile and region.
// It follows the standard AWS SDK credential chain:
// 1. Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN)
// 2. EC2/ECS/Lambda instance profiles
// 3. Shared credentials file (~/.aws/credentials)
// 4. SSO cache
func LoadConfig(profile, region string) (aws.Config, error) {
	ctx := context.TODO()
	opts := []func(*config.LoadOptions) error{}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	return config.LoadDefaultConfig(ctx, opts...)
}
