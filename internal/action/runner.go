package action

import (
	"fmt"
	"path/filepath"

	"github.com/jpalaniselvam/myna/internal/auth"
	"github.com/jpalaniselvam/myna/internal/payload"
	baseTypes "github.com/jpalaniselvam/myna/internal/types"
)

// RunAction reads a file and executes the action with full variable resolution
func RunAction(filePath string, envName string) error {
	execCtx, err := BuildExecutionContext(filePath, envName)
	if err != nil {
		return err
	}

	action, err := LoadAction(filePath, execCtx.Resolver)
	if err != nil {
		return fmt.Errorf("failed to load action: %w", err)
	}

	baseDir := filepath.Dir(filePath)
	payloadData, err := payload.Resolve(action.base.Payload, baseDir, execCtx.Resolver)
	if err != nil {
		return fmt.Errorf("unsupported payload: %s", action.base.Payload)
	}

	profile, region := execCtx.InheritConfig(action.base.Profile, action.base.Region)

	cfg, err := auth.LoadConfig(profile, region)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	switch action.base.Kind {
	case baseTypes.KindLambdaInvoke:
		return runLambdaAction(cfg, action.content, payloadData)
	case baseTypes.KindSNSPublish:
		return runSNSAction(cfg, action.content, payloadData)
	case baseTypes.KindSQSSend, baseTypes.KindSQSReceive, baseTypes.KindSQSDelete, baseTypes.KindSQSPurge, baseTypes.KindSQSMoveTask:
		return runSQSAction(cfg, action.content, payloadData)
	case baseTypes.KindEventBridgePutEvents:
		return runEventBridgeAction(cfg, action.content, payloadData)
	default:
		return fmt.Errorf("unsupported action kind: %s", action.base.Kind)
	}
}
