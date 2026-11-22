package httpx

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// BindJSON unmarshals the request body into T and reports whether it succeeded.
func BindJSON[T any](c *gin.Context) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil {
		return req, false
	}
	return req, true
}

// JSON writes a JSON body with the given status.
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

// Error writes a standardized RFC 7807 problem JSON body.
func Error(c *gin.Context, status int, code, detail string) {
	c.JSON(status, gin.H{
		"type":   "about:blank",
		"title":  http.StatusText(status),
		"status": status,
		"detail": detail,
	})
}

// StatusText returns the RFC 7807 title for a status code.
func StatusText(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusMethodNotAllowed:
		return "method_not_allowed"
	case http.StatusConflict:
		return "conflict"
	case http.StatusInternalServerError:
		return "internal_error"
	default:
		return http.StatusText(status)
	}
}

// UnmarshalTypeError is the error shape ShouldBindJSON can produce.
type UnmarshalTypeError struct {
	FieldValue string
	Kind      string
	Err       error
}

func (u UnmarshalTypeError) Error() string {
	return fmt.Sprintf("json: cannot unmarshal %s into Go value of type %s", u.Err, u.Kind)
}
