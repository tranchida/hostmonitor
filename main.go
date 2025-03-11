package main

import (
	"embed"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/echotest/internal/handler"
)

//go:embed template
var contentFS embed.FS

func newEngine() *gin.Engine {

	e := gin.New()
	e.Use(gin.Logger())

	//fs, _ := fs.Sub(contentFS, "static")
	//e.StaticFS("/static", http.FS(fs))

	e.SetHTMLTemplate(template.Must(template.ParseFS(contentFS, "template/*.html")))

	e.GET("/", handler.IndexHandler)
	e.GET("/host", handler.HostInfoHandler)

	return e
}

func main() {

	fmt.Println("open browser on : http://localhost:8080")

	if err := newEngine().Run(":8080"); err != nil {
		panic(err)
	}

}
