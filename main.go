package main

import (
	"embed"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HostInfo struct {
	CurrentTime string
	Hostname    string
	Port        string
}

//go:embed static templates
var contentFS embed.FS

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newEcho() (*echo.Echo, error) {

	staticFS, err := fs.Sub(contentFS, "static")
	if err != nil {
		return nil, err
	}

	templates, err := template.ParseFS(contentFS, "templates/*")
	if err != nil {
		return nil, err
	}

	e := echo.New()

	e.Renderer = &templateRenderer{templates: templates}

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_rfc3339} ${method} ${uri} ${status} ${latency_human} ${bytes_out}\n",
	}))

	e.GET("/*", echo.WrapHandler(http.FileServer(http.FS(staticFS))))

	e.GET("/host", func(c echo.Context) error {
		hostname, _ := os.Hostname()
		port := "8080" // You can change this to your desired port

		hostInfo := HostInfo{
			CurrentTime: time.Now().Format(time.RFC3339),
			Hostname:    hostname,
			Port:        port,
		}

		return c.Render(http.StatusOK, "main", hostInfo)

	})

	return e, nil
}

func main() {

	e, err := newEcho()
	if err != nil {
		panic(err)
	}
	
	if err := e.Start(":8080"); err != nil {
		panic(err)
	}

}
