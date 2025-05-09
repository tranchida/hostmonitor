package handler

import (
	"github.com/a-h/templ"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/hostmonitor/internal/model"
	"github.com/tranchida/hostmonitor/internal/telemetry"
	"github.com/tranchida/hostmonitor/template"
	"go.opentelemetry.io/otel/attribute"	
)

func render(ctx *gin.Context, status int, template templ.Component) error {
	// Create a span for rendering the template
	renderCtx, span := telemetry.StartSpan(ctx.Request.Context(), "handler.render")
	defer span.End()
	
	ctx.Status(status)
	return template.Render(renderCtx, ctx.Writer)
}

// IndexHandler handles the / route.
func IndexHandler(c *gin.Context) {
	// Get the current context with tracing information
	ctx := c.Request.Context()
	
	// Create a span for this handler
	ctx, span := telemetry.StartSpan(ctx, "handler.IndexHandler")
	defer span.End()
	
	render(c, http.StatusOK, template.Index())
}

// HostInfoHandler handles the /host route.
func HostInfoHandler(c *gin.Context) {
	// Get the current context with tracing information
	ctx := c.Request.Context()
	
	// Create a span for this handler
	ctx, span := telemetry.StartSpan(ctx, "handler.HostInfoHandler")
	defer span.End()
	
	// Get host information with the traced context
	hostInfo, err := model.GetHostInfo(ctx)
	if err != nil {
		// Record the error in the current span
		telemetry.RecordError(ctx, err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	
	// Add some attributes to the span about the host info
	telemetry.AddAttributes(ctx, 
		attribute.String("host.name", hostInfo.Hostname),
		attribute.String("host.os", hostInfo.OS),
		attribute.String("host.platform", hostInfo.Platform))
	
	render(c, http.StatusOK, template.Host(hostInfo))
}
