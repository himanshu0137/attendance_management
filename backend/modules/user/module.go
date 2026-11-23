package user

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"

	"hunrx.com/attendance_management/pkg/app"
	"hunrx.com/attendance_management/pkg/db"
	"hunrx.com/attendance_management/pkg/httpx"
)

// init registers the module with the core.
func init() {
	app.Register(NewModule())
}

// Module is the feature module registration.
type Module struct {
	service *Service
}

// NewModule returns the module instance.
func NewModule() app.FeatureModule {
	return Module{}
}

// Name returns the module name.
func (m Module) Name() string {
	return "user"
}

// Migrate runs the user schema migrations.
func (m Module) Migrate(ctx context.Context, db *db.DB) error {
	return db.Migrate(ctx, "user", "../migrations")
}

// Mount wires the handlers.
func (m Module) Mount(core *app.CoreContext, rg *gin.RouterGroup) {
	repo := NewRepository(core.Pool)
	svc := NewService(repo)
	m.service = svc
	rg.GET("/me", m.getMe)
	rg.PUT("/me", m.updateMe)
	rg.GET("/me/preferences", m.getPreferences)
	rg.PUT("/me/preferences", m.updatePreferences)
}

// GetMyDetails gets current user details
// @Summary Get user details
// @Description Retrieve current user's personal details
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} User
// @Router /api/user/me [get]
func (m Module) getMe(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Error(c, 401, "missing_user", "no user_id in token")
		return
	}
	
	user, err := m.service.GetMe(c.Request.Context(), tenantID, userID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, user)
}

// UpdateMyDetails updates current user details
// @Summary Update user details
// @Description Update current user's personal details
// @Tags User
// @Accept json
// @Produce json
// @Param user body UpdateUserReq true "User details"
// @Success 200 {object} User
// @Router /api/user/me [put]
func (m Module) updateMe(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Error(c, 401, "missing_user", "no user_id in token")
		return
	}
	
	req, ok := httpx.BindJSON[UpdateUserReq](c)
	if !ok {
		return
	}
	
	user, err := m.service.UpdateMe(c.Request.Context(), tenantID, userID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, user)
}

// GetUserPreferences gets current user preferences
// @Summary Get user preferences
// @Description Retrieve current user's application preferences
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} Preferences
// @Router /api/user/me/preferences [get]
func (m Module) getPreferences(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Error(c, 401, "missing_user", "no user_id in token")
		return
	}
	
	prefs, err := m.service.GetPreferences(c.Request.Context(), tenantID, userID)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, prefs)
}

// UpdateUserPreferences updates current user preferences
// @Summary Update user preferences
// @Description Update current user's application preferences
// @Tags User
// @Accept json
// @Produce json
// @Param preferences body Preferences true "User preferences"
// @Success 200 {object} Preferences
// @Router /api/user/me/preferences [put]
func (m Module) updatePreferences(c *gin.Context) {
	tenantID, ok := httpx.TenantID(c)
	if !ok {
		httpx.Error(c, 401, "missing_tenant", "no tenant_id in token")
		return
	}
	
	userID, ok := httpx.UserID(c)
	if !ok {
		httpx.Error(c, 401, "missing_user", "no user_id in token")
		return
	}
	
	req, ok := httpx.BindJSON[Preferences](c)
	if !ok {
		return
	}
	
	prefs, err := m.service.UpdatePreferences(c.Request.Context(), tenantID, userID, req)
	if err != nil {
		httpx.Error(c, 500, "db_error", err.Error())
		return
	}
	
	httpx.JSON(c, 200, prefs)
}

// UpdateUserReq represents the request for updating user details
type UpdateUserReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Language string `json:"language"`
}

// Preferences represents user application preferences
type Preferences struct {
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
	Language      string `json:"language"`
	DateFormat    string `json:"date_format"`
	AutoClockIn   bool   `json:"auto_clock_in"`
}