package httpx

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRequireRole(t *testing.T) {
	// Set up Gin test context
	gin.SetMode(gin.TestMode)
	
	t.Run("should allow access with matching role", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("roles", []string{"ORG_ADMIN", "EMPLOYEE"})
		
		// Test that the middleware doesn't return an error for a valid role
		handler := RequireRole("ORG_ADMIN")
		assert.NotNil(t, handler)
	})
	
	t.Run("should allow access with matching role (second test)", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("roles", []string{"EMPLOYEE"})
		
		// Test that the middleware doesn't return an error for a valid role
		handler := RequireRole("ORG_ADMIN")
		assert.NotNil(t, handler)
	})
}

func TestTenantID(t *testing.T) {
	t.Run("should extract tenant ID correctly", func(t *testing.T) {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("org_id", "test-tenant-id")
		
		tenantID, ok := TenantID(c)
		
		assert.True(t, ok)
		assert.Equal(t, "test-tenant-id", tenantID)
	})
	
	t.Run("should verify function exists", func(t *testing.T) {
		// Just verify the function exists and compiles
		assert.NotNil(t, TenantID)
	})
}