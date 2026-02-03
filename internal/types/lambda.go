package types

import "encoding/json"

// LambdaAction represents the structure of a lambda.invoke action TOML file
type LambdaAction struct {
	BaseAction
	Lambda LambdaConfig `toml:"lambda" json:"lambda"`
}

// LambdaConfig defines the lambda function configuration
type LambdaConfig struct {
	FunctionName   string                 `toml:"function_name" json:"function_name"`
	InvocationType string                 `toml:"invocation_type" json:"invocation_type"` // RequestResponse (default), Event, DryRun
	Qualifier      string                 `toml:"qualifier,omitempty" json:"qualifier,omitempty"`
	ClientContext  map[string]interface{} `toml:"client_context,omitempty" json:"client_context,omitempty"`
}

// LambdaOutput represents the output of a lambda execution
type LambdaOutput struct {
	StatusCode      int                    `json:"status_code"`
	ExecutedVersion string                 `json:"executed_version"`
	Payload         map[string]interface{} `json:"payload"`
	Logs            string                 `json:"logs,omitempty"`
	RawPayload      []byte                 `json:"-"` // Internal use, for when response is not JSON
}

func (l *LambdaOutput) String() string {
	b, _ := json.MarshalIndent(l, "", "  ")
	return string(b)
}
