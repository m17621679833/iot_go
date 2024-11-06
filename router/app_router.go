package router

import (
	"github.com/gin-gonic/gin"
	systemRouter "iot_go/app/system/router"
)

func RegisterRouter(router *gin.Engine) {
	systemRouter.Router(router)
}
