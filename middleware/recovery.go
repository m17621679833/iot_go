package middleware

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	conf "iot_go/base/conf"
	logger "iot_go/base/trace"
	"iot_go/base/util"
	"runtime/debug"
)

// RecoveryMiddleware捕获所有panic，并且返回错误信息
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.NiceLogger.Error(util.GetGinTraceContext(c), "_com_panic", map[string]interface{}{
					"error": fmt.Sprint(err),
					"stack": string(debug.Stack()),
				})
				if conf.Conf.BaseInfo.DebugMode != "debug" {
					ResponseError(c, 500, errors.New("内部错误"))
					return
				} else {
					ResponseError(c, 500, errors.New(fmt.Sprint(err)))
					return
				}
			}
		}()
		c.Next()
	}
}
