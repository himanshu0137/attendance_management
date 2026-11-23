package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/db"
)

// CoreContext holds the core application context.
type CoreContext struct {
	Engine   *gin.Engine
	Modules  []FeatureModule
}

// FeatureModule represents a feature module in the application.
type FeatureModule interface {
	// Name returns the name of the module.
	Name() string
	
	// Migrate runs database migrations for this module.
	Migrate(ctx context.Context, db *db.DB) error
	
	// Mount mounts the module's routes onto the given router group.
	Mount(core *CoreContext, group *gin.RouterGroup)
}

// DB represents a database connection with tenant isolation.
type DB struct {
	Pool *pgxpool.Pool
	DSN  string
}

// Register registers a feature module.
func Register(module FeatureModule) {
	modules = append(modules, module)
}

// ResolveModules resolves the modules based on the requested names.
func ResolveModules(requested []string) ([]FeatureModule, error) {
	if len(requested) == 0 {
		return modules, nil
	}
	
	var resolved []FeatureModule
	for _, name := range requested {
		for _, module := range modules {
			if module.Name() == name {
				resolved = append(resolved, module)
				break
			}
		}
	}
	
	return resolved, nil
}

// Core modules registry.
var modules []FeatureModule