package types

// Kind represents the type of action (e.g., lambda.invoke)
type Kind string

const (
	KindLambdaInvoke         Kind = "lambda.invoke"
	KindSQSSend              Kind = "sqs.send_message"
	KindSQSReceive           Kind = "sqs.receive_message"
	KindSQSDelete            Kind = "sqs.delete_message"
	KindSQSPurge             Kind = "sqs.purge_queue"
	KindSQSMoveTask          Kind = "sqs.start_message_move_task"
	KindSNSPublish           Kind = "sns.publish"
	KindEventBridgePutEvents Kind = "eventbridge.put_events"
)

// Metadata contains common fields for all actions
type Metadata struct {
	Version     string `toml:"version"`
	Kind        Kind   `toml:"kind"`
	Description string `toml:"description"`
	Region      string `toml:"region,omitempty"`
	Profile     string `toml:"profile,omitempty"`
	RoleARN     string `toml:"role_arn,omitempty"`
	TimeoutMs   int    `toml:"timeout_ms,omitempty"`
}

// Vars contains pre. no support for post as of now
type Vars struct {
	Pre map[string]interface{} `toml:"pre"`
}

type Payload struct {
	Data string `toml:"data,omitempty"`
	File string `toml:"file,omitempty"`
}

// BaseAction combines Metadata and ActionVariables for initial parsing
type BaseAction struct {
	Metadata
	Vars
	Payload Payload `toml:"payload"`
}

type Collection struct {
	Metadata
	Vars
}
