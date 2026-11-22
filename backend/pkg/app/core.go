package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/db"
	"hunrx.com/attendance_management/pkg/httpx"
)

// Setup wires the core, resolves the active modules and mounts them.
func Setup(ctx context.Context, cfg config.ServerConfig, pool *pgxpool.Pool, core *CoreContext) error {
	active, err := ResolveModules(cfg.EnabledModules)
	if err != nil {
		return err
	}
	dbPool := &db.DB{Pool: pool, DSN: cfg.DBDSN}
	if err != nil {
		return err
	}
	core.Modules = active
	core.Engine = gin.New()
	core.Engine.Use(gin.Recovery())

	// Unauthenticated auth proxy endpoints.
	core.Engine.POST("/api/auth/login", func(c *gin.Context) {
		token, err := httpx.Login(c, cfg.Issuer, cfg.ClientID, cfg.ClientSecret)
		if err != nil {
			httpx.Error(c, 401, "invalid_credentials", err.Error())
			return
		}
		httpx.JSON(c, 200, gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
			"expires_in":    token.ExpiresIn,
		})
	})
	core.Engine.POST("/api/auth/refresh", func(c *gin.Context) {
		token, err := httpx.Refresh(c, cfg.Issuer, cfg.ClientID, cfg.ClientSecret)
		if err != nil {
			httpx.Error(c, 401, "invalid_credentials", err.Error())
			return
		}
		httpx.JSON(c, 200, gin.H{
			"access_token":  token.AccessToken,
			"refresh_token": token.RefreshToken,
			"expires_in":    token.ExpiresIn,
		})
	})
	core.Engine.POST("/api/auth/logout", func(c *gin.Context) {
		if err := httpx.Logout(c.Request.Context(), cfg.Issuer, cfg.ClientID, cfg.ClientSecret,
			c.PostForm("refresh_token")); err != nil {
			httpx.Error(c, 502, "logout_failed", err.Error())
			return
		}
		httpx.JSON(c, 204, nil)
	})

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
