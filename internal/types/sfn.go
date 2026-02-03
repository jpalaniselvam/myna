package types

// SFNAction represents the structure of a Step Functions action TOML file
type SFNAction struct {
	BaseAction
	SFN SFNConfig `toml:"sfn" json:"sfn"`
}

// SFNConfig defines the Step Functions configuration parameters
type SFNConfig struct {
	StateMachineARN string `toml:"state_machine_arn,omitempty" json:"state_machine_arn,omitempty"`
	Name            string `toml:"name,omitempty" json:"name,omitempty"`
	ExecutionARN    string `toml:"execution_arn,omitempty" json:"execution_arn,omitempty"`
	Error           string `toml:"error,omitempty" json:"error,omitempty"`
	Cause           string `toml:"cause,omitempty" json:"cause,omitempty"`
}
