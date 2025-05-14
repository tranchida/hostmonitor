package main

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tranchida/hostmonitor/internal/handler"
)

//go:embed template static
var contentFS embed.FS

func newEngine() *echo.Echo {

	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	staticfs, _ := fs.Sub(contentFS, "static")
	e.StaticFS("/static", staticfs)

	e.GET("/", handler.IndexHandler)
	e.GET("/host", handler.HostInfoHandler)

	return e
}

func main() {
	fmt.Println("open browser on this url : http://localhost:8080")

	if err := newEngine().Start(":8080"); err != nil {
		panic(err)
	}

}
