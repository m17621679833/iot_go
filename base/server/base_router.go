package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	config "iot_go/base/conf"
	"iot_go/base/docs"
	"iot_go/middleware"
	appRouter "iot_go/router"
)

func InitRoute() *gin.Engine {
	docs.SwaggerInfo.Title = config.Conf.SwaggerConfig.Title
	docs.SwaggerInfo.Description = config.Conf.SwaggerConfig.Desc
	docs.SwaggerInfo.Version = config.Conf.SwaggerConfig.Version
	docs.SwaggerInfo.Host = config.Conf.SwaggerConfig.Host
	docs.SwaggerInfo.BasePath = config.Conf.SwaggerConfig.BasePath
	docs.SwaggerInfo.Schemes = config.Conf.SwaggerConfig.Schemes
	router := gin.Default()
	router.Use(middleware.Cors())
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	appRouter.RegisterRouter(router)
	return router
}
