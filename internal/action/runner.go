package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/smithy-go"
	"github.com/jpalaniselvam/myna/internal/auth"
	mynaContext "github.com/jpalaniselvam/myna/internal/context"
	"github.com/jpalaniselvam/myna/internal/payload"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

// RunAction reads a file and executes the action with full variable resolution
// It prints the ActionOutput JSON to stdout.
func RunAction(filePath string, envName string) error {
	var (
		actionKind    baseTypes.Kind
		actionProfile string
		actionRegion  string
		actionResult  interface{}
		actionErr     error
	)

	// Ensure we always print the output, regardless of panic or return
	defer func() {
		printActionOutput(actionKind, actionProfile, actionRegion, actionResult, actionErr)
	}()

	execCtx, err := mynaContext.BuildExecutionContext(filePath, envName)
	if err != nil {
		actionErr = err
		return nil
	}

	action, err := mynaContext.LoadAction(filePath, execCtx.Resolver)
	if err != nil {
		actionErr = fmt.Errorf("failed to load action: %w", err)
		return nil
	}

	baseAction := action.Base()

	actionKind = baseAction.Kind
	profile, region := execCtx.InheritConfig(baseAction.Credentials.Profile, baseAction.Credentials.Region)
	actionProfile = profile
	actionRegion = region

	baseDir := filepath.Dir(filePath)
	payloadData, err := payload.Resolve(baseAction.Payload, baseDir, execCtx.Resolver)
	if err != nil {
		actionErr = fmt.Errorf("unsupported payload: %s", baseAction.Payload)
		return nil
	}

	cfg, err := auth.LoadConfig(profile, region)
	if err != nil {
		actionErr = fmt.Errorf("failed to load configuration: %w", err)
		return nil
	}

	switch actionKind {
	// Lambda
	case baseTypes.KindLambdaInvoke:
		actionResult, actionErr = runLambdaAction(cfg, action.Content(), payloadData)
	// SNS
	case baseTypes.KindSNSPublish:
		actionResult, actionErr = runSNSAction(cfg, action.Content(), payloadData)
	// SQS
	case baseTypes.KindSQSSend, baseTypes.KindSQSReceive, baseTypes.KindSQSDelete, baseTypes.KindSQSPurge, baseTypes.KindSQSMoveTask:
		actionResult, actionErr = runSQSAction(cfg, action.Content(), payloadData)
	// EventBridge
	case baseTypes.KindEventBridgePutEvents:
		actionResult, actionErr = runEventBridgeAction(cfg, action.Content(), payloadData)
	// EC2
	case baseTypes.KindEC2DescribeInstances, baseTypes.KindEC2StartInstances, baseTypes.KindEC2StopInstances, baseTypes.KindEC2RebootInstances, baseTypes.KindEC2TerminateInstances:
		actionResult, actionErr = runEC2Action(cfg, action.Content())
	// SES
	case baseTypes.KindSESSendEmail, baseTypes.KindSESVerifyEmailIdentity:
		actionResult, actionErr = runSESAction(cfg, action.Content(), payloadData)
	// RDS
	case baseTypes.KindRDSDescribeDBInstances, baseTypes.KindRDSStartDBInstance, baseTypes.KindRDSStopDBInstance, baseTypes.KindRDSRebootDBInstance:
		actionResult, actionErr = runRDSAction(cfg, action.Content())
	// SFN
	case baseTypes.KindSFNListStateMachines, baseTypes.KindSFNStartExecution, baseTypes.KindSFNDescribeExecution, baseTypes.KindSFNStopExecution:
		actionResult, actionErr = runSFNAction(cfg, action.Content(), payloadData)
	default:
		actionErr = fmt.Errorf("unsupported action kind: %s", actionKind)
	}

	return nil
}

func printActionOutput(kind baseTypes.Kind, profile, region string, data interface{}, err error) {
	status := "success"
	var actionErr *baseTypes.ActionError

	if err != nil {
		status = "failure"
		actionErr = &baseTypes.ActionError{
			Message: err.Error(),
			Code:    "InternalError",
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			actionErr.Code = apiErr.ErrorCode()
			actionErr.Message = apiErr.ErrorMessage()
		}
	}

	output := baseTypes.ActionOutput{
		Metadata: baseTypes.ResponseMetadata{
			Status:    status,
			Action:    kind,
			Timestamp: time.Now(),
			Profile:   profile,
			Region:    region,
		},
		Data:  data,
		Error: actionErr,
	}

	b, marshalErr := json.MarshalIndent(output, "", "  ")
	if marshalErr != nil {
		fmt.Printf(`{"metadata":{"status":"failure"},"error":{"code":"JSONMarshalError","message":"%s"}}`, marshalErr.Error())
		return
	}

	fmt.Println(string(b))
}
