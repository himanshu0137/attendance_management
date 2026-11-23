package app

import (
	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/httpx"
)

// SetupAuth sets up the authentication endpoints.
func SetupAuth(core *CoreContext, cfg *config.Config) {
	// Unauthenticated auth proxy endpoints.
	
	// Login user
	// @Summary User login
	// @Description Authenticate user with Keycloak and return tokens
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param credentials body gin.H true "Login credentials"
	// @Success 200 {object} gin.H
	// @Router /api/auth/login [post]
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
	
	// Refresh authentication token
	// @Summary Refresh token
	// @Description Refresh authentication token using refresh token
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param credentials body gin.H true "Refresh credentials"
	// @Success 200 {object} gin.H
	// @Router /api/auth/refresh [post]
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
	
	// Logout user
	// @Summary User logout
	// @Description Logout user and invalidate refresh token
	// @Tags Authentication
	// @Accept json
	// @Produce json
	// @Param refresh_token body gin.H true "Refresh token"
	// @Success 204 {object} gin.H
	// @Router /api/auth/logout [post]
	core.Engine.POST("/api/auth/logout", func(c *gin.Context) {
		if err := httpx.Logout(c.Request.Context(), cfg.KeycloakIssuerURL, cfg.KeycloakClientID, cfg.KeycloakClientSecret,
			c.PostForm("refresh_token")); err != nil {
			httpx.Error(c, 502, "logout_failed", err.Error())
			return
		}
		httpx.JSON(c, 204, nil)
	})
}