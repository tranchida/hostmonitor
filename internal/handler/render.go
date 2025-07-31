package handler

import (
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func render(ctx echo.Context, status int, template templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := template.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(status, buf.String())
}