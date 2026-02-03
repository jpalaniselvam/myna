package types

// RDSAction represents the structure of an RDS action TOML file
type RDSAction struct {
	BaseAction
	RDS RDSConfig `toml:"rds" json:"rds"`
}

// RDSConfig defines the RDS configuration parameters
type RDSConfig struct {
	DBInstanceIdentifier string              `toml:"db_instance_identifier" json:"db_instance_identifier"`
	Filters              map[string][]string `toml:"filters,omitempty" json:"filters,omitempty"`
	DBSnapshotIdentifier string              `toml:"db_snapshot_identifier,omitempty" json:"db_snapshot_identifier,omitempty"`
	ForceFailover        bool                `toml:"force_failover" json:"force_failover"`
}
