package app

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"hunrx.com/attendance_management/pkg/config"
	"hunrx.com/attendance_management/pkg/httpx"
)

// SetupAuth sets up the authentication endpoints.
func SetupAuth(core *CoreContext, cfg config.ServerConfig) {
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
}

// Server is the single binary entrypoint.
type Server struct {
	Engine *gin.Engine
}

// New creates the server from the core.
func New(core *CoreContext) (*Server, error) {
	return &Server{Engine: core.Engine}, nil
}

// Run serves the engine on host:port.
func (s *Server) Run(host string, port int) error {
	srv := &http.Server{Addr: host + ":" + strconv.Itoa(port), Handler: s.Engine}
	return srv.ListenAndServe()
}