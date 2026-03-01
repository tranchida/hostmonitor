package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/tranchida/hostmonitor/internal/model"
	"github.com/tranchida/hostmonitor/template"
)

// IndexHandler handles the / route.
func IndexHandler(c echo.Context) error {
	return render(c, http.StatusOK, template.Index())
}

// HostInfoHandler handles the /host route (htmx HTML fragment).
func HostInfoHandler(c echo.Context) error {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		return err
	}
	return render(c, http.StatusOK, template.Host(hostInfo))
}

// HostInfoJSONHandler handles the /api/hostinfo route (JSON).
func HostInfoJSONHandler(c echo.Context) error {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, hostInfo)
}
