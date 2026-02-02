package main

import (
	"context"
	"fmt"

	"github.com/jpalaniselvam/myna/internal/collection"
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
func (a *App) UpdateCollection(collectionPath, name, desc string, pre map[string]interface{}) error {
	return collection.Update(collectionPath, name, desc, pre)
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

// UpdateEnvVar updates a key-value pair in an environment
func (a *App) UpdateEnvVar(collectionPath, envName, key, value string) error {
	return collection.UpdateEnvVar(collectionPath, envName, key, value)
}

// DeleteEnvVar deletes a key-value pair from an environment
func (a *App) DeleteEnvVar(collectionPath, envName, key string) error {
	return collection.DeleteEnvVar(collectionPath, envName, key)
}
