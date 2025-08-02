package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tranchida/hostmonitor/internal/model"
	"go.uber.org/zap"
)

//go:embed template static
var contentFS embed.FS

type Template struct {
    templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func newEngine() *echo.Echo {

	e := echo.New()

	logformat := "${remote_ip} - - [${time_custom}] \"${method} ${path} ${protocol}\" ${status} ${bytes_out} \"${user_agent}\"\n"
	customTimeFormat := "2/Jan/2006:15:04:05 -0700"

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           logformat,
		CustomTimeFormat: customTimeFormat,
	}))

	staticfs, _ := fs.Sub(contentFS, "static")
	e.StaticFS("/static", staticfs)
	e.Renderer = &Template{
		templates: template.Must(template.ParseFS(contentFS, "template/*.tmpl")),
	}

	e.GET("/", IndexHandler)
	e.GET("/host", HostInfoHandler)

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

// IndexHandler handles the / route.
func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.tmpl", nil)
}

// HostInfoHandler handles the /host route.
func HostInfoHandler(c echo.Context) error {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		return err
	}
	return c.Render(http.StatusOK, "host.tmpl", hostInfo)
}
