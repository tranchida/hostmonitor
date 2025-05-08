package handler

import (
	"github.com/a-h/templ"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/hostmonitor/internal/model"
	"github.com/tranchida/hostmonitor/template"	
)

func render(ctx *gin.Context, status int, template templ.Component) error {
	ctx.Status(status)
	return template.Render(ctx.Request.Context(), ctx.Writer)
}

// IndexHandler handles the / route.
func IndexHandler(c *gin.Context) {
	render(c, http.StatusOK, template.Index())
}

// HostInfoHandler handles the /host route.
func HostInfoHandler(c *gin.Context) {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	render(c, http.StatusOK, template.Host(hostInfo))
}
