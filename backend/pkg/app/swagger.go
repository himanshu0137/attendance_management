package app

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	docs "hunrx.com/attendance_management/docs"
)

// SetupSwagger sets up Swagger documentation for the API
func SetupSwagger(engine *gin.Engine) {
	docs.SwaggerInfo.Version = "v1"
	// Swagger route
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
