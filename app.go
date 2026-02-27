package main

import (
	"context"
	"fmt"

	"github.com/jpalaniselvam/myna/internal/collection"
	"github.com/jpalaniselvam/myna/internal/types"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// CreateCollection creates a new collection
func (a *App) CreateCollection(baseDir, name, desc string) error {
	return collection.Create(baseDir, name, desc)
}

// UpdateCollection updates an existing collection
func (a *App) UpdateCollection(collectionPath, name, desc string, pre map[string]interface{}, creds *types.Credentials) error {
	return collection.Update(collectionPath, name, desc, pre, creds)
}

// DeleteCollection deletes a completely
func (a *App) DeleteCollection(collectionPath string) error {
	return collection.Delete(collectionPath)
}

// GetCollection gets the collection details
func (a *App) GetCollection(collectionPath string) (*collection.CollectionResponse, error) {
	return collection.Get(collectionPath)
}

// CreateEnvironment creates a new environment in a collection
func (a *App) CreateEnvironment(collectionPath, envName string) error {
	return collection.CreateEnvironment(collectionPath, envName)
}

// AddEnvVar adds a key-value pair to an environment
func (a *App) AddEnvVar(collectionPath, envName, key, value string) error {
	return collection.AddEnvVar(collectionPath, envName, key, value)
}

// GetEnvironment retrieves the variables of an environment
func (a *App) GetEnvironment(collectionPath, envName string) (map[string]interface{}, error) {
	return collection.GetEnvironment(collectionPath, envName)
}

// UpdateEnvVar updates a key-value pair in an environment
func (a *App) UpdateEnvVar(collectionPath, envName, key, value string) error {
	return collection.UpdateEnvVar(collectionPath, envName, key, value)
}

// DeleteEnvVar deletes a key-value pair from an environment
func (a *App) DeleteEnvVar(collectionPath, envName, key string) error {
	return collection.DeleteEnvVar(collectionPath, envName, key)
}

// CreateAction creates a new action in a collection
func (a *App) CreateAction(input collection.CreateActionInput) error {
	return collection.CreateAction(input)
}

// UpdateAction updates an existing action within a collection
func (a *App) UpdateAction(input collection.UpdateActionInput) error {
	return collection.UpdateAction(input)
}

// DeleteAction deletes an action from a collection
func (a *App) DeleteAction(input collection.GetActionInput) error {
	return collection.DeleteAction(input)
}

// GetAction retrieves an action from a collection
func (a *App) GetAction(input collection.GetActionInput) (interface{}, error) {
	return collection.GetAction(input)
}

// SelectDirectory opens a directory selection dialog
func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Collection Directory",
	})
}

// GetActionKinds returns all supported action kinds
func (a *App) GetActionKinds() []types.Kind {
	return []types.Kind{
		types.KindLambdaInvoke,
		types.KindSQSSend,
		types.KindSQSReceive,
		types.KindSQSDelete,
		types.KindSQSPurge,
		types.KindSQSMoveTask,
		types.KindSNSPublish,
		types.KindEventBridgePutEvents,
		types.KindEC2DescribeInstances,
		types.KindEC2StartInstances,
		types.KindEC2StopInstances,
		types.KindEC2RebootInstances,
		types.KindEC2TerminateInstances,
		types.KindRDSDescribeDBInstances,
		types.KindRDSStartDBInstance,
		types.KindRDSStopDBInstance,
		types.KindRDSRebootDBInstance,
		types.KindSFNListStateMachines,
		types.KindSFNStartExecution,
		types.KindSFNDescribeExecution,
		types.KindSFNStopExecution,
		types.KindSESSendEmail,
		types.KindSESVerifyEmailIdentity,
	}
}
