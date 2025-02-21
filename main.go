package main

import (
	"embed"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tranchida/echotest/internal/handler"
)

//go:embed static
var contentFS embed.FS

func newEcho() *echo.Echo {

	e := echo.New()

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           "${remote_ip} - - [${time_custom}] \"${method} ${uri} ${protocol}\" ${status} ${bytes_in} ${bytes_out} ${latency_human}\n",
		CustomTimeFormat: "02/Jan/2006:15:04:05 -0700",
	}))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Filesystem: http.FS(contentFS),
	}))

	e.GET("/host", handler.HostInfoHandler)

	return e
}

func main() {

	fmt.Println("open browser on : http://localhost:8080")

	if err := newEcho().Start(":8080"); err != nil {
		panic(err)
	}

}
