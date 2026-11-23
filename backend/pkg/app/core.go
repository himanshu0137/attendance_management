package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/db"
)

// Setup wires the core, resolves the active modules and mounts them.
func Setup(ctx context.Context, cfg config.ServerConfig, pool *pgxpool.Pool, core *CoreContext) error {
	active, err := ResolveModules(cfg.EnabledModules)
	if err != nil {
		return err
	}
	dbPool := &db.DB{Pool: pool, DSN: cfg.DatabaseURL}
	if err != nil {
		return err
	}
	core.Modules = active
	core.Engine = gin.New()
	core.Engine.Use(gin.Recovery())

	// Setup authentication endpoints
	SetupAuth(core, cfg)

	for _, m := range active {
		m.Migrate(ctx, dbPool)
		// core.Engine.GET("/healthz", func(c *gin.Context) {
		// 	httpx.JSON(c, 200, gin.H{
		// 		"status":         "ok",
		// 		"active_modules": active,
		// 	})
		// })
		m.Mount(core, core.Engine.Group("/api/"+m.Name()))
	}
	return nil
}
