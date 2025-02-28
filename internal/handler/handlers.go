package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tranchida/echotest/internal/component"
	"github.com/tranchida/echotest/internal/model" // Updated import path
)

// IndexHandler handles the / route.
func IndexHandler(c echo.Context) error {
	return Render(c, http.StatusOK, component.Index())
}

// HostInfoHandler handles the /host route.
func HostInfoHandler(c echo.Context) error {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		return err // Proper error handling; return the error to Echo
	}
	return Render(c, http.StatusOK, component.HostDisplay(hostInfo))
}
