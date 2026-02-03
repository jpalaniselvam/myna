package types

// EC2Action represents the structure of an EC2 action TOML file
type EC2Action struct {
	BaseAction
	EC2 EC2Config `toml:"ec2" json:"ec2"`
}

// EC2Config defines the EC2 configuration parameters
type EC2Config struct {
	InstanceIds []string            `toml:"instance_ids" json:"instance_ids"`
	Filters     map[string][]string `toml:"filters,omitempty" json:"filters,omitempty"`
	Hibernate   bool                `toml:"hibernate" json:"hibernate"`
	Force       bool                `toml:"force" json:"force"`
	DryRun      bool                `toml:"dry_run" json:"dry_run"`
}
