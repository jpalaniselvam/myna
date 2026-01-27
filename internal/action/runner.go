package action

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/aws/smithy-go"
	"github.com/jpalaniselvam/myna/internal/auth"
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

	execCtx, err := BuildExecutionContext(filePath, envName)
	if err != nil {
		actionErr = err
		return nil
	}

	action, err := LoadAction(filePath, execCtx.Resolver)
	if err != nil {
		actionErr = fmt.Errorf("failed to load action: %w", err)
		return nil
	}

	actionKind = action.base.Kind
	profile, region := execCtx.InheritConfig(action.base.Profile, action.base.Region)
	actionProfile = profile
	actionRegion = region

	baseDir := filepath.Dir(filePath)
	payloadData, err := payload.Resolve(action.base.Payload, baseDir, execCtx.Resolver)
	if err != nil {
		actionErr = fmt.Errorf("unsupported payload: %s", action.base.Payload)
		return nil
	}

	cfg, err := auth.LoadConfig(profile, region)
	if err != nil {
		actionErr = fmt.Errorf("failed to load configuration: %w", err)
		return nil
	}

	switch actionKind {
	case baseTypes.KindLambdaInvoke:
		actionResult, actionErr = runLambdaAction(cfg, action.content, payloadData)
	case baseTypes.KindSNSPublish:
		actionResult, actionErr = runSNSAction(cfg, action.content, payloadData)
	case baseTypes.KindSQSSend, baseTypes.KindSQSReceive, baseTypes.KindSQSDelete, baseTypes.KindSQSPurge, baseTypes.KindSQSMoveTask:
		actionResult, actionErr = runSQSAction(cfg, action.content, payloadData)
	case baseTypes.KindEventBridgePutEvents:
		actionResult, actionErr = runEventBridgeAction(cfg, action.content, payloadData)
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
