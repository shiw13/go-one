package middleware

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shiw13/go-one/pkg/logger"
	"github.com/shiw13/go-one/pkg/util/bytesconv"
	"go.uber.org/zap"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func AccessLogMiddleware(l *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBodyBytes []byte
		if c.Request.Body != nil {
			reqBodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
		}

		w := &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
		c.Writer = w

		begin := time.Now()
		c.Next()
		latency := time.Since(begin)

		statusCode := c.Writer.Status()

		var reqBodyStr string
		if reqBodyBytes != nil {
			reqBodyStr = bytesconv.UnsafeBytesToString(reqBodyBytes)
		}

		const maxPrintLen = 10 * 1024
		var resBodyStr string
		if c.Writer.Header().Get("Content-Encoding") != "" {
			resBodyStr = "compressed body"
		} else {
			if w.body.Len() < maxPrintLen {
				resBodyStr = w.body.String()
			} else {
				resBodyStr = "body len over limit"
			}
		}

		logFunc := l.Zap().Info
		if statusCode >= 400 {
			logFunc = l.Zap().Warn
		}

		logFunc("http server",
			zap.String(logger.SrcAddr, c.Request.RemoteAddr),
			zap.String(logger.HTTPMethod, c.Request.Method),
			zap.String(logger.URLOriginal, c.Request.RequestURI),
			zap.String(logger.UserAgent, c.Request.UserAgent()),
			zap.String(logger.HTTPReqBody, reqBodyStr),
			zap.Int(logger.HTTPStatusCode, statusCode),
			zap.String(logger.HTTPResBody, resBodyStr),
			zap.Int64("latency", latency.Milliseconds()),
		)
	}
}
