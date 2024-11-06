package system

import (
	"github.com/gin-gonic/gin"
	SystemLoginController "iot_go/app/system/api"
	"iot_go/middleware"
)

func Router(baseRouter *gin.Engine) {
	group := baseRouter.Group("/system")
	group.Use(middleware.RequestLog())
	SystemLoginController.RegisterSysLoginApi(group)
}
