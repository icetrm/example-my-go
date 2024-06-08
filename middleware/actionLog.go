package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func ActionLog() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		hostName, err := os.Hostname()
		if err != nil {
			hostName = ""
		}
		return fmt.Sprintf("%s|%s|%s|%s|%s|%s|%d|%s|%s|%s\n",
			hostName,
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
