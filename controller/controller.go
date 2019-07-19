package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

/* 自定义Logger */
func ggLogger() gin.HandlerFunc {
	logClient := logrus.New()

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()

		//执行时间
		latency := end.Sub(start)

		path := c.Request.URL.Path

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		logClient.Infof("| %3d | %13v | %15s | %s  %s |",
			statusCode,
			latency,
			clientIP,
			method, path,
		)
	}
}

/* 服务路由 */
func ServeRouter() *gin.Engine {
	router := gin.New()
	router.Use(ggLogger())
	router.Use(gin.Recovery())
	initController(router)
	return router
}

/* 调试路由 */
func DebugRouter() *gin.Engine {
	router := gin.New()
	router.Use(ggLogger())
	router.Use(gin.Recovery())
	initController(router)
	return router
}

func initController(r *gin.Engine) {
	appGroup := r.Group("/abLin", func(context *gin.Context) {
		setRequestId(context)
	})
	initProductController(appGroup)
}
