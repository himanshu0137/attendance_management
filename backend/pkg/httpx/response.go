package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON binds JSON from the request body to a struct and returns false if it fails.
func BindJSON[T any](c *gin.Context) (T, bool) {
	var result T
	if err := c.ShouldBindJSON(&result); err != nil {
		Error(c, http.StatusBadRequest, "invalid_request", err.Error())
		return result, false
	}
	return result, true
}

// JSON sends a JSON response with the given status code and data.
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

// Error sends a standardized RFC 7807 Problem Details JSON error response.
func Error(c *gin.Context, status int, code, detail string) {
	c.JSON(status, gin.H{
		"type":   "https://httpstatuses.com/" + http.StatusText(status),
		"title":  http.StatusText(status),
		"status": status,
		"detail": detail,
		"instance": c.Request.URL.Path,
	})
}