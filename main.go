package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tranchida/hostmonitor/internal/handler"
	"go.uber.org/zap"
)

//go:embed template static
var contentFS embed.FS

//go:embed all:frontend/build
var frontendFS embed.FS

func newEngine() *echo.Echo {

	e := echo.New()

	logformat := "${remote_ip} - - [${time_custom}] \"${method} ${path} ${protocol}\" ${status} ${bytes_out} \"${user_agent}\"\n"
	customTimeFormat := "2/Jan/2006:15:04:05 -0700"

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           logformat,
		CustomTimeFormat: customTimeFormat,
	}))

	// CORS middleware for SvelteKit dev server
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.OPTIONS},
	}))

	// Legacy htmx routes (kept for backward compatibility)
	staticfs, _ := fs.Sub(contentFS, "static")
	e.StaticFS("/static", staticfs)
	e.GET("/legacy", handler.IndexHandler)
	e.GET("/host", handler.HostInfoHandler)

	// JSON API
	e.GET("/api/hostinfo", handler.HostInfoJSONHandler)

	// SvelteKit frontend (static adapter build)
	frontendBuild, _ := fs.Sub(frontendFS, "frontend/build")
	e.StaticFS("/", frontendBuild)

	return e
}

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	fmt.Println("open browser on this url : http://localhost:" + port)

	if err := newEngine().Start(":" + port); err != nil {
		panic(err)
	}

}
