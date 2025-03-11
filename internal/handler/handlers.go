package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranchida/hostmonitor/internal/model"
)

// IndexHandler handles the / route.
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", nil)
}

// HostInfoHandler handles the /host route.
func HostInfoHandler(c *gin.Context) {
	hostInfo, err := model.GetHostInfo()
	if err != nil {
		c.Error(err) 
	}
	c.HTML(http.StatusOK, "host", hostInfo)
}
