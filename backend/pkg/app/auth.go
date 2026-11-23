package app

import (
	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/httpx"
)

// SetupAuth sets up the authentication endpoints.
func SetupAuth(core *CoreContext, cfg *config.Config) {
	// Unauthenticated auth proxy endpoints.
	core.Engine.POST("/api/auth/login", func(c *gin.Context) {
		token, err := httpx.Login(c, cfg.KeycloakIssuerURL, cfg.KeycloakClientID, cfg.KeycloakClientSecret)
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
		token, err := httpx.Refresh(c, cfg.KeycloakIssuerURL, cfg.KeycloakClientID, cfg.KeycloakClientSecret)
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
		if err := httpx.Logout(c.Request.Context(), cfg.KeycloakIssuerURL, cfg.KeycloakClientID, cfg.KeycloakClientSecret,
			c.PostForm("refresh_token")); err != nil {
			httpx.Error(c, 502, "logout_failed", err.Error())
			return
		}
		httpx.JSON(c, 204, nil)
	})
}