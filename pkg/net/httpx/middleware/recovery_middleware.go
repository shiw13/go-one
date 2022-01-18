package middleware

import (
	"net/http"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
	"github.com/shiw13/go-one/pkg/logger"
	"github.com/shiw13/go-one/pkg/util/bytesconv"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var reqStr string
				if c.Request != nil {
					if bs, e := httputil.DumpRequest(c.Request, true); e == nil {
						reqStr = bytesconv.UnsafeBytesToString(bs)
					}
				}
				logger.DPanicf("http server panic: %s\n%v\n", reqStr, err)
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
