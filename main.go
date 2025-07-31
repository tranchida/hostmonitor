package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/hostmonitor/internal/handler"
	"go.uber.org/zap"
)

//go:embed template static
var contentFS embed.FS

func newEngine(logger *zap.Logger) *gin.Engine {

	e := gin.New()
	e.Use(LoggerMiddleware(logger))
	e.Use(gin.Recovery())

	staticfs, _ := fs.Sub(contentFS, "static")
	e.StaticFS("/static", http.FS(staticfs))

	e.GET("/", handler.IndexHandler)
	e.GET("/host", handler.HostInfoHandler)

	return e
}

func main() {
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	fmt.Println("ouvrir le browser sur : http://localhost:8080")

	if err := newEngine(logger).Run(":8080"); err != nil {
		panic(err)
	}

}

func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		logger.Info(
			"HTTP",
			zap.String("method", c.Request.Method),
			zap.Int("status", c.Writer.Status()),
			zap.String("path", c.Request.URL.Path),
			zap.Duration("duration", time.Since(start)),
		)

	}
}
