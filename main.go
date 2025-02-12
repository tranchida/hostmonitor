package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/shirou/gopsutil/v3/host"
)

type HostInfo struct {
	CurrentTime     string
	Hostname        string
	Port            string
	Uptime          string
	OS              string
	Platform        string
	PlatformVersion string
}

//go:embed static templates
var contentFS embed.FS

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newEcho() (*echo.Echo, error) {

	logger := zerolog.New(os.Stdout)

	e := echo.New()
	e.HideBanner = false

	e.Renderer = &Template{
		templates: template.Must(template.ParseFS(contentFS, "templates/*")),
	}

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info().
				Str("URI", v.URI).
				Int("status", v.Status).
				Msg("request")

			return nil
		},
	}))

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:       "static",
		Filesystem: http.FS(contentFS),
	}))

	e.GET("/host", func(c echo.Context) error {
		hostname, _ := os.Hostname()
		port := "8080"

		// extract information from the current host using gopsutil
		uptime, _ := host.Uptime()
		info, _ := host.Info()

		hostInfo := HostInfo{
			CurrentTime:     time.Now().Format(time.RFC3339),
			Hostname:        hostname,
			Port:            port,
			Uptime:          formatDuration(time.Duration(time.Duration(uptime).Seconds())),
			OS:              info.OS,
			Platform:        info.Platform,
			PlatformVersion: info.PlatformVersion,
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

func formatDuration(d time.Duration) string {
	days := d / (24 * time.Hour)
	d -= days * 24 * time.Hour
	hours := d / time.Hour
	d -= hours * time.Hour
	minutes := d / time.Minute

	return fmt.Sprintf("%d jours, %d heures, %d minutes", days, hours, minutes)
}
