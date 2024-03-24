package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/iBoBoTi/ats/server"
	"net/http"
	"os"
)

type healthController struct {
	srv *server.Server
}

func NewHealthController(srv *server.Server) *healthController {
	return &healthController{srv: srv}
}
func (c *healthController) HealthCheck(ctx *gin.Context) {

	srv := c.srv

	healthOutput := struct {
		Status    string `json:"status"`
		Env       string `json:"env"`
		Host      string `json:"host"`
		Version   string `json:"version"`
		BuildTime string `json:"build_time"`
	}{
		Status:    "ok",
		Env:       srv.Config.Environment,
		Host:      getOSHostName(),
		Version:   srv.Version,
		BuildTime: srv.BuildTime,
	}
	ctx.JSON(http.StatusOK, healthOutput)
}

func getOSHostName() string {
	osHostName, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return osHostName
}
