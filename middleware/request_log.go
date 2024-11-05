package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	config "iot_go/base/conf"
	trace "iot_go/base/trace"
	"time"
)

func RequestInLog(c *gin.Context) {
	traceContext := trace.NewTrace()
	if rId := c.Request.Header.Get("com-header-rid"); rId != "" {
		traceContext.TraceId = rId
	}
	if spanId := c.Request.Header.Get("com-header-spanid"); spanId != "" {
		traceContext.SpanId = spanId
	}
	c.Set("startExecTime", time.Now())
	c.Set("trace", traceContext)

	bodyBytes, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	trace.NiceLogger.Info(traceContext, "_com_request_in", map[string]interface{}{
		"uri":    c.Request.RequestURI,
		"method": c.Request.Method,
		"args":   c.Request.PostForm,
		"body":   string(bodyBytes),
		"from":   c.ClientIP(),
	})
}

func RequestOutLog(c *gin.Context) {
	// after request
	endExecTime := time.Now()
	response, _ := c.Get("response")
	st, _ := c.Get("startExecTime")
	context, _ := c.Get("trace")
	startExecTime, _ := st.(time.Time)
	trace.NiceLogger.Info(context.(*trace.TraceContext), "_com_request_out", map[string]interface{}{
		"uri":       c.Request.RequestURI,
		"method":    c.Request.Method,
		"args":      c.Request.PostForm,
		"from":      c.ClientIP(),
		"response":  response,
		"proc_time": endExecTime.Sub(startExecTime).Seconds(),
	})
}

func RequestLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Conf.LogConfig.FW.On {
			RequestInLog(c)
			defer RequestOutLog(c)
		}
		c.Next()
	}
}
