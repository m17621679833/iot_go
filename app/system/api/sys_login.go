package system

import (
	"github.com/gin-gonic/gin"
	base "iot_go/base/db"
	"iot_go/middleware"
)

type SysLoginApi struct {
}

func RegisterSysLoginApi(systemGroup *gin.RouterGroup) {
	systemLogin := &SysLoginApi{}
	systemGroup.GET("/login", systemLogin.UserLogin)
}

// UserLogin godoc
// @Summary 用户登陆
// @Description 用户登陆
// @Tags 用户登陆接口
// @ID /system/login
// @Accept  json
// @Produce  json
// @Param body dto.UserLoginInput true "body"
// @Success 200 {object} middleware.Response{data=dto.UserTokenOutput} "success"
// @Router /system/auth/login [post]
func (login *SysLoginApi) UserLogin(c *gin.Context) {
	base.GetDefaultDB(c)
	middleware.ResponseSuccess(c, true)
}
