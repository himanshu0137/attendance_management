package app

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// CoreContext holds the core application context.
type CoreContext struct {
	Engine  *gin.Engine
	Modules []FeatureModule
	Pool    *pgxpool.Pool
}

// FeatureModule represents a feature module in the application.
type FeatureModule interface {
	// Name returns the name of the module.
	Name() string

	// Mount mounts the module's routes onto the given router group.
	Mount(core *CoreContext, group *gin.RouterGroup)
}

// DB represents a database connection with tenant isolation.
type DB struct {
	Pool *pgxpool.Pool
	DSN  string
}
