package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/httpx"
)

// Setup wires the core, resolves the active modules and mounts them.
func Setup(ctx context.Context, cfg *config.Config, pool *pgxpool.Pool, core *CoreContext) error {
	core.Engine = gin.New()
	core.Engine.Use(gin.Recovery())

	// Setup authentication endpoints
	SetupAuth(core, cfg)

	core.Engine.GET("/healthz", func(c *gin.Context) {
		httpx.JSON(c, 200, gin.H{
			"status": "ok",
		})
	})

	// Setup Swagger documentation
	SetupSwagger(core.Engine)

	return nil
}
