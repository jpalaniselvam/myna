package types

// Kind represents the type of action (e.g., lambda.invoke)
type Kind string

const (
	KindLambdaInvoke Kind = "lambda.invoke"

	KindSQSSend     Kind = "sqs.send_message"
	KindSQSReceive  Kind = "sqs.receive_message"
	KindSQSDelete   Kind = "sqs.delete_message"
	KindSQSPurge    Kind = "sqs.purge_queue"
	KindSQSMoveTask Kind = "sqs.start_message_move_task"

	KindSNSPublish Kind = "sns.publish"

	KindEventBridgePutEvents Kind = "eventbridge.put_events"

	KindEC2DescribeInstances  Kind = "ec2.describe_instances"
	KindEC2StartInstances     Kind = "ec2.start_instances"
	KindEC2StopInstances      Kind = "ec2.stop_instances"
	KindEC2RebootInstances    Kind = "ec2.reboot_instances"
	KindEC2TerminateInstances Kind = "ec2.terminate_instances"

	KindRDSDescribeDBInstances Kind = "rds.describe_db_instances"
	KindRDSStartDBInstance     Kind = "rds.start_db_instance"
	KindRDSStopDBInstance      Kind = "rds.stop_db_instance"
	KindRDSRebootDBInstance    Kind = "rds.reboot_db_instance"

	KindSFNListStateMachines Kind = "sfn.list_state_machines"
	KindSFNStartExecution    Kind = "sfn.start_execution"
	KindSFNDescribeExecution Kind = "sfn.describe_execution"
	KindSFNStopExecution     Kind = "sfn.stop_execution"

	KindSESSendEmail           Kind = "ses.send_email"
	KindSESVerifyEmailIdentity Kind = "ses.verify_email_identity"

	KindDynamoDBListTables Kind = "dynamodb.list_tables"
	KindDynamoDBPutItem    Kind = "dynamodb.put_item"
	KindDynamoDBGetItem    Kind = "dynamodb.get_item"
	KindDynamoDBUpdateItem Kind = "dynamodb.update_item"
	KindDynamoDBDeleteItem Kind = "dynamodb.delete_item"
	KindDynamoDBQuery      Kind = "dynamodb.query"
	KindDynamoDBScan       Kind = "dynamodb.scan"
)

// Metadata contains common fields for all actions
type Metadata struct {
	Version     string `toml:"version" json:"version"`
	Kind        Kind   `toml:"kind" json:"kind"`
	Description string `toml:"description" json:"description"`
	TimeoutMs   int    `toml:"timeout_ms,omitempty" json:"timeout_ms,omitempty"`
}

type Credentials struct {
	Region  string `toml:"region,omitempty" json:"region,omitempty"`
	Profile string `toml:"profile,omitempty" json:"profile,omitempty"`
	RoleARN string `toml:"role_arn,omitempty" json:"role_arn,omitempty"`
}

// Vars contains pre. no support for post as of now
type Vars struct {
	Pre map[string]interface{} `toml:"pre" json:"pre"`
}

type Payload struct {
	Data string `toml:"data,omitempty" json:"data,omitempty"`
	File string `toml:"file,omitempty" json:"file,omitempty"`
}

// BaseAction combines Metadata and ActionVariables for initial parsing
type BaseAction struct {
	Metadata
	Vars
	Credentials Credentials `toml:"credentials" json:"credentials"`
	Payload     Payload     `toml:"payload" json:"payload"`
}

type Collection struct {
	Metadata
	Credentials Credentials `toml:"credentials" json:"credentials"`
	Vars
}
